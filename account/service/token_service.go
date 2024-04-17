package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

// Token Service used for injection an implementation of tokenRepository
// / for use in service methods along with keys and secrets for
// signing jwts

// New Token Service is a factory function for
// initializing a userservice with its repository layer dependencies

type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source

	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevTokenID)

			return nil, err
		}
	}

	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// delete user's current refresh token (used when refreshing idToken)

	return &model.TokenPair{
		IDToken:      model.IDToken{SS: idToken},
		RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID, UID: u.UID},
	}, nil
}

// TODO signout
func (s *tokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	return s.TokenRepository.DeleteUserRefreshTokens(ctx, uid.String())
}

// TODO VALIDATEIDTOKEN

func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, s.PubKey)

	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error : %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}
	return claims.User, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	// validate acutal jwt with string a secret
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")

	}

	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
