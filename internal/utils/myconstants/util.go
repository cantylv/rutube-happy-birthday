// Copyright © ivanlobanov. All rights reserved.
package myconstants

import (
	"time"
)

// Default config parameters.
const (
	DefaultInterval = 1
	TimeExpDur      = 14 * 24 * time.Hour
)

// Naming.
const (
	RequestId     = "request-id"
	UserId        = "user-id"
	CsrfHeader    = "X-CSRF-Token"
	JwtCookie     = "jwt-token"
	EnvCsrfSecret = "CSRF_SECRET"
	EnvJwtSecret  = "JWT_SECRET"
	EnvSecretKey  = "KEY_SECRET"
)
