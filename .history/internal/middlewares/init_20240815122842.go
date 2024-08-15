// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

type MiddlewaresProps struct {
	Router *mux.Router
}

// Init
// Initializes the chain of middlewares.
func Init(props MiddlewaresProps) (h http.Handler) {
	h = JwtVerification(props.Router)
	h = Csrf(h)
	h = Cors(h)
	h = Access(h)
	return h
}
