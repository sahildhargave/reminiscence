package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// TODO VALIDATE ERROR
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// Authorization User extract from the authorization
func AuthUser(s model.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		// 😎😎😎
		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			// otherwise error type is unknown
			err := apperrors.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		// TODO validate
		user, err := s.ValidateIDToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
