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

func NewGroupRepository(db *sql.DB) Repository {
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
	current_time, _ := helper.GetTime()

	query2 := `
	INSERT INTO group_members (g_id,u_id,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING id
	`
	// adding the members into the group
	for _, member := range group.Members {
		if err := r.db.QueryRowContext(ctx, query2, group.ID, member, current_time, current_time); err.Err() != nil {
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
	for _, member := range members {

		// checking user already added or not
		if ok := r.CheckUserAlreadyInTheGroup(ctx, groupid, member); ok {

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

	return members, nil
}

func (r *repository) CheckUserAlreadyInTheGroup(ctx context.Context, GroupId, UserId int64) bool {
	var id int64
	query := `
	SELECT id FROM group_members WHERE g_id=$1 AND u_id=$2
	`

	err := r.db.QueryRowContext(ctx, query, GroupId, UserId).Scan(&id)
	if err != nil || id <= 0 {
		return true
	}

	return false
}

// Get all the group in which user in it
func (r *repository) GetAllGroupByUserID(ctx context.Context, userId int64) ([]int64, error) {
	type ID struct {
		Id int64 `json:"g_id"`
	}

	query := `
	SELECT g_id FROM group_members WHERE u_id=$1
	`

	var groups []int64

	row, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return groups, err
	}

	for row.Next() {
		var id ID
		if err := row.Scan(&id.Id); err != nil {
			log.Printf("error occurs while scanning gid %d", id)
			continue
		}
		groups = append(groups, id.Id)
	}
	// log.Println(groups)
	return groups, nil
}

// remove member from group (slightly same as add to group) only admin can do this
func (r *repository) RemoveMemberFromGroup(ctx context.Context, groupId, userId int64) error {

	query := `
	DELETE FROM group_members WHERE g_id=$1 AND u_id=$2
	`

	if ok := r.CheckUserAlreadyInTheGroup(ctx, groupId, userId); ok {
		return errors.New("user doesn't belongs to the group")
	}

	_, err := r.db.ExecContext(ctx, query, groupId, userId)
	if err != nil {
		return err
	}

	return nil
}

// update group details
func (r *repository) UpdateGroupDetails(ctx context.Context, group Group) (*Group, error) {

	group.Updated_at, _ = helper.GetTime()
	query := `
	UPDATE groups 
	SET admin_id=$1 , name=$2 , about=$3, image=$4 , updated_at=$5
	WHERE id=$6 
	`

	if err := r.db.QueryRowContext(ctx, query, group.AdminID, group.Name, group.About, group.Image, group.Updated_at, group.ID); err.Err() != nil {
		return nil, err.Err()
	}

	return &group, nil
}

// delete group then remove all the people also from group_members
func (r *repository) DeleteGroupByID(ctx context.Context, groupID int64) error {

	// begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
	DELETE FROM groups WHERE id=$1
	`

	// check group already exist or not
	_, err = r.GetGroupByID(ctx, groupID)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return errors.New("group doesn't exist")
	}

	//  delete first all the members related the group
	queryDeleteAllMembers := `
	DELETE FROM group_members 
	WHERE g_id=$1
	`

	// delete all group members
	_, err = r.db.ExecContext(ctx, queryDeleteAllMembers, groupID)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// delete group
	_, err = r.db.ExecContext(ctx, query, groupID)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("....Transaction committed")
	}

	return nil
}

// TODO : if admin leaves then make admin as after just added user (make common func for all who want to leave)
