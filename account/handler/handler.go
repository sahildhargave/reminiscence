package handler

import (
	"net/http"
	"time"

	"github.com/sahildhargave/memories/account/handler/middleware"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

	// Create a account
	g := c.R.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.Details)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
		g.PUT("/details", h.Details)

	}

	//g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	//g.POST("/signout", h.Signout)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	//g.PUT("/details", h.Details)

	// signup
	//signin
	//signout
	//tokens
	//image
	//deleteImage
	//details
}

//
//func (h *Handler) Signin(c *gin.Context) {
//	time.Sleep(1 * time.Second)
//	c.JSON(http.StatusOK, gin.H{
//		"hello": "it's signin",
//	})
//}

//func (h *Handler) Signout(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"hello": "it's signout",
//	})
//}

//func (h *Handler) Tokens(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"hello": "it's tokens",
//	})
//}

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

//func (h *Handler) Details(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"hello": "it's Details",
//	})
//}
