// Copyright © ivanlobanov. All rights reserved.
package route

import (
	"net/http"

	"github.com/cantylv/service-happy-birthday/internal/middlewares"
	"github.com/cantylv/service-happy-birthday/internal/route/auth"
	"github.com/cantylv/service-happy-birthday/internal/route/sub"
	"github.com/cantylv/service-happy-birthday/services"
	"github.com/gorilla/mux"
)

type RouterProps struct {
	R *mux.Router
	S services.Services
}

func Initialize(p RouterProps) http.Handler {
	auth.Init(auth.AuthProps{
		Router: p.R,
	})
	sub.Init(sub.SubProps{
		Router: p.R,
	})
	return middlewares.Init(middlewares.MiddlewaresProps{
		Router: p.R,
	})
}
