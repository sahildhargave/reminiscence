package service

import (
	"context"
	"crypto/rsa"

	"github.com/sahildhargave/memories/account/model"
)

// Token Service used for injection an implementation of tokenRepository
// / for use in service methods along with keys and secrets for
// signing jwts
type TokenService struct {
	// Token Respository model.Token  Repository
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
}

// New Token Service is a factory function for
// initializing a userservice with its repository layer dependencies

type TSConfig struct {
	// Token Respository model.TokenRepository
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
}

// New Token Service is a factory function for
// initializing a userService with its repository layer dependencies

func NewTokenService(c *TSConfig) model.TokenService {
	return &TokenService{
		PrivKey:       c.PrivKey,
		PubKey:        c.PubKey,
		RefreshSecret: c.RefreshSecret,
	}
}

// NewpAIRfROMuSER CREATES  FRESH ID AND REFRESH ID AND REFRESH TOKENS FOR THE CURRENT USER
// IF A PREVIOS TOKEN IS INCLUDED, THE PREVIOUS TOKEN IS REMOVED FROM
// THE TOKEN REPOSITORY

func (s *TokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	// Need to use a repository for idToken as it is unrelated to any data source
	idToken, err := generateIDToken(u, s.PrivKey)

}
