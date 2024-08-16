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
	RequestId     = "Request-ID"
	UserId        = "User-ID"
	CsrfHeader    = "X-CSRF-Token"
	JwtCookie     = "Jwt-Token"
	EnvCsrfSecret = "CSRF_SECRET"
	EnvJwtSecret  = "JWT_SECRET"
)
