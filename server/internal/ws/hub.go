package ws

import (
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	WriteMessages chan *msg.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		WriteMessages: make(chan *msg.Message, 5),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Message)
			}
		case message := <-h.WriteMessages:
			for client := range h.clients {
				// send only message when connected client and sender client's requested receiver id is same
				if message.ReceiverID == client.ID {
					select {
					case client.Message <- message:
					default:
						close(client.Message)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
