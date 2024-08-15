package functions

import (
	"net/http"
	"strconv"

	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/mailru/easyjson"
)

type ErrorResponseProps struct {
	W          http.ResponseWriter
	Msg        string
	CodeStatus int
}

type JsonResponseProps struct {
	W          http.ResponseWriter
	Payload    easyjson.Marshaler
	CodeStatus int
}

// ErrorResponse
// Uses for error response from the server. For marshalling uses easyjson (mailru).
func ErrorResponse(props ErrorResponseProps) {
	props.W.Header().Add("Content-Type", "application/json")
	errObject := &entity.ErrorDetail{Error: props.Msg}
	body, err := easyjson.Marshal(errObject)
	if err != nil {
		props.W.Header().Add("Content-Length", "0")
		props.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	props.W.WriteHeader(props.CodeStatus)
	contentLength, err := props.W.Write(body)
	if err != nil {
		props.W.WriteHeader(http.StatusInternalServerError)
	}
	props.W.Header().Add("Content-Length", strconv.Itoa(contentLength))
}

// JsonResponse
// Forms server response in JSON format. For marshalling uses easyjson (mailru).
func JsonResponse(props JsonResponseProps) {
	props.W.Header().Add("Content-Type", "application/json")
	body, err := easyjson.Marshal(props.Payload)
	if err != nil {
		props.W.Header().Add("Content-Length", "0")
		props.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	props.W.WriteHeader(props.CodeStatus)
	contentLength, err := props.W.Write(body)
	if err != nil {
		props.W.WriteHeader(http.StatusInternalServerError)
	}
	props.W.Header().Add("Content-Length", strconv.Itoa(contentLength))
}
