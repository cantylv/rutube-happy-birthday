// Copyright © ivanlobanov. All rights reserved.
package auth

import (
	"github.com/gorilla/mux"
)

type AuthProps struct {
	Router *mux.Router
}

func Init(props AuthProps) {
	// need to define delivery struct for auth handlers
	props.Router.HandleFunc("/api/v1/signup", auth.SignUp)
	props.Router.HandleFunc("/api/v1/signin", auth.SignIn)
	props.Router.HandleFunc("/api/v1/signout", auth.SignOut)
}
