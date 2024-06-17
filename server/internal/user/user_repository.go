package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int

	query := "INSERT INTO users (name,mobile,about,image,is_online,token,refresh_token,created_at,updated_at,last_seen) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Mobile, user.About, user.Image, user.Is_Online, user.Token, user.Refresh_Token, user.Created_at, user.Updated_at, user.Last_Seen).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByMobile(ctx context.Context, mobile int64) (string, error) {
	var userId int
	query := "SELECT id FROM users WHERE mobile = $1"
	if err := r.db.QueryRowContext(ctx, query, mobile).Scan(&userId); err != nil {
		return "", err
	}
	return strconv.Itoa(userId), nil
}

func (r *repository) GetUserByMobileInt64(ctx context.Context, mobile int64) (int64, error) {
	var userId int
	query := "SELECT id FROM users WHERE mobile = $1"
	if err := r.db.QueryRowContext(ctx, query, mobile).Scan(&userId); err != nil {
		return 0, err
	}
	return int64(userId), nil
}

func (r *repository) GetUserById(ctx context.Context, id int64) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = $1"
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Mobile, &user.About, &user.Image, &user.Last_Seen, &user.Is_Online, &user.Token, &user.Refresh_Token, &user.Created_at, &user.Updated_at); err != nil {
		log.Println(err.Error())
		return nil, errors.New("can't fetch the user")
	}
	return &user, nil
}

func (r *repository) AddUserOTP(ctx context.Context, otp *UserOTP) error {
	query := "INSERT INTO otps (id,otp,expires_at) VALUES ($1,$2,$3)"
	if err := r.db.QueryRowContext(ctx, query, otp.Uid, otp.OTP, otp.Expires_at); err.Err() != nil {
		log.Println("is nil", err.Err())
		return errors.New("can't add otp")
	}
	return nil
}

func (r *repository) VerifyOTP(ctx context.Context, otp *UserOTP) (int64, error) {
	query := "SELECT expires_At FROM otps WHERE id=$1 AND otp=$2"
	if err := r.db.QueryRowContext(ctx, query, otp.Uid, otp.OTP).Scan(&otp.Expires_at); err != nil {
		return 0, err
	}
	return otp.Expires_at, nil
}

func (r *repository) UpdateTokens(ctx context.Context, token string, refresh_token string, updated_at time.Time, id int64) (int64, error) {
	var mobile int64
	query := "UPDATE users SET token=$1, refresh_token=$2, updated_at=$3 WHERE id=$4 RETURNING mobile"

	if err := r.db.QueryRowContext(ctx, query, token, refresh_token, updated_at, id).Scan(&mobile); err != nil {
		return 0, errors.New("can't update the tokens")
	}

	return mobile, nil
}
