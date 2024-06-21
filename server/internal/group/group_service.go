package group

import (
	"context"
	"errors"
	"log"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewGroupService(r Repository) Service {
	return &service{Repository: r, timeout: time.Duration(100) * time.Second}
}

func (s *service) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.CreateGroup(ctx, group)
}

func (s *service) AddMemberToGroup(ctx context.Context, req *AddMemberReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	//  check group exits (by id)
	group, err := s.Repository.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return errors.New("group with id doesn't exist")
	}

	// also check that admin is is same as sender id
	if group.AdminID != req.AdminId {
		return errors.New("only admin can add new members")
	}

	// only 20 people per group
	members, err := s.Repository.GetMemberByGroupID(ctx, group.ID)
	if err != nil {
		return err
	}
	if len(members) >= 20 {
		return errors.New("only 20 people can add in the group")
	}

	return s.Repository.AddMemberToGroup(ctx, req.GroupID, req.Members)
}

func (s *service) GetAllGroupByUserID(ctx context.Context, userId int64) (*[]Group, error) {
	var groups []Group
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	groupIds, err := s.Repository.GetAllGroupByUserID(ctx, userId)
	if err != nil {
		return nil, nil
	}

	for _, id := range groupIds {
		group, err := s.Repository.GetGroupByID(ctx, id)
		if err != nil {
			log.Printf("error while getting group info %d", id)
		}
		groups = append(groups, *group)
	}

	return &groups, nil
}

func (s *service) RemoveMemberFromGroup(ctx context.Context, groupId, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.RemoveMemberFromGroup(ctx, groupId, userId)
}

func (s *service) GetGroupDetailsByID(ctx context.Context, groupID int64) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.GetGroupByID(ctx, groupID)
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
		return nil, err
	}

	return group, nil
}

func (s *service) DeleteGroupByID(ctx context.Context, groupID int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repository.DeleteGroupByID(ctx, groupID)
}
