// Copyright © ivanlobanov. All rights reserved.
package sub

import (
	"github.com/gorilla/mux"
)

type SubProps struct {
	Router *mux.Router
}

func Init(props SubProps) {
	// need to define delivery struct for sub handlers
	props.Router.HandleFunc("/api/v1/sub/{user_id}", sub.Sub)
	props.Router.HandleFunc("/api/v1/unsub/{user_id}", sub.Unsub)
}
