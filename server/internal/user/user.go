package user

import (
	"context"
	"time"
)

type UserOTP struct {
	Uid        int64  `json:"id"`
	OTP        string `json:"otp"`
	Expires_at int64  `json:"expires_at"`
}

type User struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Mobile        int64     `json:"mobile"`
	About         string    `json:"about"`
	Image         string    `json:"image"`
	Last_Seen     time.Time `json:"last_seen"`
	Is_Online     bool      `json:"is_online"`
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type UserLoginReq struct {
	Mobile int64 `json:"mobile"`
}

// after sending login req user gives back otp as res
type UserLoginRes struct {
	ID  string `json:"id"`
	OTP string `json:"otp"`
}

// returns when use authenticated with otp

type CreateUserReq struct {
	Name   string `json:"name"`
	Mobile int64  `json:"mobile"`
	About  string `json:"about"`
	Image  string `json:"image"`
}

type CreateUserRes struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Mobile int64  `json:"mobile"`
	About  string `json:"about"`
	Image  string `json:"image"`
}

type OTPVerificationReq struct {
	ID  int64  `json:"id"`
	Otp string `json:"otp"`
}

type OTPVerificationRes struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Mobile        int64     `json:"mobile"`
	About         string    `json:"about"`
	Image         string    `json:"image"`
	Last_Seen     time.Time `json:"last_seen"`
	Is_Online     bool      `json:"is_online"`
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
}

type UserContactsRes struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Mobile    int64     `json:"mobile"`
	About     string    `json:"about"`
	Image     string    `json:"image"`
	Last_Seen time.Time `json:"last_seen"`
	Is_Online bool      `json:"is_online"`
}

type TokenGenRes struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

// deals with database
type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByMobile(ctx context.Context, mobile int64) (string, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByMobileInt64(ctx context.Context, mobile int64) (int64, error)
	AddUserOTP(ctx context.Context, otp *UserOTP) error
	VerifyOTP(ctx context.Context, otp *UserOTP) (int64, error)
	UpdateTokens(ctx context.Context, token string, refresh_token string, updated_at time.Time, id int64) (int64, error)
	UpdateToken(ctx context.Context, id int64, token string) (int64, error)
}

// act as bridge between repository and handlers
type Service interface {
	CreateUser(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error)
	Login(ctx context.Context, req *UserLoginReq) (*UserLoginRes, error)
	VerifyOTP(ctx context.Context, otp *OTPVerificationReq) (*OTPVerificationRes, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
	UpdateToken(ctx context.Context, id int64, token string) error
}
