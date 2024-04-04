//😂🤣😂🤣

package mocks

import (
	"context"
	"memories/model"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	//args that will be passed to "Return" in the test  will when function

	// is called with a uif .Hence the name ret

	ret := m.Called(ctx, uid)

	// first value passed to "Return"
	var r0 *model.User
	if ret.Get(0) != nil {
		// we can just return this if we know we won;t be passing function to "Return"
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) ClearProfileImage(ctx context.Context, uid uuid.UUID) error {
	ret := m.Called(ctx, uid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Signup is a mock of UserService.SignUp
func (m *MockUserService) Signup(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Signin is a mock of userService.Signin
func (m *MockUserService) Signin(ctx context.Context, u *model.User) error {

	ret := m.Called(ctx, u)
	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// 😁😂🤣😃😄😅
// UPDATEDETAILS IS A MOCK OF USERSERVICES UPDATEDETAILS
func (m *MockUserService) UpdateDetails(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}

// SetProfileImage is a mock of userService.SetProfileImage
func (m *MockUserService) SetProfileImage(
	ctx context.Context,
	uid uuid.UUID,
	imageFileHeader *multipart.FileHeader,
) (*model.User, error) {

	ret := m.Called(ctx, uid, imageFileHeader)

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
