package functions

import (
	"crypto/hmac"
	"crypto/sha256"
	"os"

	"github.com/satori/uuid"
	"github.com/spf13/viper"
)

type HashProps struct {
	EnvName   string
	Statement string
}

// HashWithStatement
// Returns hash that is transmitted in the client-server model by custom header.
func HashWithStatement(props HashProps) (string, error) {
	secretKey := getSecretToken(props.EnvName)
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(props.Statement))
	if err != nil {
		return "", err
	}
	return string(mac.Sum(nil)), nil
}

// getSecretToken
// Returns token from programm environment.
func getSecretToken(envName string) (csrfSecretValue string) {
	csrfSecret := viper.Get(envName)
	if csrfSecret == nil {
		csrfSecretValue = uuid.NewV4().String()
		os.Setenv(envName, csrfSecretValue)
	} else {
		csrfSecretValue = csrfSecret.(string)
	}
	return csrfSecretValue
}
