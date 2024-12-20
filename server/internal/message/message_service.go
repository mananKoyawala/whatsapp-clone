package msg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log/slog"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/group"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
)

type service struct {
	Repository
	userRepo  user.Repository
	groupRepo group.Repository
	timeout   time.Duration
	logger    *slog.Logger
}

func NewMsgService(r Repository, userRepo user.Repository, groupRepo group.Repository, logger *slog.Logger) Service {
	return &service{Repository: r, userRepo: userRepo, groupRepo: groupRepo, timeout: time.Duration(100) * time.Second, logger: logger}
}

func (s *service) AddMessage(ctx context.Context, msg *CreateMesReq) (*CreateMesRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if sender and receivier's are exits or not
	_, err := s.userRepo.GetUserById(ctx, msg.SenderID)
	if err != nil {
		s.logger.Warn("sender doesn't exist", slog.String("userid", helper.Int64ToStirng(msg.SenderID)))
		return nil, errors.New("sender does not exist")
	}

	_, err = s.userRepo.GetUserById(ctx, msg.ReceiverID)
	if err != nil {
		s.logger.Warn("receiver doesn't exist", slog.String("userid", helper.Int64ToStirng(msg.ReceiverID)))
		return nil, errors.New("receiver does not exist")
	}

	current_time, _ := helper.GetTime()

	newMsg := &Message{
		SenderID:    msg.SenderID,
		ReceiverID:  msg.ReceiverID,
		MessageType: msg.MessageType,
		MessageText: msg.MessageText,
		MediaUrl:    msg.MediaUrl,
		IsRead:      false,
		Created_at:  current_time,
		Updated_at:  current_time,
	}

	r, err := s.Repository.AddMessage(ctx, *newMsg)
	if err != nil {
		s.logger.Error("failed to add message", slog.String("error", err.Error()))
		return nil, err
	}

	res := &CreateMesRes{
		ID:          r.ID,
		SenderID:    r.SenderID,
		ReceiverID:  r.ReceiverID,
		MessageType: r.MessageType,
		MessageText: r.MessageText,
		MediaUrl:    r.MediaUrl,
		IsRead:      r.IsRead,
		Created_at:  r.Created_at,
		Updated_at:  r.Updated_at,
	}

	s.logger.Info("message added successfully", slog.String("messageid", helper.Int64ToStirng(newMsg.ID)))
	return res, nil
}

func (s *service) PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if sender and receivier's are exits or not
	_, err := s.userRepo.GetUserById(ctx, req.SenderID)
	if err != nil {
		s.logger.Warn("sender doesn't exist", slog.String("userid", helper.Int64ToStirng(req.SenderID)))
		return nil, errors.New("sender does not exist")
	}

	_, err = s.userRepo.GetUserById(ctx, req.ReceiverID)
	if err != nil {
		s.logger.Warn("receiver doesn't exist", slog.String("userid", helper.Int64ToStirng(req.ReceiverID)))
		return nil, errors.New("receiver does not exist")
	}

	if req.FromDate == "" {
		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		req.FromDate = firstOfMonth.Format("01-02-2006 15:04:05")
	}

	if req.ToDate == "" {
		req.ToDate = time.Now().Format("01-02-2006 15:04:05")
	}

	res, err := s.Repository.PullAllMessages(ctx, req)
	if err != nil {
		s.logger.Error("failed to pull all messages", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("pulled all messages successfully.")
	return res, nil
}

func (s *service) UpdateIsReadMessage(ctx context.Context, req *[]MessageReq) error {

	for _, msg := range *req {
		if err := s.Repository.UpdateIsReadMessage(ctx, &msg); err != nil {
			s.logger.Error("failed to update read message", slog.String("error", err.Error()))
			return err
		}
	}

	s.logger.Info("read messages updated successfully.")
	return nil
}

func (s *service) DeleteMessage(ctx context.Context, msg *MessageReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.IsMsgExist(ctx, msg); err != nil {
		msg := fmt.Sprintf("message with %d id doesn't exist", msg.ID)
		s.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	if err := s.Repository.DeleteMessage(ctx, msg); err != nil {
		s.logger.Error("failed to delete message", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("message deleted successfully", slog.String("messageid", helper.Int64ToStirng(msg.ID)))
	return nil
}

func (s *service) PullAllGroupMessages(ctx context.Context, req *GetAllGroupMessageReq) (*[]Message, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check gorup exists or not
	_, err := s.groupRepo.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		s.logger.Warn("group doesn't exist", slog.String("groupid", helper.Int64ToStirng(req.GroupID)))
		return nil, errors.New("group does not exist")
	}

	if req.FromDate == "" {
		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		req.FromDate = firstOfMonth.Format("01-02-2006 15:04:05")
	}

	if req.ToDate == "" {
		req.ToDate = time.Now().Format("01-02-2006 15:04:05")
	}

	res, err := s.Repository.PullAllGroupMessages(ctx, req)
	if err != nil {
		s.logger.Error("failed to pull all group messages", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("all group messages pulled successfully.", slog.String("groupid", helper.Int64ToStirng(req.GroupID)))
	return res, nil
}

func (s *service) DeleteGroupMessage(ctx context.Context, msg *MessageGroupReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.groupRepo.GetGroupByID(ctx, msg.GroupID)
	if err != nil {
		s.logger.Warn("group doesn't exist", slog.String("groupid", helper.Int64ToStirng(msg.GroupID)))
		return errors.New("group doesn't exist")
	}

	if err := s.Repository.DeleteGroupMessage(ctx, msg); err != nil {
		s.logger.Error("failed to delete group message", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("group message delete successfully.", slog.String("groupid", helper.Int64ToStirng(msg.GroupID)))
	return nil
}
