package model

import (
	"context"
	"time"
	"mime/multipart"

	"github.com/google/uuid"
)

type UserService interface {
	ClearProfileImage(ctx context.Context, uid uuid.UUID) error
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
	Signin(ctx context.Context, u *User) error
	UpdateDetails(ctx context.Context, u *User) error
	SetProfileImage(ctx context.Context, uid uuid.UUID, imageFileHeader *multipart.FileHeader) (*User, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	UpdateImage(ctx context.Context, uid uuid.UUID, imageURL string) (*User, error)
}

type ImageRepository interface {
	DeleteProfile(ctx context.Context, objName string) error
	UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}

// TokenService define methods the handler layer expects to Interact
// With in regards to producing jwts as string

type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error)
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
	Signout(ctx context.Context, uid uuid.UUID) error
}

// TODO Token Repository adding
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}
