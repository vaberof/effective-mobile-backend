package apiv1

import (
	"encoding/json"
)

type ResponseStatus string
type ErrorDescription map[string]any

var (
	statusOk    = "Ok"
	statusError = "Error"
)

type Response struct {
	Status  ResponseStatus  `json:"status"`
	Payload json.RawMessage `json:"payload"`
}

type ErrorResponsePayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(payload json.RawMessage) *Response {
	return &Response{
		Status:  ResponseStatus(statusOk),
		Payload: payload,
	}
}

func Error(code string, message string) *Response {
	payload, _ := json.Marshal(&ErrorResponsePayload{
		Code:    code,
		Message: message,
	})

	return &Response{
		Status:  ResponseStatus(statusError),
		Payload: payload,
	}
}
