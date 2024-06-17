package contact

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddContact(ctx context.Context, uid, cid int64) (*Contact, error) {

	var c Contact

	query := `
	INSERT INTO contacts (uid,cid) VALUES ($1,$2) RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, uid, cid).Scan(&c.ID); err != nil {
		return nil, err
	}

	c.Uid = uid
	c.Cid = cid
	return &c, nil
}

func (r *repository) ContactAlreadyExist(ctx context.Context, uid, cid int64) bool {

	var id int64
	query := `
	SELECT id FROM contacts WHERE uid=$1 AND cid=$2
	`

	r.db.QueryRowContext(ctx, query, uid, cid).Scan(&id)
	return id > 0
}
