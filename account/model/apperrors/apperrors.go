package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization   Type = "AUTHORIZATION"   // Authentication Failures -
	BadRequest      Type = "BADREQUEST"      // validation errors / BadInput
	Conflict        Type = "CONFLICT"        // Already exists (eg : =create accoujnt with existent email) - 409
	Internal        Type = "INTERNAL"        // Server (500) and fallback errors
	NotFound        Type = "NOTFOUND"        // Resource not found
	PayloadTooLarge Type = "PAYLOADTOOLARGE" // for uploading tons of JSON, or an image over the limit - 413
)

// 🤣😃😄😅😆😉😊😋😎

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Error satisfies the error for the application
// which is helpful in returning a consistent
// error type/ message from API endpoints

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}

	return http.StatusInternalServerError
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason : %v", reason),
	}
}

func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("Internal server error."),
	}
}

// NewNotFound to create an error for  404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}
