package group

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type service struct {
	Repository
	timeout time.Duration
	logger  *slog.Logger
}

func NewGroupService(r Repository, logger *slog.Logger) Service {
	return &service{Repository: r, timeout: time.Duration(100) * time.Second, logger: logger}
}

func (s *service) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.Repository.CreateGroup(ctx, group)
	if err != nil {
		s.logger.Error("failed to create group", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("group was created", slog.String("groupid", helper.Int64ToStirng(group.ID)))
	return res, nil
}

func (s *service) AddMemberToGroup(ctx context.Context, req *AddMemberReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	//  check group exits (by id)
	group, err := s.Repository.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		msg := fmt.Sprintf("group doesn't exist with id %d", req.GroupID)
		s.logger.Warn(msg, slog.String("error", err.Error()))
		return errors.New(msg)
	}

	// also check that admin is is same as sender id
	if group.AdminID != req.AdminId {
		s.logger.Error("only admin can add new members", slog.String("reqid", helper.Int64ToStirng(req.AdminId)))
		return errors.New("only admin can add new members")
	}

	// only 20 people per group
	members, err := s.Repository.GetMemberByGroupID(ctx, group.ID)
	if err != nil {
		return err
	}
	if len(members) >= 20 {
		s.logger.Error("only 20 people allowed per group", slog.String("groupid", helper.Int64ToStirng(req.GroupID)))
		return errors.New("only 20 people allowed per group")
	}

	if err := s.Repository.AddMemberToGroup(ctx, req.GroupID, req.Members); err != nil {
		s.logger.Error("failed to add member to group", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("member was added into group", slog.String("groupid", helper.Int64ToStirng(group.ID)))
	return nil
}

func (s *service) GetAllGroupByUserID(ctx context.Context, userId int64) (*[]Group, error) {
	var groups []Group
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	groupIds, err := s.Repository.GetAllGroupByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("failed to get all groups by users id "+helper.Int64ToStirng(userId), slog.String("error", err.Error()))
		return nil, err
	}

	for _, id := range groupIds {
		group, err := s.Repository.GetGroupByID(ctx, id)
		if err != nil {

			msg := fmt.Sprintf("error while getting group info %d", id)
			s.logger.Error(msg, slog.String("groupid", helper.Int64ToStirng(group.ID)))
		}
		groups = append(groups, *group)
	}

	s.logger.Info("got all group by user successfully", slog.String("userid", helper.Int64ToStirng(userId)))
	return &groups, nil
}

func (s *service) RemoveMemberFromGroup(ctx context.Context, groupId, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.RemoveMemberFromGroup(ctx, groupId, userId); err != nil {
		msg := fmt.Sprintf("failed to remove member from groupid %d", groupId)
		s.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("memeber was removed from group successfully", slog.String("groupid", helper.Int64ToStirng(groupId)))
	return nil
}

func (s *service) GetGroupDetailsByID(ctx context.Context, groupID int64) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.Repository.GetGroupByID(ctx, groupID)
	if err != nil {
		msg := fmt.Sprintf("failed to get group by id %d", groupID)
		s.logger.Error(msg, slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("got group details", slog.String("groupid", helper.Int64ToStirng(groupID)))
	return res, nil
}

func (s *service) UpdateGroupDetails(ctx context.Context, req UpdateGroup) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	g := Group{
		ID:      req.ID,
		AdminID: req.AdminID,
		Name:    req.Name,
		About:   req.About,
		Image:   req.Image,
	}

	group, err := s.Repository.UpdateGroupDetails(ctx, g)
	if err != nil {
		msg := fmt.Sprintf("failed to update group details by group id %d", g.ID)
		s.logger.Error(msg, slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("group details was updated", slog.String("groupid", helper.Int64ToStirng(g.ID)))
	return group, nil
}

func (s *service) DeleteGroupByID(ctx context.Context, groupID int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.DeleteGroupByID(ctx, groupID); err != nil {
		msg := fmt.Sprintf("failed to delete group by id %d", groupID)
		s.logger.Error(msg, slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("group was deleted", slog.String("groupid", helper.Int64ToStirng(groupID)))
	return nil
}
