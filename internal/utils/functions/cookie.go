// Copyright © ivanlobanov. All rights reserved.
package functions

import (
	"net/http"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
)

type SetCookieProps struct {
	W   http.ResponseWriter
	Uid string
}

// SetCookieAndHeaders
// Sets up cookie header and csrf header.
func SetCookieAndHeaders(props SetCookieProps) (http.ResponseWriter, error) {
	expiration := time.Now().Add(myconstants.TimeExpDur)
	jwt, err := NewJwtToken(NewJwtTokenProps{
		UserId: props.Uid,
		Time:   expiration,
	})
	if err != nil {
		return props.W, err
	}
	cookie := http.Cookie{
		Name:     myconstants.JwtCookie,
		Value:    jwt,
		Expires:  expiration,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(props.W, &cookie)

	csrfToken, err := NewCsrfToken(jwt)
	if err != nil {
		return props.W, err
	}
	props.W.Header().Set(myconstants.CsrfHeader, csrfToken)
	return props.W, nil
}

func FlashCookie(w http.ResponseWriter, r *http.Request) {
	sessionCookie := &http.Cookie{
		Name:     myconstants.JwtCookie,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, sessionCookie)

	w.Header().Set(myconstants.CsrfHeader, "")
}
