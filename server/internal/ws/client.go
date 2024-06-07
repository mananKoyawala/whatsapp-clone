package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

var (
	pongwait = 10 * time.Second

	pingInterval = (pongwait * 9) / 10
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

	if err := c.conn.SetReadDeadline(time.Now().Add(pongwait)); err != nil {
		log.Println(err)
		return
	}

	c.conn.SetReadLimit(512) // called jambo framming

	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// means client closed the connection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error :- %v", err)
			}
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

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case msg, ok := <-c.Message:
			if !ok {
				return
			}
			c.conn.WriteJSON(msg)

		case <-ticker.C:
			// log.Println("ping")

			// send a ping to the clinet
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				log.Println("writemsg err :-", err.Error())
				return
			}
		}
	}

}

func (c *Client) pongHandler(pongMsg string) error {
	// log.Println("pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongwait))
}
