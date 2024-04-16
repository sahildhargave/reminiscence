// 😎😎😍😍😘🥰
// useservice acts as a struct for injecting an implementation of userRepository
// for use in service methods
package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

type userService struct {
	UserRepository model.UserRepository
	//ImageRepository model.ImageRepository
}

type USConfig struct {
	UserRepository model.UserRepository
	//ImageRepository model.ImageRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
		//ImageRepository: c.ImageRepository,
	}
}

func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

func (s *userService) Signup(ctx context.Context, u *model.User) error {

	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to signup user for email:%v\n", u.Email)
		return apperrors.NewInternal()
	}

	//TODO create A user . Its un-natural to mutate the user here

	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	//if err := s.UserRepository.Create(ctx, u); err != nil {
	//	return err
	//}

	// Adding events
	//TODO learn Message broker

	return nil
}

func (s *userService) Signin(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)

	// return NotAuthorized to client to omit details of
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verifying password
	match, err := comparePasswords(uFetched.Password, u.Password)
	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")

	}

	*u = *uFetched
	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	err := s.UserRepository.Update(ctx, u)

	if err != nil {
		return err
	}

	return nil
}
