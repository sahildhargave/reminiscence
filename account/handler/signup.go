// 😁😂🤣😎😍😘🥰
package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, gt=6,lte=30"`
}

func (h *Handler) Signup(c *gin.Context) {
	// TODO Define a variable to which we 'll bind incomiing

	// json body, { email, password}
	var req signupReq

	//Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.UserService.Signup(ctx, u)

	if err != nil {
		log.Printf("Failed to sign up  user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	//create token pair as strings

	//	tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")
	//
	//	if err != nil {
	//		log.Printf("Failed to create tokens for user: %v\n", err.Error())
	//
	//		//may eventually implement rollback logic here
	//		// means if we fail to crate tokens after create a user
	//		// we make sure to clear/delete the created user in the database
	//
	//		c.JSON(apperrors.Status(err), gin.H{
	//			"error" : err,
	//		})
	//		return
	//	}
	//
	//	c.JSON(http.StatusCreated, gin.H{
	//		"token": tokens,
	//	})
	//
	//

}
