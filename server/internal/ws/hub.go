package ws

import (
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	WriteMessages chan *msg.Message

	// Inbound group messages
	WriteGroupMessages chan *msg.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		WriteMessages:      make(chan *msg.Message, 5),
		WriteGroupMessages: make(chan *msg.Message, 5),
		Register:           make(chan *Client),
		Unregister:         make(chan *Client),
		clients:            make(map[*Client]bool),
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
		case message := <-h.WriteGroupMessages:

			// get all the connected group clients to ws
			groupClients := findClients(message.Members, h.clients)
			// sending messages to all connected group members
			for client := range h.clients {
				for _, id := range groupClients {
					if client.ID == id {
						select {
						case client.Message <- message:
						default:
							close(client.Message)
							delete(h.clients, client)
						}
					}
				}
			}

			// based on connected group clients send the message
		case message := <-h.WriteMessages:
			for client := range h.clients {
				// send only message when connected client and sender client's requested receiver id is same
				if message.ReceiverID == client.ID || message.SenderID == client.ID {
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

func findClients(clientIDs []int64, clients map[*Client]bool) []int64 {
	var foundIDs []int64
	for _, clientID := range clientIDs {
		if clientExists(clientID, clients) {
			foundIDs = append(foundIDs, clientID)
		}
	}
	return foundIDs
}

func clientExists(clientID int64, clients map[*Client]bool) bool {
	// Find a client with matching ID
	for client := range clients {
		if client.ID == clientID {
			return true
		}
	}
	return false
}
