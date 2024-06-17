package contact

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type Handler struct {
	Service
}

func NewContactHan(s Service) Handler {
	return Handler{Service: s}
}

func (h *Handler) AddContact(c *gin.Context) (int, error) {

	var contact CreateContactReq

	if err := c.BindJSON(&contact); err != nil {
		return http.StatusBadRequest, err
	}

	// TODO : user can request it's own resources

	res, err := h.Service.AddContact(c.Request.Context(), &contact)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) GetContacts(c *gin.Context) (int, error) {
	id, _ := strconv.Atoi(c.Param("id"))

	// TODO : user can request it's own resources

	res, err := h.Service.GetContacts(c.Request.Context(), int64(id))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}
