package msg

import "time"

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
