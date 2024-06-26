package contact

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type Handler struct {
	Service
	logger *slog.Logger
}

func NewContactHan(s Service, logger *slog.Logger) Handler {
	return Handler{Service: s, logger: logger}
}

func (h *Handler) AddContact(c *gin.Context) (int, error) {

	var contact CreateContactReq

	if err := c.BindJSON(&contact); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	// TODO : user can request it's own resources

	res, err := h.Service.AddContact(c.Request.Context(), &contact)
	if err != nil {
		h.logger.Error("failed to add contact", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("contact added successfully", slog.String("contactid", helper.Int64ToStirng(res.ID)))
	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) GetContacts(c *gin.Context) (int, error) {
	id, _ := strconv.Atoi(c.Param("id"))

	// TODO : user can request it's own resources

	res, err := h.Service.GetContacts(c.Request.Context(), int64(id))
	if err != nil {
		h.logger.Error("failed to get contacts", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("contacts were got", slog.String("userid", helper.IntToStirng(id)))
	return api.WriteData(c, http.StatusOK, res)
}
