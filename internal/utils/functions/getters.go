package functions

import (
	"net/http"
	"strings"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
)

// Context

// Request-ID
func GetCtxRequestId(r *http.Request) string {
	ctxRequestId := r.Context().Value(myconstants.RequestId)
	if ctxRequestId == nil {
		return ""
	}
	return ctxRequestId.(string)
}

// HTTP Headers

// Authorization
func GetJWtToken(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	// Header must start from "Bearer "
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authorizationHeader, "Bearer ")
}
