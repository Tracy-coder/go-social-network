package domain

import (
	"context"
)

type UserRegisterReq struct {
	Username string
	Password string
	Email    string
}

type UserLoginResp struct {
	UserID   int64
	Username string
}

type UserInfoResp struct {
	ID        int64
	Username  string
	Email     string
	Followers int32
	Following int32
	Posts     int32
	Signup    uint64
}

type StatusInfo struct {
	ID       int64
	UserID   int64
	Username string
	Message  string
	Posted   uint64
}

type User interface {
	Reset(ctx context.Context)
	Register(ctx context.Context, req UserRegisterReq) error
	Login(ctx context.Context, username string, password string) (*UserLoginResp, error)
	UserInfo(ctx context.Context, userID int64) (*UserInfoResp, error)
	PostStatus(ctx context.Context, userID int64, message string) (*StatusInfo, error)
	DeleteStatus(ctx context.Context, userID int64, postID int64) error
	GetTimeline(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*StatusInfo, error)
	FollowAction(ctx context.Context, userID int64, otherID int64) error
	UnFollowAction(ctx context.Context, userID int64, otherID int64) error
}
