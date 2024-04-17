//😎😎😎😎

package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/sahildhargave/memories/account/model"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

// Create a mocks concrete NewPairFromUser

func (m *MockTokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {

	ret := m.Called(ctx, u, prevTokenID)

	var r0 *model.TokenPair

	if ret.Get(0) != nil {
		//
		r0 = ret.Get(0).(*model.TokenPair)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}


//😭😭😭
func (m *MockTokenService) Signout(ctx context.Context, uid uuid.UUID) error{
	ret := m.Called(ctx, uid)
	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}


// Validate ID Token mocks
func (m *MockTokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	ret := m.Called(tokenString)

	// TODO
	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)

	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// validating refresh token validaterefresh token
// 😎😎😎
func (m *MockTokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {

	ret := m.Called(refreshTokenString)

	var r0 *model.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.RefreshToken)

	}
	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)

	}
	return r0, r1
}
