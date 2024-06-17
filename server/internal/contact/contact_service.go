package contact

import (
	"errors"
	"time"

	"github.com/mananKoyawala/whatsapp-clone/internal/user"
	"golang.org/x/net/context"
)

type service struct {
	Repository
	userRespo user.Repository
	timeout   time.Duration
}

func NewContactServ(r Repository, ur user.Repository) Service {
	return &service{
		Repository: r,
		userRespo:  ur,
		timeout:    time.Duration(100) * time.Second,
	}
}

func (s *service) AddContact(ctx context.Context, req *CreateContactReq) (*Contact, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// check if contact exist in the system or not
	cid, err := s.userRespo.GetUserByMobileInt64(ctx, req.CMobile)
	if err != nil {
		return nil, errors.New("contact doesn't exist in the system")
	}

	// check uid and cid already exists
	isExist := s.Repository.ContactAlreadyExist(ctx, req.Uid, cid)
	if isExist {
		return nil, errors.New("contact already added")
	}

	res, err := s.Repository.AddContact(ctx, req.Uid, cid)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetContacts(ctx context.Context, id int64) (*[]user.UserContactsRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.GetContacts(ctx, id)
}
