package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

var (
	pongwait = 10 * time.Second

	pingInterval = (pongwait * 9) / 10
)

type Client struct {
	conn         *websocket.Conn
	Message      chan *msg.Message
	GroupMessage chan *msg.Message
	ID           int64 `json:"id"`
	GroupID      int64 `json:"group_id"`
}

func (c *Client) readMessage(hub *Hub, mr msg.Repository, gr group.Repository, logger *slog.Logger) {
	defer func() {
		hub.Unregister <- c
		c.conn.Close()
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongwait)); err != nil {
		logger.Error("read dead line exceed", slog.String("error", err.Error()))
		return
	}

	c.conn.SetReadLimit(512) // called jambo framming

	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// means client closed the connection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				msg := fmt.Sprintf("error :- %v", err)
				logger.Debug(msg)
			}
			break
		}

		var nmsg *msg.Message

		// unmarshal the message form the client
		if err := json.Unmarshal(message, &nmsg); err != nil {
			logger.Error("failed to bind JSON", slog.String("error", err.Error()))
			break
		}

		// hub.WriteMessages <- newMessage
		// get all the group members based on the group id and send them

		newMessage := msg.NewMessage(nmsg) // get new message

		// check message is type of group
		if newMessage.IsGroupMessage {
			// group message
			sendGroupMessage(newMessage, hub, mr, gr, logger)
		} else {
			// one-one message
			sendOneOneMessage(*nmsg, newMessage, c, hub, mr, logger)
		}

	}
}

func sendOneOneMessage(nmsg msg.Message, newMessage *msg.Message, c *Client, hub *Hub, mr msg.Repository, logger *slog.Logger) {
	newMessage.GroupID = 0 // it's make easy when we retrive the messages based on group id , it's prevent from getting other messages
	resmessage, err := mr.AddMessage(context.Background(), *newMessage)
	if err == nil {
		// if both sender and receiver exist
		newMessage.ID = resmessage.ID
		// sending the message if client and sender client id is same
		if nmsg.SenderID == c.ID {
			hub.WriteMessages <- newMessage
		}
	} else {
		logger.Error("failed to sent message", slog.String("error", err.Error()))
	}
}

func sendGroupMessage(newMessage *msg.Message, hub *Hub, mr msg.Repository, gr group.Repository, logger *slog.Logger) {
	_, err := gr.GetGroupByID(context.Background(), newMessage.GroupID)
	if err == nil {
		newMessage.ReceiverID = newMessage.SenderID // only groups chat has receiver id = senderid
		resmessage, err := mr.AddMessage(context.Background(), *newMessage)
		if err == nil {
			newMessage.ID = resmessage.ID

			// get all the group memebers
			members, err := gr.GetMemberByGroupID(context.Background(), newMessage.GroupID)
			if err != nil {
				logger.Error("error occurs while getting the members of group", slog.String("groupid", helper.Int64ToStirng(newMessage.GroupID)))
			}
			// var msg msg.GroupMessage
			newMessage.Members = members

			// send message to group members
			hub.WriteGroupMessages <- newMessage
		} else {
			logger.Error("failed to sent message", slog.String("error", err.Error()))
		}
	}
}

func (c *Client) writeMessage(logger *slog.Logger) {
	defer c.conn.Close()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case msg, ok := <-c.Message:
			if !ok {
				logger.Error("failed to write message")
				return
			}
			c.conn.WriteJSON(msg)

		case <-ticker.C:
			// log.Println("ping")

			// send a ping to the clinet
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				logger.Error("failed to write ping message", slog.String("error", err.Error()))
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	// log.Println("pong")
	return c.conn.SetReadDeadline(time.Now().Add(pongwait))
}
