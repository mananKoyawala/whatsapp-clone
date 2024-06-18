package group

import (
	"context"
	"errors"
	"time"
)

type service struct {
	Repositroy
	timeout time.Duration
}

func NewGroupService(r Repositroy) Service {
	return &service{Repositroy: r, timeout: time.Duration(100) * time.Second}
}

func (s *service) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.Repositroy.CreateGroup(ctx, group)
}

func (s *service) AddMemberToGroup(ctx context.Context, req *AddMemberReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	//  check group exits (by id)
	group, err := s.Repositroy.GetGroupByID(ctx, req.GroupID)
	if err != nil {
		return errors.New("group with id doesn't exist")
	}

	// also check that admin is is same as sender id
	if group.AdminID != req.AdminId {
		return errors.New("only admin can add new members")
	}

	// only 20 people per group
	members, err := s.Repositroy.GetMemberByGroupID(ctx, group.ID)
	if err != nil {
		return err
	}
	if len(members) >= 20 {
		return errors.New("only 20 people can add in the group")
	}

	return s.Repositroy.AddMemberToGroup(ctx, req.GroupID, req.Members)
}
