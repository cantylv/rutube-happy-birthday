package middlewares

import (
	"fmt"
	"net/http"

	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"go.uber.org/zap"
)

type tokens struct {
	JWTtoken  string
	CsrfToken string
}

// Csrf
// Middleware is used for prevent Cross Site Request Forgery(CSRF) attack.
// We use 'Double Submit Cookie pattern' (for stateless app). Implementation --> Signed Double-Submit Cookie.
// For hashing we use Hash-based Message Authentication (HMACSHA256).
func Csrf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zap.Must(zap.NewProduction()).Sugar()
		requestId := functions.GetCtxRequestId(r)
		jwtToken := functions.GetJWtToken(r)

		isMutatingMethod := false
		for _, method := range []string{"PUT", "POST"} {
			if r.Method == method {
				isMutatingMethod = true
			}
		}
		if isMutatingMethod && jwtToken != "" {
			csrfToken := r.Header.Get(myconstants.CsrfHeader)
			if csrfToken == "" {
				logger.Error("No X-CSRF-Token in headers of the HTTP request. User was redirected to the authorization form.",
					zap.String(myconstants.RequestId, requestId))
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "You were redirected to the authorization form.",
					CodeStatus: http.StatusForbidden,
				})
				w.Header().Add("Location", "/api/v1/signin")
				return
			}
			// Csrf-Token validation.
			isValid, err := isValidCsrfToken(tokens{
				JWTtoken:  jwtToken,
				CsrfToken: csrfToken,
			})
			if err != nil {
				logger.Error(
					fmt.Sprintf("CSRF-token validation error: %v", err),
					zap.String(myconstants.RequestId, requestId),
				)
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "Unexpected internal server error. Try again, please!",
					CodeStatus: http.StatusInternalServerError,
				})
				return
			}
			if !isValid {
				logger.Error("Invalid CSRF-Token. User was redirected to the authorization form.",
					zap.String(myconstants.RequestId, requestId),
				)
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "Invalid CSRF-Token. You were redirected to the authorization form.",
					CodeStatus: http.StatusForbidden,
				})
				w.Header().Add("Location", "/api/v1/signin")
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func isValidCsrfToken(ts tokens) (bool, error) {
	hash, err := functions.HashWithStatement(functions.HashProps{
		EnvName:   myconstants.EnvCsrfSecret,
		Statement: ts.JWTtoken,
	})
	if err != nil {
		return false, err
	}
	if hash != ts.CsrfToken {
		return false, nil
	}
	return true, nil
}
