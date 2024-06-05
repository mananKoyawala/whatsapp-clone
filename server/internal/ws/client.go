package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

type Client struct {
	conn    *websocket.Conn
	Message chan *msg.Message
	ID      int64 `json:"id"`
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error :- %v", err)
			}
			// means client closed the connection
			break
		}

		log.Println(string(message))

		var msg *msg.Message

		// unmarshal the message form the client
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println(err.Error())
			break
		}

		// TODO : add message to database

		// sendign the message if client and sender client id is same
		if msg.SenderID == c.ID {
			hub.WriteMessages <- msg
		}

	}
}

func (c *Client) writeMessage() {
	defer c.conn.Close()

	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}

		c.conn.WriteJSON(msg)
	}
}
