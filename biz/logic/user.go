package logic

import (
	"context"
	"errors"
	"fmt"
	"go-social-network/biz/domain"
	"go-social-network/data"
	"go-social-network/pkg/encrypt"
	"strconv"
	"time"
)

type User struct {
	Data *data.Data
}

func NewUser(data *data.Data) *User {
	return &User{Data: data}
}

func (u *User) Register(ctx context.Context, req domain.UserRegisterReq) error {
	if u.Data.Redis.HGet(ctx, "users:", req.Username).Val() != "" {
		return errors.New("Username already exists")
	}

	id := u.Data.Redis.Incr(ctx, "user:id:").Val()
	password, _ := encrypt.BcryptEncrypt(req.Password)
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.HSet(ctx, "users:", req.Username, id)
	pipeline.HMSet(ctx, fmt.Sprintf("user:%s", strconv.Itoa(int(id))), "username", req.Username,
		"id", id, "followers", 0, "following", 0, "posts", 0, "password", password,
		"signup", time.Now().UnixNano())
	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline err in CreateUser: %s", err)
	}
	return nil
}
