// Copyright © ivanlobanov. All rights reserved.
package auth

import (
	"github.com/cantylv/service-happy-birthday/internal/delivery/auth"
	"github.com/cantylv/service-happy-birthday/internal/repository/user"
	uAuth "github.com/cantylv/service-happy-birthday/internal/usecase/auth"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthProps struct {
	Router     *mux.Router
	Collection *mongo.Collection
}

func Init(props AuthProps) {
	repoUser := user.NewRepoLayer(props.Collection)
	usecaseAuth := uAuth.NewUsecaseLayer(&repoUser)
	deliveryAuth := auth.NewDeliveryLayer(&usecaseAuth)
	// need to define delivery struct for auth handlers
	props.Router.HandleFunc("/api/v1/signup", deliveryAuth.SignUp).Methods("POST")
	props.Router.HandleFunc("/api/v1/signin", deliveryAuth.SignIn).Methods("POST")
	props.Router.HandleFunc("/api/v1/signout", deliveryAuth.SignOut).Methods("POST")
}
