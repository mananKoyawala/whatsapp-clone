package group

import (
	"context"
	"time"

	helper "github.com/mananKoyawala/whatsapp-clone/helpers"
)

type Group struct {
	ID         int64     `json:"id"`
	AdminID    int64     `json:"admin_id"`
	Name       string    `json:"name"`
	About      string    `json:"about"`
	Image      string    `json:"image"`
	Members    []int64   `json:"members"`
	Created_at time.Time `json:"create_at"`
	Updated_at time.Time `json:"updated_at"`
}

func NewGroup(req CreateGroupReq) *Group {
	currentTime, _ := helper.GetTime()
	return &Group{
		AdminID:    req.AdminID,
		Name:       req.Name,
		About:      req.About,
		Image:      req.Image,
		Members:    []int64{req.AdminID},
		Created_at: currentTime,
		Updated_at: currentTime,
	}
}

type CreateGroupReq struct {
	AdminID int64   `json:"admin_id"`
	Name    string  `json:"name"`
	About   string  `json:"about"`
	Image   string  `json:"image"`
	Members []int64 `json:"members"`
}

type AddMemberReq struct {
	AdminId int64   `json:"admin_id"`
	GroupID int64   `json:"group_id"`
	Members []int64 `json:"members"`
}

type Repositroy interface {
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	AddMemberToGroup(ctx context.Context, groupid int64, members []int64) error
	GetGroupByID(ctx context.Context, groupId int64) (*Group, error)
	GetMemberByGroupID(ctx context.Context, groupId int64) ([]int64, error)
	CheckUserAlreadyInTheGroup(ctx context.Context, GroupId, UserId int64) bool
}

type Service interface {
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	AddMemberToGroup(ctx context.Context, req *AddMemberReq) error
}
