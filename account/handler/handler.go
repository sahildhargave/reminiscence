package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

type config struct {
	R *gin.Engine
}

func NewHandler(c *config) {
	h := &Handler{}

	// Create a account
	g := c.R.Group("ACCOUNT_API_URL")

	g.GET("/me", h.Me)
	// signup
	//signin
	//signout
	//tokens
	//image
	//deleteImage
	//details
}

func (h *Handler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's me",
	})
}
