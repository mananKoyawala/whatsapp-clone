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

	msg, ok := validateMessage(msgReq)
	if !ok {
		return api.WriteMessage(c, http.StatusBadRequest, msg)
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

	// validating json
	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	// check wheather user can access resourse or not
	/* it prevents from situation like where hacker knows two or more peopels id and token and want to manuplating data but he can't manuplate data of x user using y's token */
	if ok := checkRequestUserAuthenticated(c, req.SenderID); !ok {
		return http.StatusUnauthorized, api.Unauthorized
	}

	// pull all messages
	res, err := h.Service.PullAllMessages(c.Request.Context(), &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) UpdateIsReadMessage(c *gin.Context) (int, error) {
	var req struct {
		Message []MessageReq `json:"msg"`
	}

	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	if err := h.Service.UpdateIsReadMessage(c.Request.Context(), &req.Message); err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteMessage(c, http.StatusOK, "all the messages are read updated")
}

func (h *Handler) DeleteMessage(c *gin.Context) (int, error) {
	var req MessageReq
	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	// check wheather user can access resourse or not
	if ok := checkRequestUserAuthenticated(c, req.SenderID); !ok {
		return http.StatusUnauthorized, api.Unauthorized
	}

	if err := h.Service.DeleteMessage(c.Request.Context(), &req); err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteMessage(c, http.StatusOK, "message deleted")
}

func validateMessage(req CreateMesReq) (string, bool) {
	var str string

	// * if message_type="text" then message_text != ""
	// * if message_type="media" then media_url != ""

	if req.SenderID <= 0 {
		str += " / sender_id required / "
	}

	if req.ReceiverID <= 0 {
		str += " / receiver_id required / "
	}

	if !(req.MessageType == "text" || req.MessageType == "media") {
		str += " / message_type must be text or media / "
	}

	if req.MessageType == "text" {
		if req.MessageText == "" {
			str += " / message_text required / "
		}
	}

	if req.MessageType == "media" {
		if req.MediaUrl == "" {
			str += " / media_url required / "
		}
	}

	if str != "" {
		return str, false
	}

	return "", true
}

func checkRequestUserAuthenticated(c *gin.Context, senderId int64) bool {
	reqUserId, ok := c.Get("id")
	if !ok {
		return false
	}

	return reqUserId == senderId
}

func (h *Handler) PullAllGroupMessages(c *gin.Context) (int, error) {
	var req GetAllGroupMessageReq

	// validating json
	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	// pull all messages
	res, err := h.Service.PullAllGroupMessages(c.Request.Context(), &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) DeleteGroupMessage(c *gin.Context) (int, error) {
	var req MessageGroupReq
	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	if err := h.Service.DeleteGroupMessage(c.Request.Context(), &req); err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteMessage(c, http.StatusOK, "message deleted")
}
