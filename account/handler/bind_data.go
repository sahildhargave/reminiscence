package handler

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/sahildhargave/memories/account/model/apperrors"

	"github.com/gin-gonic/gin"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param`
}

// bindData is helper function , returns false if data is not bound

func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type applicaiton/json", c.FullPath())

		err := apperrors.NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			//could probably extract this, it is also in middleware_auth_user
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invaild request parameters. See invalidArgs")

			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		//later er'll add code for validating max body size here!
		// if we are't able to properly extract validation errors,
		//we'll fall back and return an internal server error

		fallBack := apperrors.NewInternal()

		c.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}

	return true
}
