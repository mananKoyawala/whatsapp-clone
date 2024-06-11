package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
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

	str, ok := validateCreateUser(userReq)
	if !ok {
		return api.WriteMessage(c, http.StatusBadRequest, str)
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

	if helper.CheckLength(int(userReq.Mobile), 10) {
		return api.WriteMessage(c, http.StatusBadRequest, " / mobile number must be 10 digit / ")
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
