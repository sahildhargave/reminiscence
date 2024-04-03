package handler

import (
	"memories/model"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService model.UserService
}

type Config struct {
	R           *gin.Engine
	UserService model.UserService
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService: c.UserService,
	}

	// Create a account
	g := c.R.Group(os.Getenv("ACCOUNT_API_URL"))

	g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/signout", h.Signout)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("/details", h.Details)

	// signup
	//signin
	//signout
	//tokens
	//image
	//deleteImage
	//details
}

//func (h *Handler) Me(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"hello": "it's me",
//	})
//}

func (h *Handler) Signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signup",
	})
}

func (h *Handler) Signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signin",
	})
}

func (h *Handler) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

func (h *Handler) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

func (h *Handler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's Delete Image",
	})
}

func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's Details",
	})
}
