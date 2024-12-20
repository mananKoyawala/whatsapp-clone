package msg

import (
	"context"
	"time"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type Message struct {
	ID             int64     `json:"id"`
	SenderID       int64     `json:"sender_id"`
	ReceiverID     int64     `json:"receiver_id"`
	GroupID        int64     `json:"group_id"`
	IsGroupMessage bool      `json:"is_group_msg"`
	Members        []int64   `json:"members"`
	MessageType    string    `json:"message_type"`
	MessageText    string    `json:"message_text"`
	MediaUrl       string    `json:"media_url"`
	IsRead         bool      `json:"is_read"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
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

type GetAllGroupMessageReq struct {
	GroupID  int64  `json:"group_id"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type MessageReq struct {
	ID         int64 `json:"id"`
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
}

type MessageGroupReq struct {
	ID      int64 `json:"id"`
	GroupID int64 `json:"group_id"`
}

func NewMessage(msg *Message) *Message {
	currentTime, _ := helper.GetTime()
	return &Message{
		SenderID:       msg.SenderID,
		ReceiverID:     msg.ReceiverID,
		GroupID:        msg.GroupID,
		IsGroupMessage: msg.IsGroupMessage,
		Members:        msg.Members,
		MessageType:    msg.MessageType,
		MessageText:    msg.MessageText,
		MediaUrl:       msg.MediaUrl,
		IsRead:         false,
		Created_at:     currentTime,
		Updated_at:     currentTime,
	}
}

type Repository interface {
	AddMessage(ctx context.Context, msg Message) (*Message, error)
	PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error)
	PullAllGroupMessages(ctx context.Context, req *GetAllGroupMessageReq) (*[]Message, error)
	UpdateIsReadMessage(ctx context.Context, req *MessageReq) error
	DeleteMessage(ctx context.Context, msg *MessageReq) error
	DeleteGroupMessage(ctx context.Context, msg *MessageGroupReq) error
	IsMsgExist(ctx context.Context, msg *MessageReq) error
}

type Service interface {
	AddMessage(ctx context.Context, msg *CreateMesReq) (*CreateMesRes, error)
	PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error)
	PullAllGroupMessages(ctx context.Context, req *GetAllGroupMessageReq) (*[]Message, error)
	UpdateIsReadMessage(ctx context.Context, req *[]MessageReq) error
	DeleteGroupMessage(ctx context.Context, msg *MessageGroupReq) error
	DeleteMessage(ctx context.Context, msg *MessageReq) error
}
