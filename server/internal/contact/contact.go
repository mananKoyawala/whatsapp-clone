package contact

import "context"

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
}

type Service interface {
	AddContact(ctx context.Context, req *CreateContactReq) (*Contact, error)
}
