package group

import (
	"errors"
	"fmt"
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

func NewGroupHandler(s Service, logger *slog.Logger) Handler {
	return Handler{Service: s, logger: logger}
}

func (h *Handler) CreateGroup(c *gin.Context) (int, error) {
	var groupReq CreateGroupReq

	if err := c.BindJSON(&groupReq); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	req := NewGroup(groupReq)

	groupRes, err := h.Service.CreateGroup(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to create group", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("group was created", slog.String("groupid", helper.Int64ToStirng(groupRes.ID)))
	return api.WriteData(c, http.StatusOK, groupRes)
}

func (h *Handler) AddMemberToGroup(c *gin.Context) (int, error) {
	var req AddMemberReq

	if err := c.BindJSON(&req); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	err := h.Service.AddMemberToGroup(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to add member to group", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("memebers were added into the group", slog.String("groupid", helper.Int64ToStirng(req.GroupID)))
	return api.WriteMessage(c, http.StatusOK, "members were added into the group")
}

func (h *Handler) GetAllGroupByUserID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("uid"))
	if err != nil || id <= 0 {
		h.logger.Warn("requested id was wrong", slog.String("userid", helper.IntToStirng(id)))
		return http.StatusBadRequest, errors.New("permission denied")
	}

	groups, err := h.Service.GetAllGroupByUserID(c.Request.Context(), int64(id))
	if err != nil {
		h.logger.Error("failed to get all group user by id", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("all group members were got by userid", slog.String("userid", helper.IntToStirng(id)))
	return api.WriteData(c, http.StatusOK, groups)
}

func (h *Handler) RemoveMemberFromGroup(c *gin.Context) (int, error) {
	gid, _ := strconv.Atoi(c.Param("gid"))
	uid, _ := strconv.Atoi(c.Param("uid"))

	if err := h.Service.RemoveMemberFromGroup(c.Request.Context(), int64(gid), int64(uid)); err != nil {
		h.logger.Error("failed to remove memebers form group", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("member was removed from group", slog.String("userid", helper.IntToStirng(uid)), slog.String("groupid", helper.IntToStirng(gid)))
	return api.WriteMessage(c, http.StatusOK, "user removed from group.")
}

func (h *Handler) GetGroupDetailsByID(c *gin.Context) (int, error) {
	groupID, _ := strconv.Atoi(c.Param("gid"))

	res, err := h.Service.GetGroupDetailsByID(c.Request.Context(), int64(groupID))
	if err != nil {
		h.logger.Error("failed to get group details by id", slog.String("error", err.Error()))
		return http.StatusNotFound, err
	}

	h.logger.Info("got group details by id", slog.String("groupid", helper.IntToStirng(groupID)))
	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) UpdateGroupDetails(c *gin.Context) (int, error) {
	var req UpdateGroup

	if err := c.BindJSON(&req); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	res, err := h.Service.UpdateGroupDetails(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to update group details", slog.String("error", err.Error()))
		return http.StatusOK, err
	}

	h.logger.Info("updated group details", slog.String("groupid", helper.Int64ToStirng(req.ID)))
	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) DeleteGroupByID(c *gin.Context) (int, error) {
	groupId, _ := strconv.Atoi(c.Param("gid"))

	if err := h.Service.DeleteGroupByID(c.Request.Context(), int64(groupId)); err != nil {
		h.logger.Error("failed to delete group by id", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	msg := fmt.Sprintf("group deleted with id %d", groupId)
	h.logger.Info("group was deleted", slog.String("groupid", helper.IntToStirng(groupId)))
	return api.WriteMessage(c, http.StatusOK, msg)
}
