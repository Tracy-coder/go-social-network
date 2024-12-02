package domain

import "context"

type UserRegisterReq struct {
	Username string
	Password string
	Email    string
}

type UserLoginResp struct {
	UserID   int64
	Username string
}
type User interface {
	Register(ctx context.Context, req UserRegisterReq) error
	Login(ctx context.Context, username string, password string) (*UserLoginResp, error)
}
