package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	msg "github.com/mananKoyawala/whatsapp-clone/internal/message"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("origin")
		// return origin == "http://localhost:8080"
	},
}

type Handler struct {
	hub *Hub
}

func NewWsHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) WsConnector(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("uid"))
	userId := int64(id)

	// TODO : check if the client exits in the system do it in the middleware

	// check client already register with id
	for client := range h.hub.clients {
		if client.ID == userId {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already connected",
			})
			return
		}
		break
	}

	log.Println("client id -> ", userId)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	client := &Client{
		conn:    conn,
		Message: make(chan *msg.Message, 10),
		ID:      userId,
	}

	// register client
	h.hub.Register <- client

	go client.readMessage(h.hub)
	go client.writeMessage()

}
