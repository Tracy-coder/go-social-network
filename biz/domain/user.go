package domain

import "context"

type UserRegisterReq struct {
	Username string
	Password string
	Email    string
}

type User interface {
	Register(ctx context.Context, req UserRegisterReq) error
}
