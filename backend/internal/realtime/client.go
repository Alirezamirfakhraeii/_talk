package realtime

import "github.com/gorilla/websocket"

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID int64
	Send   chan []byte
}

func NewClient(
	hub *Hub,
	connection *websocket.Conn,
	userID int64,
) *Client {
	return &Client{
		Hub:    hub,
		Conn:   connection,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (client *Client) writePump() {
	defer client.Conn.Close()

	for {
		message, ok := <-client.Send

		if !ok {
			_ = client.Conn.WriteMessage(
				websocket.CloseMessage,
				[]byte{},
			)
			return
		}

		err := client.Conn.WriteMessage(
			websocket.TextMessage,
			message,
		)
		if err != nil {
			return
		}
	}
}
