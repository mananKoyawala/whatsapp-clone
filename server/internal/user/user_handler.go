package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type Handler struct {
	Service
}

func NewUserHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateUser(c *gin.Context) (int, error) {
	var userReq CreateUserReq

	if err := c.BindJSON(&userReq); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &userReq)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) LoginUser(c *gin.Context) (int, error) {
	var userReq UserLoginReq
	if err := c.BindJSON(&userReq); err != nil {
		return http.StatusBadRequest, err
	}

	res, err := h.Service.Login(c.Request.Context(), &userReq)
	if err != nil {
		return http.StatusNotFound, err
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *Handler) VerifyUserOTP(c *gin.Context) (int, error) {
	id, _ := strconv.Atoi(c.Query("uid"))
	otp := c.Query("otp")

	if id <= 0 || otp == "" {
		return http.StatusUnauthorized, errors.New("unauthorized access")
	}

	req := &OTPVerificationReq{
		ID:  int64(id),
		Otp: otp,
	}

	res, err := h.Service.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		return http.StatusUnauthorized, errors.New("unauthorized access")
	}

	return api.WriteData(c, http.StatusOK, res)
}
