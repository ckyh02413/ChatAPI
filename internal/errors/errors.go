package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrMessageNotFound    = errors.New("message not found")
	ErrChatroomNotFound   = errors.New("chatroom not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrForbidden          = errors.New("permission denied")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func ErrorToStatus(err error) int {
	switch err {
	case ErrChatroomNotFound, ErrMessageNotFound, ErrUserNotFound:
		return http.StatusNotFound
	case ErrForbidden:
		return http.StatusForbidden
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrInvalidCredentials:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
