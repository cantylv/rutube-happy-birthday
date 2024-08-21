// Copyright © ivanlobanov. All rights reserved.
package sub

import (
	"github.com/cantylv/service-happy-birthday/internal/delivery/sub"
	rSub "github.com/cantylv/service-happy-birthday/internal/repository/sub"
	rUser "github.com/cantylv/service-happy-birthday/internal/repository/user"
	uSub "github.com/cantylv/service-happy-birthday/internal/usecase/sub"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubProps struct {
	Router     *mux.Router
	Collection *mongo.Collection
}

func Init(props SubProps) {
	repoSub := rSub.NewRepoLayer(props.Collection)
	repoUser := rUser.NewRepoLayer(props.Collection)
	usecaseSub := uSub.NewUsecaseLayer(&repoSub, &repoUser)
	deliverySub := sub.NewDeliveryLayer(&usecaseSub)

	props.Router.HandleFunc("/api/v1/sub/{employee_id}", deliverySub.Sub).Methods("POST")
	props.Router.HandleFunc("/api/v1/sub/{employee_id}/new_interval/{interval}", deliverySub.ChangeSubInterval).Methods("PUT")
	props.Router.HandleFunc("/api/v1/unsub/{employee_id}", deliverySub.Unsub).Methods("POST")
}
