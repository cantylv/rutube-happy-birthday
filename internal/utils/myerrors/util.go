// Copyright © ivanlobanov. All rights reserved.
package myerrors

import "errors"

var (
	// Repository level
	ErrInvalidObjectId  = errors.New("object ID was not received")
	ErrUpdateFailed     = errors.New("error while updating document")
	ErrUserNotExist     = errors.New("user was not found")
	ErrUserAlreadyExist = errors.New("user already exists")
	ErrNoSubscription   = errors.New("user doesn't have subscription")

	// Delivery level
	ErrInternal               = errors.New("unexpected internal server error, please try again in one minute")
	ErrAlreadyRegistered      = errors.New("you're already registered")
	ErrAlreadyAuthenticated   = errors.New("you're already authenticated")
	ErrAuth                   = errors.New("you're not authenticated")
	ErrUserExist              = errors.New("you're not registered in our system")
	ErrInvalidJwt             = errors.New("invalid jwt-token")
	ErrInvalidRequestData     = errors.New("incorrect data received, please try again")
	ErrEmailIsReserved        = errors.New("user with this email already exist, failed")
	ErrBadCredentials         = errors.New("incorrect password or login")
	ErrBadVarPathEmployeeId   = errors.New("provided wrong value of path parameter 'employee_id'")
	ErrNoSubscriptionEmployee = errors.New("you don't have subscription")
)
