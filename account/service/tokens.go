package service

import (
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sahildhargave/memories/account/model"
)

// idTokenCustomClaims holds structure of jwt claims of token
type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

// generate ID Token an IDToken which is a jwt with mycustomClaims
// Could call this Generate ID Token String , but the signature makes this fairly clear

func generateIDToken(u *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := idTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

// TODO // refreshToken holds the actual signed jwt string along with ID
// TODO // return id it can be used without re-parsing the jwt from signed string

type refreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

// TODO Refresh Token Cutom Claims holds the payload of a refresh token
// This used to extract user id for subsequent
// TODO operation (IE, Fetch user in Redis)

type refreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// generate Refresh Token creates a refresh token

// the refresh TOken stores only the user's ID , A string

func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &refreshTokenData{
		SS:        ss,
		ID:        tokenID,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

// TODO validating ID Token returns the token claims if the token is valid

func validateIDToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {

	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// TODO Handle logging service
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but could't parse claims")
	}

	return claims, nil
}

// validating refresh Token uses the secret key
func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {

	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	//😥😥😥
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Refresh token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("Refresh token valid but couldn't parse claims")
	}

	return claims, nil
}
