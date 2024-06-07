package msg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type Handler struct {
	Service
}

func NewMsgHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) AddMessage(c *gin.Context) (int, error) {

	var msgReq CreateMesReq

	if err := c.BindJSON(&msgReq); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := h.Service.AddMessage(c.Request.Context(), &msgReq)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

// if from and to date are not provided then it takes first date of month as from_date and current date as to_date
func (h *Handler) PullAllMessages(c *gin.Context) (int, error) {
	var req GetAllMessageReq
	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := h.Service.PullAllMessages(c.Request.Context(), &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}
