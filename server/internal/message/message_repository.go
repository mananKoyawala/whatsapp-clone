package msg

import (
	"context"
	"database/sql"
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

	query := `
	SELECT * 
	FROM messages 
	WHERE (sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1) 
	OR (sender_id != $1 AND sender_id != $2)
	AND created_at BETWEEN $3 AND $4
	ORDER BY created_at`
	// OR (sender_id != $1 AND sender_id != $2) it acts as a filter so it is not going to find the messages like 1-1
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

// here i taked receiver and sender id for more security beause anyone can update with only message id
func (r *repository) UpdateIsReadMessage(ctx context.Context, req *ReadMessage) error {

	query := `
	UPDATE messages 
	SET is_read=true 
	WHERE id=$1 AND sender_id=$2 AND receiver_id=$3`

	if err := r.db.QueryRowContext(ctx, query, req.ID, req.SenderID, req.ReceiverID); err.Err() != nil {
		return err.Err()
	}

	return nil
}
