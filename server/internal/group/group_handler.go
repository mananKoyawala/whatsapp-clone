package group

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type Handler struct {
	Service
}

func NewGroupHandler(s Service) Handler {
	return Handler{Service: s}
}

func (h *Handler) CreateGroup(c *gin.Context) (int, error) {
	var groupReq CreateGroupReq

	if err := c.BindJSON(&groupReq); err != nil {
		return http.StatusBadRequest, err
	}

	req := NewGroup(groupReq)

	groupRes, err := h.Service.CreateGroup(c.Request.Context(), req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, groupRes)
}

func (h *Handler) AddMemberToGroup(c *gin.Context) (int, error) {
	var req AddMemberReq

	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	err := h.Service.AddMemberToGroup(c.Request.Context(), &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteMessage(c, http.StatusOK, "members added into the group")
}

func (h *Handler) GetAllGroupByUserID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("uid"))
	if err != nil || id <= 0 {
		return http.StatusBadRequest, errors.New("invalid id")
	}

	groups, err := h.Service.GetAllGroupByUserID(c.Request.Context(), int64(id))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, groups)
}

func (h *Handler) RemoveMemberFromGroup(c *gin.Context) (int, error) {
	gid, _ := strconv.Atoi(c.Param("gid"))
	uid, _ := strconv.Atoi(c.Param("uid"))

	if err := h.Service.RemoveMemberFromGroup(c.Request.Context(), int64(gid), int64(uid)); err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, err
	}
	return api.WriteMessage(c, http.StatusOK, "user removed from group.")
}

func (h *Handler) GetGroupDetailsByID(c *gin.Context) (int, error) {
	groupID, _ := strconv.Atoi(c.Param("gid"))

	res, err := h.Service.GetGroupDetailsByID(c.Request.Context(), int64(groupID))
	if err != nil {
		return http.StatusNotFound, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) UpdateGroupDetails(c *gin.Context) (int, error) {
	var req UpdateGroup

	if err := c.BindJSON(&req); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := h.Service.UpdateGroupDetails(c.Request.Context(), req)
	if err != nil {
		return http.StatusOK, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) DeleteGroupByID(c *gin.Context) (int, error) {
	groupId, _ := strconv.Atoi(c.Param("gid"))

	if err := h.Service.DeleteGroupByID(c.Request.Context(), int64(groupId)); err != nil {
		return http.StatusInternalServerError, err
	}
	msg := fmt.Sprintf("group deleted with id %d", groupId)
	return api.WriteMessage(c, http.StatusOK, msg)
}
