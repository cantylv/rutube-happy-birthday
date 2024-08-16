// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/utils/functions"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/myerrors"
	"go.uber.org/zap"
)

// JWT --> header.payload.signature
// header --> base64(meta_information) - rsc - random secret
// payload --> base64(payload_data)
// signature --> hmacsha256(header + . + payload + secret)

//// e.g. header
// {
// 	"exp": "02.01.2006 15:04:05 UTC-07"
// }
//// e.g. payload
// {
// 	"id": "66b89cea43ad0d6f8cf3f54e",
// }

// JwtVerification
// Needed for authentication.
func JwtVerification(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zap.Must(zap.NewProduction()).Sugar()
		requestID := functions.GetCtxRequestId(r)

		jwtToken, err := functions.GetJWtToken(r)
		if err != nil {
			logger.Error(fmt.Sprintf("Error while jwt getting: %v", err),
				zap.String(myconstants.RequestId, requestID))
			functions.ErrorResponse(functions.ErrorResponseProps{
				W:          w,
				Msg:        myerrors.Internal,
				CodeStatus: http.StatusInternalServerError,
			})
			return
		}
		if jwtToken != "" {
			isValid, uId, err := jwtTokenIsValid(jwtToken)
			if err != nil {
				logger.Error(fmt.Sprintf("Error while jwt verification: %v", err),
					zap.String(myconstants.RequestId, requestID))
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        myerrors.Internal,
					CodeStatus: http.StatusInternalServerError,
				})
				return
			}
			if !isValid {
				logger.Info(fmt.Sprintf("Invalid jwt-token: %v", err),
					zap.String(myconstants.RequestId, requestID))
				functions.ErrorResponse(functions.ErrorResponseProps{
					W:          w,
					Msg:        myerrors.InvalidJwt,
					CodeStatus: http.StatusUnauthorized,
				})
				w.Header().Add("Location", "/api/v1/signin")
				return
			}
			ctx := context.WithValue(r.Context(), myconstants.UserId, uId)
			r = r.WithContext(ctx)
		}
		// Decode payload and use data.
		h.ServeHTTP(w, r)
	})
}

// jwtTokenIsValid
// Needed for validation jwt-token.
func jwtTokenIsValid(token string) (bool, string, error) {
	// check time validation of token
	// if all is okey, return true
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, "", nil
	}
	signature, err := functions.HashWithStatement(functions.HashProps{
		EnvName:   myconstants.EnvJwtSecret,
		Statement: parts[0] + "." + parts[1], // header + "." + payload
	})
	if err != nil {
		return false, "", err
	}
	if signature != parts[2] {
		return false, "", nil
	}

	dataHeader, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, "", err
	}
	var h entity.JwtTokenHeader
	err = json.Unmarshal(dataHeader, &h)
	if err != nil {
		return false, "", err
	}

	dataPayload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, "", err
	}
	var p entity.JwtTokenPayload
	err = json.Unmarshal(dataPayload, &p)
	if err != nil {
		return false, "", err
	}

	// "02.01.2006 15:04:05 UTC-07" template
	jwtDate, err := time.Parse("02.01.2006 15:04:05 UTC-07", h.Exp)
	if err != nil {
		return false, "", err
	}
	dateNow := time.Now()
	if jwtDate.Equal(dateNow) || jwtDate.After(dateNow) {
		return false, "", nil
	}
	return true, p.Id, nil
}
