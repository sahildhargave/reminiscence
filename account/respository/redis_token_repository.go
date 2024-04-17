package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

// redis Token Repository is data/ repository implementation
// of service layer Token Repository

type redisTokenRepository struct {
	Redis *redis.Client
}

// New Token Repository is factory for initializing user repository
func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// Set Refresh Token is refresh Token with an expiry time

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	// We'll store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

// Delete Refresh Token used to delete old refresh tokens
// servicesaccess this to resolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

// Deleteing User refresh tokens
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprint("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iter.Val())

	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}
	return nil
}
