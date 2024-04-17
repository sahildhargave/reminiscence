package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/sahildhargave/memories/account/model"

	"github.com/sahildhargave/memories/account/model/apperrors"
	"github.com/sahildhargave/memories/account/model/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid := uuid.New()

		mockUser := &model.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// Configure the expected method call on the mock repository
		mockUserRepository.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
			userArg := args.Get(1).(*model.User)
			userArg.UID = uid
		}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.NoError(t, err)

		// Assert that the user now has a UID assigned
		assert.NotEqual(t, uuid.Nil, mockUser.UID)

		// Assert other properties of the user as needed
		assert.Equal(t, "bob@bob.com", mockUser.Email)

		// Verify the expected method call on the mock repository
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		// Configure the expected method call on the mock repository to return an error
		mockUserRepository.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		// Assert that the expected error is returned from the signup function
		assert.EqualError(t, err, mockErr.Error())

		// Verify the expected method call on the mock repository
		mockUserRepository.AssertExpectations(t)
	})
}
