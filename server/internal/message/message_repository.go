package msg

import (
	"context"
	"database/sql"
	"log"
)

type repository struct {
	db *sql.DB
}

func NewMsgReposritory(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AddMessage(ctx context.Context, msg Message) (*Message, error) {
	var msgID int64

	query := "INSERT INTO messages (sender_id,receiver_id,message_type,message_text,media_url,is_read,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"

	if err := r.db.QueryRowContext(ctx, query, msg.SenderID, msg.ReceiverID, msg.MessageType, msg.MessageText, msg.MediaUrl, msg.IsRead, msg.Created_at, msg.Updated_at).Scan(&msgID); err != nil {
		return nil, err
	}

	msg.ID = msgID

	return &msg, nil
}

func (r *repository) PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error) {
	log.Println(req.FromDate, req.ToDate)
	query := `
	SELECT * 
	FROM messages 
	WHERE (sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1) 
	AND created_at BETWEEN $3 AND $4
	ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, query, req.SenderID, req.ReceiverID, req.FromDate, req.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.MessageType, &msg.MessageText, &msg.MediaUrl, &msg.Created_at, &msg.Updated_at, &msg.IsRead); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &messages, nil
}
