package route

import (
	"github.com/cantylv/service-happy-birthday/internal/route/auth"
	"github.com/cantylv/service-happy-birthday/internal/route/sub"
	"github.com/cantylv/service-happy-birthday/services"
	"github.com/gorilla/mux"
)

type RouterProps struct {
	R *mux.Router
	S services.Services
}

func Initialize(p RouterProps) {
	// auth
	.HandleFunc("/api/v1/signup", auth.SignUp)
	m.HandleFunc("/api/v1/signin", auth.SignIn)
	m.HandleFunc("/api/v1/signout", auth.SignOut)

	// subsciption
	m.HandleFunc("/api/v1/sub/{user_id}", sub.Sub)
	m.HandleFunc("/api/v1/unsub/{user_id}", sub.Unsub)
}
