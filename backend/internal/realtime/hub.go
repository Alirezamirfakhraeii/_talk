package realtime

import "context"

type Delivery struct {
	UserID  int64
	Payload []byte
}

type Hub struct {
	clients    map[int64]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	deliver    chan Delivery
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		deliver:    make(chan Delivery),
	}
}

func (hub *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			hub.closeAllClients()
			return

		case client := <-hub.register:
			hub.addClient(client)

		case client := <-hub.unregister:
			hub.removeClient(client)

		case delivery := <-hub.deliver:
			hub.deliverToUser(delivery)
		}
	}
}

func (hub *Hub) Register(client *Client) {
	hub.register <- client
}

func (hub *Hub) Unregister(client *Client) {
	hub.unregister <- client
}

func (hub *Hub) Deliver(userID int64, payload []byte) {
	hub.deliver <- Delivery{
		UserID:  userID,
		Payload: payload,
	}
}

func (hub *Hub) addClient(client *Client) {
	if _, exists := hub.clients[client.UserID]; !exists {
		hub.clients[client.UserID] = make(map[*Client]struct{})
	}

	hub.clients[client.UserID][client] = struct{}{}
}

func (hub *Hub) removeClient(client *Client) {
	userClients, exists := hub.clients[client.UserID]
	if !exists {
		return
	}

	if _, exists := userClients[client]; !exists {
		return
	}

	delete(userClients, client)
	close(client.Send)

	if len(userClients) == 0 {
		delete(hub.clients, client.UserID)
	}
}

func (hub *Hub) deliverToUser(delivery Delivery) {
	userClients, exists := hub.clients[delivery.UserID]
	if !exists {
		return
	}

	for client := range userClients {
		select {
		case client.Send <- delivery.Payload:
		default:
		}
	}
}

func (hub *Hub) closeAllClients() {
	for _, userClients := range hub.clients {
		for client := range userClients {
			close(client.Send)
		}
	}

	hub.clients = make(map[int64]map[*Client]struct{})
}
