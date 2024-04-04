package service

import (
	"context"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

//😎😎😍😍😘🥰
// useservice acts as a struct for injecting an implementation of userRepository
// for use in service methods

type userService struct {
	UserRepository  model.UserRespository
	ImageRepository model.ImageRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer

type USConfig struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		ImageRepository: c.ImageRepository,
	}

}

func (s *userService) clearProfileImage(
	ctx context.Context,
	uid uuid.UUID,
) error {
	user
}
