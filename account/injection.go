package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sahildhargave/memories/account/handler"
	repository "github.com/sahildhargave/memories/account/respository"

	//repository "github.com/sahildhargave/memories/account/respository"
	"github.com/sahildhargave/memories/account/service"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer

func inject(d *dataSources) (*gin.Engine, error) {

	log.Println("aianjecting  data  sources")

	// 😎😎 repository layer
	userRepository := repository.NewUserRepository(d.DB)

	// 😎😎 service layer
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	// loading rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("Could Not Read Private Key Pem File: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)

	if err != nil {
		return nil, fmt.Errorf("Could Not Read Public Key Pem File: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("Could Not Parse Public Key: %w", err)
	}

	// load refresh token scret from env variable

	refreshSecret := os.Getenv("REFRESH_SECRET")

	tokenService := service.NewTokenService(&service.TSConfig{
		PrivKey:       privKey,
		PubKey:        pubKey,
		RefreshSecret: refreshSecret,
	})

	// initialize gin.Engine
	router := gin.Default()

	// reading in ACCOUNT_API_URL

	//baseURL := os.Getenv("ACCOUNT_API_URL")

	handler.NewHandler(&handler.Config{
		R:            router,
		UserService:  userService,
		TokenService: tokenService,
	})

	return router, nil

}
