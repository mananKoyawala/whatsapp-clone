package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"time"

	"github.com/google/uuid"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type service struct {
	Repository
	timeout time.Duration
	logger  *slog.Logger
	layer   string
}

func NewUserService(repository Repository, logger *slog.Logger) Service {
	return &service{Repository: repository, timeout: time.Duration(100) * time.Second, logger: logger, layer: "userService"}
}

func (s *service) CreateUser(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id, _ := s.Repository.GetUserByMobile(ctx, user.Mobile)
	if id != "" {
		msg := fmt.Sprintf("user mobile number is %d", user.Mobile)
		s.logger.Error("user already exists", slog.String("error", msg))
		return nil, errors.New("user already register with mobile number")
	}

	current_time, _ := helper.GetTime()
	u := &User{
		Name:          user.Name,
		Mobile:        user.Mobile,
		About:         user.About,
		Image:         user.Image,
		Created_at:    current_time,
		Updated_at:    current_time,
		Token:         "",
		Refresh_Token: "",
		Last_Seen:     current_time,
		Is_Online:     false,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		s.logger.Error("failed to create user", slog.String("error", err.Error()))
		return nil, err
	}

	res := &CreateUserRes{
		ID:     int64(r.ID),
		Name:   r.Name,
		Mobile: r.Mobile,
		About:  r.About,
		Image:  r.Image,
	}

	s.logger.Info("user created successfully", slog.String("userid", helper.Int64ToStirng(res.ID)))
	return res, nil
}

func (s *service) Login(ctx context.Context, req *UserLoginReq) (*UserLoginRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	userId, err := s.Repository.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		s.logger.Error("failed to get user by mobile", slog.String("error", err.Error()))
		return nil, err
	}

	o := uuid.New().String()
	id, _ := helper.StirngToInt(userId)
	expiry := time.Now().Add(15 * time.Second).Local().Unix()
	otp := &UserOTP{
		Uid:        int64(id),
		OTP:        o,
		Expires_at: expiry,
	}
	if err = s.Repository.AddUserOTP(ctx, otp); err != nil {
		s.logger.Error("error while adding otp", slog.String("error", err.Error()))
		return nil, errors.New("error while adding otp")
	}
	res := &UserLoginRes{
		ID:  userId,
		OTP: o,
	}

	s.logger.Info("user login successful", slog.String("userid", res.ID))
	return res, nil
}

func (s *service) VerifyOTP(ctx context.Context, o *OTPVerificationReq) (*OTPVerificationRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	otp := &UserOTP{
		Uid: o.ID,
		OTP: o.Otp,
	}
	expiry, err := s.Repository.VerifyOTP(ctx, otp)
	if err != nil {
		s.logger.Error("otp verification failed", slog.String("error", err.Error()))
		return nil, err
	}
	// log.Println(time.Now())
	// log.Println(expiry)
	t := time.Now().Local().Unix()
	if expiry < t {
		s.logger.Error("otp expired")
		return nil, errors.New("otp expires")
	}

	// generating tokens
	token, refresh_token, err := helper.GenerateJwtToken(o.ID)
	if err != nil {
		s.logger.Error("error occurs while generating tokens", slog.String("error", err.Error()))
		log.Println("error occurs while generating tokens")
	}

	updated_at := time.Now()

	// updates tokens
	_, err = s.Repository.UpdateTokens(ctx, token, refresh_token, updated_at, o.ID)
	if err != nil {
		s.logger.Error("failed to updated token", slog.String("error", err.Error()))
		return nil, err
	}

	// get all the data by id
	// log.Println(o.ID)
	user, err := s.Repository.GetUserById(ctx, o.ID)
	if err != nil {
		s.logger.Error("failed to get user by id", slog.String("error", err.Error()))
		return nil, err
	}

	// prepare res for user to send
	res := &OTPVerificationRes{
		ID:            o.ID,
		Name:          user.Name,
		About:         user.About,
		Image:         user.Image,
		Last_Seen:     user.Last_Seen,
		Mobile:        user.Mobile,
		Is_Online:     user.Is_Online,
		Token:         user.Token,
		Refresh_Token: user.Refresh_Token,
	}

	s.logger.Info("otp verification successful", slog.String("userid", helper.Int64ToStirng(res.ID)))
	return res, nil
}

func (s *service) GetUserById(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.Repository.GetUserById(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user by id", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("user retrieved by id", slog.String("userid", helper.Int64ToStirng(res.ID)))
	return res, err
}
