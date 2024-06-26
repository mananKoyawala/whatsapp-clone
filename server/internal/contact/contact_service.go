package contact

import (
	"errors"
	"log/slog"
	"time"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"golang.org/x/net/context"
)

type service struct {
	Repository
	userRespo user.Repository
	timeout   time.Duration
	logger    *slog.Logger
}

func NewContactServ(r Repository, ur user.Repository, logger *slog.Logger) Service {
	return &service{
		Repository: r,
		userRespo:  ur,
		timeout:    time.Duration(100) * time.Second,
		logger:     logger,
	}
}

func (s *service) AddContact(ctx context.Context, req *CreateContactReq) (*Contact, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if contact exist in the system or not
	cid, err := s.userRespo.GetUserByMobileInt64(ctx, req.CMobile)
	if err != nil {
		s.logger.Error("contact doesn't exist in the system", slog.String("error", err.Error()))
		return nil, errors.New("contact doesn't exist in the system")
	}

	// check uid and cid already exists
	isExist := s.Repository.ContactAlreadyExist(ctx, req.Uid, cid)
	if isExist {
		s.logger.Error("contact already exist", slog.String("userid", helper.Int64ToStirng(req.Uid)))
		return nil, errors.New("contact already exist")
	}

	res, err := s.Repository.AddContact(ctx, req.Uid, cid)
	if err != nil {
		s.logger.Error("failed to add contact", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("contact added", slog.String("contactid", helper.Int64ToStirng(res.ID)))
	return res, nil
}

func (s *service) GetContacts(ctx context.Context, id int64) (*[]user.UserContactsRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	contacts, err := s.Repository.GetContacts(ctx, id)
	if err != nil {
		s.logger.Error("failed to get contacts", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("contacts were got", slog.String("userid", helper.Int64ToStirng(id)))
	return contacts, nil
}
