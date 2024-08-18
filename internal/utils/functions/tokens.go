package functions

import (
	"encoding/hex"
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
	hash, err := HashWithStatement(HashProps{
		EnvName:   myconstants.EnvCsrfSecret,
		Statement: jwtToken,
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString([]byte(hash)), nil
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
	hEncoded := hex.EncodeToString(rawDataHeader)
	// Encode payload.
	p := entity.JwtTokenPayload{
		Id: props.UserId,
	}
	rawDataPayload, err := easyjson.Marshal(p)
	if err != nil {
		return "", err
	}
	pEncoded := hex.EncodeToString(rawDataPayload)
	// concatenate header and payload
	hpEncoded := hEncoded + "." + pEncoded
	signatureHash, err := HashWithStatement(HashProps{
		EnvName:   myconstants.EnvJwtSecret,
		Statement: hpEncoded,
	})
	if err != nil {
		return "", err
	}
	signature := hex.EncodeToString([]byte(signatureHash))
	return hpEncoded + "." + signature, nil
}
