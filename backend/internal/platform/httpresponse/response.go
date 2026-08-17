package httpresponse

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func Success(
	responseWriter http.ResponseWriter,
	statusCode int,
	message string,
	data any,
) {
	WriteJSON(
		responseWriter,
		statusCode,
		Response{
			Success: true,
			Message: message,
			Data:    data,
		},
	)
}

func Error(
	responseWriter http.ResponseWriter,
	statusCode int,
	code string,
	message string,
	fields map[string]string,
) {
	WriteJSON(
		responseWriter,
		statusCode,
		Response{
			Success: false,
			Error: &ErrorDetail{
				Code:    code,
				Message: message,
				Fields:  fields,
			},
		},
	)
}

func WriteJSON(
	responseWriter http.ResponseWriter,
	statusCode int,
	data any,
) {
	responseWriter.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	responseWriter.WriteHeader(statusCode)

	_ = json.NewEncoder(
		responseWriter,
	).Encode(data)
}
