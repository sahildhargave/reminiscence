package main

import (
	"context"
	"log"
	"memories/handler"
	"net/http"
	"os"

	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting account service...")

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R: router,
	})

	router.GET("/api/account", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "Fuckers",
		})
	})

	serv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	//😋😎😍😘😘🥰😙😚
	// Shutdown Code from gin
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on Port %v\n", serv.Addr)

	// Wating for close and kill the signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// this block unti; a signal is passed into the quit channel
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown server

	log.Println("shutting down server...")
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
