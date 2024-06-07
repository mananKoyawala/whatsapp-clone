package msg

import (
	"context"
	"time"
)

type Message struct {
	ID          int64     `json:"id"`
	SenderID    int64     `json:"sender_id"`
	ReceiverID  int64     `json:"receiver_id"`
	MessageType string    `json:"message_type"`
	MessageText string    `json:"message_text"`
	MediaUrl    string    `json:"media_url"`
	IsRead      bool      `json:"is_read"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type CreateMesReq struct {
	SenderID    int64  `json:"sender_id"`
	ReceiverID  int64  `json:"receiver_id"`
	MessageType string `json:"message_type"`
	MessageText string `json:"message_text"`
	MediaUrl    string `json:"media_url"`
}

type CreateMesRes struct {
	ID          int64     `json:"id"`
	SenderID    int64     `json:"sender_id"`
	ReceiverID  int64     `json:"receiver_id"`
	MessageType string    `json:"message_type"`
	MessageText string    `json:"message_text"`
	MediaUrl    string    `json:"media_url"`
	IsRead      bool      `json:"is_read"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type GetAllMessageReq struct {
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
}

type Repository interface {
	AddMessage(ctx context.Context, msg Message) (*Message, error)
	PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error)
}

type Service interface {
	AddMessage(ctx context.Context, msg *CreateMesReq) (*CreateMesRes, error)
	PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error)
}
