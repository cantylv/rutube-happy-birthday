package user

import (
	"github.com/cantylv/service-happy-birthday/internal/delivery/user"
	rUser "github.com/cantylv/service-happy-birthday/internal/repository/user"
	uUser "github.com/cantylv/service-happy-birthday/internal/usecase/user"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProps struct {
	Router     *mux.Router
	Collection *mongo.Collection
}

func Init(props UserProps) {
	repoUser := rUser.NewRepoLayer(props.Collection)
	usecaseUser := uUser.NewUsecaseLayer(&repoUser)
	deliveryUser := user.NewDeliveryLayer(&usecaseUser)
	// need to define delivery struct for auth handlers
	props.Router.HandleFunc("/api/v1/user", deliveryUser.GetUser).Methods("GET")
	props.Router.HandleFunc("/api/v1/user", deliveryUser.UpdateUser).Methods("PUT")
}
