package user

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
	layer  string
}

func NewUserHandler(s Service, logger *slog.Logger) *Handler {
	return &Handler{Service: s, logger: logger, layer: "userHandler"}
}

func (h *Handler) CreateUser(c *gin.Context) (int, error) {
	var userReq CreateUserReq

	if err := c.BindJSON(&userReq); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	str, ok := validateCreateUser(userReq)
	if !ok {
		h.logger.Warn("failed to validate details", slog.String("validation_error", str))
		return api.WriteMessage(c, http.StatusBadRequest, str)
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &userReq)
	if err != nil {
		h.logger.Error("failed to created user", slog.String("error", err.Error()))
		return http.StatusInternalServerError, err
	}

	h.logger.Info("user created successfully", slog.String("userid", helper.Int64ToStirng(res.ID)))
	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) LoginUser(c *gin.Context) (int, error) {
	var userReq UserLoginReq
	if err := c.BindJSON(&userReq); err != nil {
		h.logger.Error("failed to bind JSON", slog.String("error", err.Error()))
		return http.StatusBadRequest, err
	}

	if helper.CheckLength(int(userReq.Mobile), 10) {
		h.logger.Warn("failed to validate mobile", slog.String("validation_error", "mobile number must be 10"))
		return api.WriteMessage(c, http.StatusBadRequest, " / mobile number must be 10 digit / ")
	}

	res, err := h.Service.Login(c.Request.Context(), &userReq)
	if err != nil {
		h.logger.Error("login failed", slog.String("error", err.Error()))
		return http.StatusNotFound, err
	}

	h.logger.Info("user logged in successfully", slog.String("userid", res.ID))
	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) VerifyUserOTP(c *gin.Context) (int, error) {
	id, _ := strconv.Atoi(c.Query("uid"))
	otp := c.Query("otp")

	if id <= 0 || otp == "" {
		msg := fmt.Sprintf("userid %d and otp %s", id, otp)
		h.logger.Warn("unauthorized access attempts", slog.String("validation_error", msg))
		return http.StatusUnauthorized, errors.New("unauthorized access")
	}

	req := &OTPVerificationReq{
		ID:  int64(id),
		Otp: otp,
	}

	res, err := h.Service.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("otp verification failed", slog.String("error", err.Error()))
		return http.StatusUnauthorized, errors.New("unauthorized access")
	}

	h.logger.Info("otp verification successfull", slog.String("userid", strconv.Itoa(id)))
	return api.WriteData(c, http.StatusOK, res)
}

func validateCreateUser(req CreateUserReq) (string, bool) {
	var str string

	if req.Name == "" {
		str += " / name required / "
	}

	if helper.CheckLength(int(req.Mobile), 10) {
		str += " / mobile number must be 10 digit / "
	}

	if req.About == "" {
		str += " / about required / "
	}

	if req.Image == "" {
		str += " / image required / "
	}

	if str != "" {
		return str, false
	}

	return "", true
}
