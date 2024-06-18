package group

import (
	"context"
	"database/sql"
	"errors"
	"log"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type repository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) Repositroy {
	return &repository{db: db}
}

func (r *repository) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	query1 := `
	INSERT INTO groups (admin_id,name,about,image,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
	`

	// creating group
	if err := r.db.QueryRowContext(ctx, query1, group.AdminID, group.Name, group.About, group.Image, group.Created_at, group.Updated_at).Scan(&group.ID); err != nil {
		return nil, err
	}

	// adding initial members
	query2 := `
	INSERT INTO group_members (group_id,member_id) VALUES ($1,$2) RETURNING id
	`
	// adding the members into the group
	for member := range group.Members {
		if err := r.db.QueryRowContext(ctx, query2, group.ID, member); err != nil {
			log.Printf("error occured while adding member id %d", member)
			continue
		}
	}

	return group, nil
}

func (r *repository) AddMemberToGroup(ctx context.Context, groupid int64, members []int64) error {

	current_time, _ := helper.GetTime()

	query := `
	INSERT INTO group_members (g_id,u_id,created_at,updated_at) VALUES ($1,$2,$3,$4)
	`

	// adding the members into the group
	log.Println(groupid)
	for _, member := range members {

		// checking user already added or not
		if ok := r.CheckUserAlreadyInTheGroup(ctx, groupid, member); !ok {

			if err := r.db.QueryRowContext(ctx, query, groupid, member, current_time, current_time); err.Err() != nil {

				log.Printf("error occured while adding member id %d", member)
				return errors.New("error occured while adding member")
			}
		}

	}

	return nil
}

// Get details by group id
func (r *repository) GetGroupByID(ctx context.Context, groupId int64) (*Group, error) {
	query := `
	SELECT * FROM groups WHERE id=$1
	`

	var group Group

	err := r.db.QueryRowContext(ctx, query, groupId).Scan(&group.ID, &group.AdminID, &group.Name, &group.About, &group.Image, &group.Created_at, &group.Updated_at)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// Get group members by group id
func (r *repository) GetMemberByGroupID(ctx context.Context, groupId int64) ([]int64, error) {
	type ID struct {
		Id int64 `json:"u_id"`
	}

	log.Println(">>>", groupId)
	query := `
	SELECT u_id FROM group_members WHERE g_id=$1
	`

	var members []int64

	row, err := r.db.QueryContext(ctx, query, groupId)
	if err != nil {
		return members, err
	}

	for row.Next() {
		var id ID
		if err := row.Scan(&id.Id); err != nil {
			log.Printf("error occurs while scanning id %d", id)
			continue
		}
		members = append(members, id.Id)
	}
	log.Println(members)
	return members, nil
}

func (r *repository) CheckUserAlreadyInTheGroup(ctx context.Context, GroupId, UserId int64) bool {
	query := `
	SELECT * FROM group_members WHERE g_id=$1 AND u_id=$2
	`

	if err := r.db.QueryRowContext(ctx, query, GroupId, UserId); err.Err() != nil {
		return false
	}

	return true
}

// Get all the group in which user in it
