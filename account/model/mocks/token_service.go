//😎😎😎😎

package mocks

import (
	"context"

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
