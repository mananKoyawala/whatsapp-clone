package msg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewMsgReposritory(db *sql.DB, logger *slog.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) AddMessage(ctx context.Context, msg Message) (*Message, error) {
	var msgID int64

	query := "INSERT INTO messages (sender_id,receiver_id,group_id,is_group_msg,message_type,message_text,media_url,is_read,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id"

	if err := r.db.QueryRowContext(ctx, query, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.IsGroupMessage, msg.MessageType, msg.MessageText, msg.MediaUrl, msg.IsRead, msg.Created_at, msg.Updated_at).Scan(&msgID); err != nil {
		r.logger.Error("failed to add message", slog.String("error", err.Error()))
		return nil, err
	}

	msg.ID = msgID
	r.logger.Info("message added successfully.", slog.String("messageid", helper.Int64ToStirng(msg.ID)))
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
		r.logger.Error("failed to pull all messages", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.MessageType, &msg.MessageText, &msg.MediaUrl, &msg.Created_at, &msg.Updated_at, &msg.IsRead, &msg.GroupID, &msg.IsGroupMessage); err != nil {
			r.logger.Error("failed to scan message", slog.String("error", err.Error()))
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("failed to scan row", slog.String("error", err.Error()))
		return nil, err
	}

	r.logger.Debug("all messages are pulled successfully.")
	return &messages, nil
}

// here i taked receiver and sender id for more security beause anyone can update with only message id
func (r *repository) UpdateIsReadMessage(ctx context.Context, req *MessageReq) error {

	query := `
	UPDATE messages 
	SET is_read=true 
	WHERE id=$1 AND sender_id=$2 AND receiver_id=$3`

	if err := r.db.QueryRowContext(ctx, query, req.ID, req.SenderID, req.ReceiverID); err.Err() != nil {
		r.logger.Error("failed to update is_read message", slog.String("error", err.Err().Error()))
		return err.Err()
	}

	r.logger.Info("read message updated successfuly.")
	return nil
}

func (r *repository) DeleteMessage(ctx context.Context, msg *MessageReq) error {
	query := `
	DELETE FROM messages 
	WHERE id=$1 AND sender_id=$2 AND receiver_id=$3
	`

	if err := r.db.QueryRowContext(ctx, query, msg.ID, msg.SenderID, msg.ReceiverID); err.Err() != nil {
		r.logger.Error("failed to delete message", slog.String("error", err.Err().Error()))
		return err.Err()
	}

	r.logger.Info("message deleted successfully.")
	return nil
}

func (r *repository) IsMsgExist(ctx context.Context, msg *MessageReq) error {
	query := "SELECT * FROM messages WHERE id=$1 AND sender_id=$2 AND receiver_id=$3"

	row, err := r.db.QueryContext(ctx, query, msg.ID, msg.SenderID, msg.ReceiverID)
	if err != nil {
		msg := fmt.Sprintf("message doesn't exist with %d id", msg.ID)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}
	defer row.Close()

	hasData := row.Next()
	if !hasData {
		r.logger.Error("message doesn't exists")
		return errors.New("message doesn't exists")
	}

	r.logger.Info("message does exist")
	return nil
}

func (r *repository) PullAllGroupMessages(ctx context.Context, req *GetAllGroupMessageReq) (*[]Message, error) {

	query := `
	SELECT * 
	FROM messages 
	WHERE group_id=$1 
	AND created_at BETWEEN $2 AND $3
	ORDER BY created_at`
	// OR (sender_id != $1 AND sender_id != $2) it acts as a filter so it is not going to find the messages like 1-1
	rows, err := r.db.QueryContext(ctx, query, req.GroupID, req.FromDate, req.ToDate)
	if err != nil {
		r.logger.Error("failed to pull all group messages", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.MessageType, &msg.MessageText, &msg.MediaUrl, &msg.Created_at, &msg.Updated_at, &msg.IsRead, &msg.GroupID, &msg.IsGroupMessage); err != nil {
			r.logger.Error("failed to scan row message", slog.String("error", err.Error()))
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("failed while scanning row", slog.String("error", err.Error()))
		return nil, err
	}

	r.logger.Debug("all group messages are pulled successfully.")
	return &messages, nil
}

func (r *repository) DeleteGroupMessage(ctx context.Context, msg *MessageGroupReq) error {
	query := `
	DELETE FROM messages 
	WHERE id=$1 AND group_id=$2
	`

	if err := r.db.QueryRowContext(ctx, query, msg.ID, msg.GroupID); err.Err() != nil {
		r.logger.Error("failed to delete group", slog.String("error", err.Err().Error()))
		return err.Err()
	}

	r.logger.Info("message deleted successfully.")
	return nil
}
