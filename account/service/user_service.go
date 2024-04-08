// 😎😎😍😍😘🥰
// useservice acts as a struct for injecting an implementation of userRepository
// for use in service methods
package service

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"

	"github.com/google/uuid"
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

//func (s *userService) ClearProfileImage(
//	ctx context.Context,
//	uid uuid.UUID,
//) error {
//	user, err := s.UserRepository.FindByID(ctx, uid)
//
//	if err != nil {
//		return err
//	}
//
//	if user.ImageURL == "" {
//		return nil
//	}
//
//	objName, err := ObjNameFromURL(user.ImageURL)
//
//	if err != nil {
//		return err
//	}
//
//	err = s.ImageRepository.DeleteProfile(ctx, objName)
//	if err != nil {
//		return err
//	}
//
//	_, err = s.UserRepository.UpdateImage(ctx, uid, "")
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func ObjNameFromURL(imageURL string) (string, error) {
//	if imageURL == "" {
//		objID, _ := uuid.NewRandom()
//		return objID.String(), nil
//	}
//
//	urlPath, err := url.Parse(imageURL)
//
//	if err != nil {
//		log.Printf("Failed to parse objectName from imageURL: %v\n", imageURL)
//		return "", apperrors.NewInternal()
//	}
//
//	return path.Base(urlPath.Path), nil
//}

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
	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	return nil
}

func (s *userService) SetProfileImage(
	ctx context.Context,
	uid uuid.UUID,
	imageFileHeader *multipart.FileHeader,
) (*model.User, error) {
	return nil, nil
}
