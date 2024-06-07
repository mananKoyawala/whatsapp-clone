package msg

import (
	"context"
	"errors"
	"time"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
)

type service struct {
	Repository
	userRepo user.Repository
	timeout  time.Duration
}

func NewMsgService(r Repository, userRepo user.Repository) Service {
	return &service{Repository: r, userRepo: userRepo, timeout: time.Duration(100) * time.Second}
}

func (s *service) AddMessage(ctx context.Context, msg *CreateMesReq) (*CreateMesRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if sender and receivier's are exits or not
	_, err := s.userRepo.GetUserById(ctx, msg.SenderID)
	if err != nil {
		return nil, errors.New("sender does not exist")
	}

	_, err = s.userRepo.GetUserById(ctx, msg.ReceiverID)
	if err != nil {
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

	return res, nil
}

func (s *service) PullAllMessages(ctx context.Context, req *GetAllMessageReq) (*[]Message, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if sender and receivier's are exits or not
	_, err := s.userRepo.GetUserById(ctx, req.SenderID)
	if err != nil {
		return nil, errors.New("sender does not exist")
	}

	_, err = s.userRepo.GetUserById(ctx, req.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver does not exist")
	}

	if req.FromDate == "" {
		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		req.FromDate = firstOfMonth.Format("01-02-2006")
	}

	if req.ToDate == "" {
		req.ToDate = time.Now().Format("01-02-2006")
	}

	res, err := s.Repository.PullAllMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
