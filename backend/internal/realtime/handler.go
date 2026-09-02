package realtime

import (
	"net/http"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(request *http.Request) bool {
				origin := request.Header.Get("Origin")

				return origin == "http://localhost:5173" ||
					origin == "http://127.0.0.1:5173"
			},
		},
	}
}

func (handler *Handler) Connect(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(request.Context())
	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	connection, err := handler.upgrader.Upgrade(
		responseWriter,
		request,
		nil,
	)
	if err != nil {
		return
	}

	client := NewClient(
		handler.hub,
		connection,
		currentUser.ID,
	)

	handler.hub.Register(client)

	go client.writePump()

	client.readPump()
}
