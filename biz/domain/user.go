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
	ID         int64
	Username   string
	Email      string
	Followers  int32
	Followings int32
	Friends    int32
	Posts      int32
	Signup     uint64
}

type UserEntry struct {
	ID       int64
	Username string
	IsFollow bool
}
type StatusInfo struct {
	ID         int64
	UserID     int64
	Username   string
	Message    string
	Posted     uint64
	IsLiked    bool
	IsFollowed bool
	PutUrls    []string
	GetUrls    []string
}

type MessageInfo struct {
	ID         int64
	CreatedAt  uint64
	Content    string
	SenderID   int64
	SenderName string
	ChatID     int64
}

type ChatEntry struct {
	ID        int64
	UnseenMsg int32
}
type ChatMessageInfo struct {
	Info []MessageInfo
}
type User interface {
	Reset(ctx context.Context)
	Register(ctx context.Context, req UserRegisterReq) error
	Login(ctx context.Context, username string, password string) (*UserLoginResp, error)
	UserInfo(ctx context.Context, userID int64) (*UserInfoResp, error)
	PostStatus(ctx context.Context, userID int64, message string, numUrls int32) (*StatusInfo, error)
	DeleteStatus(ctx context.Context, userID int64, postID int64) error
	GetTimeline(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*StatusInfo, error)
	GetProfile(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*StatusInfo, error)
	FollowAction(ctx context.Context, userID int64, otherID int64) error
	UnFollowAction(ctx context.Context, userID int64, otherID int64) error
	SearchUser(ctx context.Context, userID int64, expr string) ([]*UserEntry, error)
	GetFollowings(ctx context.Context, userID int64) ([]*UserEntry, error)
	GetFollowers(ctx context.Context, userID int64) ([]*UserEntry, error)
	GetFriends(ctx context.Context, userID int64) ([]*UserEntry, error)
	CreateChat(ctx context.Context, ownerID int64, membersID []int64) (int64, error)
	PostMessage(ctx context.Context, userID int64, chatID int64, message string) (*MessageInfo, error)
	GetPendingMessage(ctx context.Context, userID int64, chatID int64) (*ChatMessageInfo, error)
	LeaveChat(ctx context.Context, userID, chatID int64) error
	GetChatList(ctx context.Context, userID int64) ([]*ChatEntry, error)
	ToggleLikeStatus(ctx context.Context, userID int64, postID int64, action bool) error
	GetHot(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*StatusInfo, error)
}
