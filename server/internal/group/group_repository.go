package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewGroupRepository(db *sql.DB, logger *slog.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	query1 := `
	INSERT INTO groups (admin_id,name,about,image,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
	`

	// creating group
	if err := r.db.QueryRowContext(ctx, query1, group.AdminID, group.Name, group.About, group.Image, group.Created_at, group.Updated_at).Scan(&group.ID); err != nil {
		r.logger.Error("failed to create group", slog.String("error", err.Error()))
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
			msg := fmt.Sprintf("error occured while adding member id %d", member)
			r.logger.Error(msg, slog.String("error", err.Err().Error()))
			continue
		}
	}

	r.logger.Info("group was created", slog.String("groupid", helper.Int64ToStirng(group.ID)))
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

				msg := fmt.Sprintf("error occured while adding member id %d", member)
				r.logger.Error(msg, slog.String("error", err.Err().Error()))
				return errors.New("error occured while adding member")
			}
		}

	}

	r.logger.Info("members was added into group")
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
		msg := fmt.Sprintf("failed to get group id %d", groupId)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return nil, err
	}

	r.logger.Debug("got group", slog.String("groupid", helper.Int64ToStirng(groupId)))
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
		msg := fmt.Sprintf("failed to get members from group id %d", groupId)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return members, err
	}

	for row.Next() {
		var id ID
		if err := row.Scan(&id.Id); err != nil {
			msg := fmt.Sprintf("error occurs while scanning id %d", id)
			r.logger.Error(msg, slog.String("error", err.Error()))
			continue
		}
		members = append(members, id.Id)
	}

	r.logger.Debug("got all members by group", slog.String("groupid", helper.Int64ToStirng(groupId)))
	return members, nil
}

func (r *repository) CheckUserAlreadyInTheGroup(ctx context.Context, GroupId, UserId int64) bool {
	var id int64
	query := `
	SELECT id FROM group_members WHERE g_id=$1 AND u_id=$2
	`

	err := r.db.QueryRowContext(ctx, query, GroupId, UserId).Scan(&id)
	if err != nil || id <= 0 {
		r.logger.Error("failed to check user already exist in the group", slog.String("error", err.Error()))
		return true
	}

	r.logger.Info("user already in group", slog.String("groupid", helper.Int64ToStirng(GroupId)))
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
		r.logger.Error("failed to get all group by user", slog.String("userid", helper.Int64ToStirng(userId)))
		return groups, err
	}

	for row.Next() {
		var id ID
		if err := row.Scan(&id.Id); err != nil {
			msg := fmt.Sprintf("error occurs while scanning gid %d", id)
			r.logger.Error(msg, slog.String("error", err.Error()))
			continue
		}
		groups = append(groups, id.Id)
	}

	r.logger.Debug("got all groups by user", slog.String("userid", helper.Int64ToStirng(userId)))
	return groups, nil
}

// remove member from group (slightly same as add to group) only admin can do this
func (r *repository) RemoveMemberFromGroup(ctx context.Context, groupId, userId int64) error {

	query := `
	DELETE FROM group_members WHERE g_id=$1 AND u_id=$2
	`

	if ok := r.CheckUserAlreadyInTheGroup(ctx, groupId, userId); ok {
		r.logger.Error("user doesn't belongs to the group")
		return errors.New("user doesn't belongs to the group")
	}

	_, err := r.db.ExecContext(ctx, query, groupId, userId)
	if err != nil {
		msg := fmt.Sprintf("failed to remove member form groupid %d", groupId)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	r.logger.Info("member was removed from group", slog.String("groupid", helper.Int64ToStirng(groupId)))
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
		r.logger.Error("failed to updated group details", slog.String("error", err.Err().Error()))
		return nil, err.Err()
	}

	r.logger.Info("group deatils updated successfully", slog.String("groupid", helper.Int64ToStirng(group.ID)))
	return &group, nil
}

// delete group then remove all the people also from group_members
func (r *repository) DeleteGroupByID(ctx context.Context, groupID int64) error {

	// begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to initiate transactions", slog.String("error", err.Error()))
		return err
	}

	query := `
	DELETE FROM groups WHERE id=$1
	`

	// check group already exist or not
	_, err = r.GetGroupByID(ctx, groupID)
	if err != nil {
		tx.Rollback()
		msg := fmt.Sprintf("group doesn't exists with id %d", groupID)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return errors.New(msg)
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
		msg := fmt.Sprintf("group doesn't exists with id %d", groupID)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	// delete group
	_, err = r.db.ExecContext(ctx, query, groupID)
	if err != nil {
		tx.Rollback()
		msg := fmt.Sprintf("group doesn't exists with id %d", groupID)
		r.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		r.logger.Error("failed to commit transactions", slog.String("error", err.Error()))
	} else {
		r.logger.Debug("transation commited")
	}

	r.logger.Info("group was deleted")
	return nil
}

// TODO : if admin leaves then make admin as after just added user (make common func for all who want to leave)
