package middlewares

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
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
		logger := zap.Must(zap.NewProduction())

		requestId := functions.GetCtxRequestId(r)
		jwtToken, err := functions.GetJWtToken(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(fmt.Sprintf("error while jwt getting: %v", err),
				zap.String(myconstants.RequestId, requestId))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.ErrInternal.Error(),
				CodeStatus: http.StatusInternalServerError,
			})
			return
		}

		isMutatingMethod := false
		for _, method := range []string{"PUT", "POST"} {
			if r.Method == method {
				isMutatingMethod = true
			}
		}
		if isMutatingMethod && jwtToken != "" {
			csrfToken := r.Header.Get(myconstants.CsrfHeader)
			if csrfToken == "" {
				logger.Info("no X-CSRF-Token in headers of the HTTP request, user was redirected to the authorization form",
					zap.String(myconstants.RequestId, requestId))
				w.Header().Add("Location", "/api/v1/signin")
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "you were redirected to the authorization form",
					CodeStatus: http.StatusForbidden,
				})
				return
			}
			// Csrf-Token validation.
			isValid, err := isValidCsrfToken(tokens{
				JWTtoken:  jwtToken,
				CsrfToken: csrfToken,
			})
			if err != nil {
				logger.Error(
					fmt.Sprintf("csrf-token validation error: %v", err),
					zap.String(myconstants.RequestId, requestId),
				)
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "unexpected internal server error, try again, please",
					CodeStatus: http.StatusInternalServerError,
				})
				return
			}
			if !isValid {
				logger.Info("invalid csrf-token, user was redirected to the authorization form",
					zap.String(myconstants.RequestId, requestId),
				)
				w.Header().Add("Location", "/api/v1/signin")
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        "invalid csrf-token, you were redirected to the authorization form",
					CodeStatus: http.StatusForbidden,
				})
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
	token := hex.EncodeToString([]byte(hash))
	if token != ts.CsrfToken {
		return false, nil
	}
	return true, nil
}
