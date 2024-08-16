package functions

import (
	"encoding/base64"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/mailru/easyjson"
)

type NewJwtTokenProps struct {
	UserId string
	Time   time.Time
}

// NewCsrfToken
// Generates csrf-token.
func NewCsrfToken(jwtToken string) (string, error) {
	return HashWithStatement(HashProps{
		EnvName:   myconstants.EnvCsrfSecret,
		Statement: jwtToken,
	})
}

// NewCsrfToken
// Generates jwt-token.
func NewJwtToken(props NewJwtTokenProps) (string, error) {
	// Encode header.
	h := entity.JwtTokenHeader{
		Exp: props.Time.Format("02.01.2006 15:04:05 UTC-07"),
	}
	rawDataHeader, err := easyjson.Marshal(h)
	if err != nil {
		return "", err
	}
	hEncoded := base64.StdEncoding.EncodeToString(rawDataHeader)
	// Encode payload.
	p := entity.JwtTokenPayload{
		Id: props.UserId,
	}
	rawDataPayload, err := easyjson.Marshal(p)
	if err != nil {
		return "", err
	}
	pEncoded := base64.StdEncoding.EncodeToString([]byte(rawDataPayload))
	// concatenate header and payload
	hpEncoded := hEncoded + "." + pEncoded
	signature, err := HashWithStatement(HashProps{
		EnvName:   myconstants.EnvJwtSecret,
		Statement: hpEncoded,
	})
	if err != nil {
		return "", err
	}
	return hpEncoded + "." + signature, nil
}
