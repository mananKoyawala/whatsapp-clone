package ws

import (
	"context"
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

func (c *Client) readMessage(hub *Hub, mr msg.Repository) {
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

		var nmsg *msg.Message

		// unmarshal the message form the client
		if err := json.Unmarshal(message, &nmsg); err != nil {
			log.Println(err.Error())
			break
		}

		newMessage := msg.NewMessage(nmsg)
		resmessage, err := mr.AddMessage(context.Background(), *newMessage)
		if err != nil {
			log.Println(err.Error())
		} // TODO : we need to make sure that both sender and receiver must be exists
		nmsg.ID = resmessage.ID

		// sending the message if client and sender client id is same
		if nmsg.SenderID == c.ID {
			hub.WriteMessages <- nmsg
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
