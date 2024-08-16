// Copyright © ivanlobanov. All rights reserved.
package myerrors

import "errors"

var (
	// Repository level
	ErrInvalidObjectId  = errors.New("object ID was not received")
	ErrUpdateFailed     = errors.New("error while updating document")
	ErrUserNotExist     = errors.New("user was not found")
	ErrUserAlreadyExist = errors.New("user already exists")

	// Delivery level
	ErrInternal           = errors.New("unexpected internal server error, please try again in one minute")
	ErrAlreadyRegistered  = errors.New("you're already registered")
	ErrAuth               = errors.New("you're not authenticated")
	ErrInvalidJwt         = errors.New("invalid jwt-token")
	ErrInvalidRequestData = errors.New("incorrect data received, please try again")
	ErrEmailIsReserved    = errors.New("user with this email already exist, failed to register")
	ErrBadCredentials     = errors.New("incorrect password or login")
)
