package contact

import (
	"context"

	"github.com/mananKoyawala/whatsapp-clone/internal/user"
)

type Contact struct {
	ID  int64 `json:"id"`
	Uid int64 `json:"uid"`
	Cid int64 `json:"cid"`
}

func NewContact(id, uid, cid int64) *Contact {
	return &Contact{
		ID:  id,
		Uid: uid,
		Cid: cid,
	}
}

type CreateContactReq struct {
	Uid     int64 `json:"uid"`
	CMobile int64 `json:"c_mobile"`
}

type Repository interface {
	AddContact(ctx context.Context, uid, cid int64) (*Contact, error)
	ContactAlreadyExist(ctx context.Context, uid, cid int64) bool
	GetContacts(ctx context.Context, id int64) (*[]user.UserContactsRes, error)
}

type Service interface {
	AddContact(ctx context.Context, req *CreateContactReq) (*Contact, error)
	GetContacts(ctx context.Context, id int64) (*[]user.UserContactsRes, error)
}
