package contact

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
	"github.com/mananKoyawala/whatsapp-clone/internal/user"
)

type repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewContactRepo(db *sql.DB, logger *slog.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) AddContact(ctx context.Context, uid, cid int64) (*Contact, error) {

	var c Contact

	query := `
	INSERT INTO contacts (uid,cid) VALUES ($1,$2) RETURNING id
	`

	if err := r.db.QueryRowContext(ctx, query, uid, cid).Scan(&c.ID); err != nil {
		r.logger.Error("failed to add contact", slog.String("error", err.Error()))
		return nil, err
	}

	// get added ids
	c.Uid = uid
	c.Cid = cid

	r.logger.Info("contact added", slog.String("contactid", helper.Int64ToStirng(c.Cid)))
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
		r.logger.Error("failed to get contacts", slog.String("error", err.Error()))
		return nil, err
	}
	defer row.Close()

	var users []user.UserContactsRes
	for row.Next() {
		var user user.UserContactsRes
		if err := row.Scan(&user.ID, &user.Name, &user.Mobile, &user.About, &user.Image, &user.Last_Seen, &user.Is_Online); err != nil {
			msg := fmt.Sprintln("Error scanning row:", err.Error())
			r.logger.Error(msg)
			continue
		}
		users = append(users, user)
	}

	r.logger.Debug("contacts were got", slog.String("userid", helper.Int64ToStirng(id)))
	return &users, nil
}
