package health

import (
	"context"
	"log"
	"net/http"
)

type Pinger interface {
	Ping(context.Context) error
}

type Handler struct {
	database Pinger
}

func NewHandler(database Pinger) *Handler {
	return &Handler{
		database: database,
	}
}

func (handler *Handler) Live(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	responseWriter.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	responseWriter.WriteHeader(http.StatusOK)

	_, err := responseWriter.Write(
		[]byte(`{"status":"ok","service":"samatalk-api"}`),
	)

	if err != nil {
		log.Printf(
			"could not write health response: %v",
			err,
		)
	}
}

func (handler *Handler) Ready(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	err := handler.database.Ping(
		request.Context(),
	)

	if err != nil {
		responseWriter.Header().Set(
			"Content-Type",
			"application/json; charset=utf-8",
		)

		responseWriter.WriteHeader(
			http.StatusServiceUnavailable,
		)

		_, _ = responseWriter.Write(
			[]byte(`{"status":"not_ready","database":"down"}`),
		)

		return
	}

	responseWriter.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	responseWriter.WriteHeader(http.StatusOK)

	_, _ = responseWriter.Write(
		[]byte(`{"status":"ready","database":"up"}`),
	)
}
