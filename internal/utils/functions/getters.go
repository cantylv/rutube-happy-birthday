package functions

import (
	"net/http"

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

// HTTP Headers "Cookie"
func GetJWtToken(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie(myconstants.JwtCookie)
	if err != nil {
		return "", err
	}
	return jwtCookie.Name, nil
}
