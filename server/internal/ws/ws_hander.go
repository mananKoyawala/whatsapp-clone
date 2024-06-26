package ws

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
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
	hub    *Hub
	mr     msg.Repository
	gr     group.Repository
	logger *slog.Logger
}

func NewWsHandler(h *Hub, mr msg.Repository, gr group.Repository, logger *slog.Logger) *Handler {
	return &Handler{
		hub:    h,
		mr:     mr,
		gr:     gr,
		logger: logger,
	}
}

func (h *Handler) WsConnector(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("uid"))
	userId := int64(id)

	// check client already register with id
	for client := range h.hub.clients {
		if client.ID == userId {
			h.logger.Error("user already connected")
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already connected",
			})
			return
		}
		break
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("client upgration failed", slog.String("error", err.Error()))
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
	h.logger.Debug("client registed", slog.String("clientid", helper.Int64ToStirng(userId)))
	h.hub.Register <- client

	// injecting logger
	go client.readMessage(h.hub, h.mr, h.gr, h.logger)
	go client.writeMessage(h.logger)
}
