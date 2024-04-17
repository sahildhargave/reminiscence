package mocks

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// mockImageRepository is mock type of model image Repository

type MockImageRepository struct {
	mock.Mock
}

// update profile is mock type for model

func (m *MockImageRepository) UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error) {

	ret := m.Called(ctx, objName, imageFile)

	var r0 string
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1
}

func (m *MockImageRepository) DeleteProfile(ctx context.Context, objName string) error {
	ret := m.Called(ctx, objName)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
