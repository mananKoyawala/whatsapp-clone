package contact

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mananKoyawala/whatsapp-clone/internal/user"
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

func (r *repository) GetContacts(ctx context.Context, id int64) (*[]user.UserContactsRes, error) {
	query := `
	SELECT u.id,u.name,u.mobile,u.about,u.image,u.last_seen,u.is_online
	FROM users u
	INNER JOIN contacts c ON u.id = c.cid
	WHERE c.uid=$1;
	`

	row, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var users []user.UserContactsRes
	for row.Next() {
		var user user.UserContactsRes
		if err := row.Scan(&user.ID, &user.Name, &user.Mobile, &user.About, &user.Image, &user.Last_Seen, &user.Is_Online); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		users = append(users, user)
	}

	return &users, nil
}
