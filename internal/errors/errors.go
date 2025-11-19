package errors

import (
	"errors"
)

var (
	// 400 Errors
	ErrInvalidInputData      = errors.New("Invalid input data. Please check your request payload.")
	ErrInvalidPlaceData      = errors.New("Invalid place data.")
	ErrInvalidUserData		 = errors.New("Invalid user data.")
	// 404 Errors
	ErrUserNotFound          = errors.New("User not found.")
	ErrPlaceNotFound         = errors.New("Place not found.")
	// 409 Errors
	ErrConflict              = errors.New("Username or email already exists.")
	// 500 Errors
	ErrInternalServer        = errors.New("An unexpected server error occurred.")
	ErrJSONMarshalFailed     = errors.New("Failed to process internal data.")
)