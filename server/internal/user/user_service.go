package user

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewUserService(repository Repository) Service {
	return &service{Repository: repository, timeout: time.Duration(100) * time.Second}
}

func (s *service) CreateUser(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id, _ := s.Repository.GetUserByMobile(ctx, user.Mobile)
	if id != "" {
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
		return nil, err
	}

	res := &CreateUserRes{
		ID:     int64(r.ID),
		Name:   r.Name,
		Mobile: r.Mobile,
		About:  r.About,
		Image:  r.Image,
	}

	return res, nil
}

func (s *service) Login(ctx context.Context, req *UserLoginReq) (*UserLoginRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	userId, err := s.Repository.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}

	o := uuid.New().String()
	id, _ := strconv.Atoi(userId)
	expiry := time.Now().Add(15 * time.Second).Local().Unix()
	otp := &UserOTP{
		Uid:        int64(id),
		OTP:        o,
		Expires_at: expiry,
	}
	if err = s.Repository.AddUserOTP(ctx, otp); err != nil {
		log.Println(err.Error())
		return nil, errors.New("error while adding otp")
	}
	res := &UserLoginRes{
		ID:  userId,
		OTP: o,
	}

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
		return nil, err
	}
	log.Println(time.Now())
	log.Println(expiry)
	t := time.Now().Local().Unix()
	if expiry < t {
		return nil, errors.New("otp expires")
	}

	// generating tokens
	token, refresh_token, err := helper.GenerateJwtToken(o.ID)
	if err != nil {
		log.Println("error occurs while generating tokens")
	}

	updated_at := time.Now()

	// updates tokens
	_, err = s.Repository.UpdateTokens(ctx, token, refresh_token, updated_at, o.ID)
	if err != nil {
		return nil, err
	}

	// get all the data by id
	log.Println(o.ID)
	user, err := s.Repository.GetUserById(ctx, o.ID)
	if err != nil {
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

	return res, nil
}
