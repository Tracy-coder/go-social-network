package logic

import (
	"context"
	"errors"
	"fmt"
	"go-social-network/biz/common"
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
	if u.Data.Redis.HGet(ctx, common.Username2ID, req.Username).Val() != "" {
		return errors.New("Username already exists")
	}

	id := u.Data.Redis.Incr(ctx, common.UserIDCounter).Val()
	password, _ := encrypt.BcryptEncrypt(req.Password)
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.HSet(ctx, common.Username2ID, req.Username, id)
	pipeline.HMSet(ctx, common.UserInfoHashTable(id), "username", req.Username,
		"id", id, "followers", 0, "following", 0, "posts", 0, "password", password,
		"signup", time.Now().UnixNano())
	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline err in CreateUser: %s", err)
	}
	return nil
}

func (u *User) Login(ctx context.Context, username string, password string) (*domain.UserLoginResp, error) {
	id := u.Data.Redis.HGet(ctx, common.Username2ID, username).Val()
	if id == "" {
		return nil, errors.New("User not exists")
	}

	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	if ok := encrypt.BcryptCheck(password, u.Data.Redis.HGet(ctx, common.UserInfoHashTable(idNum), "password").Val()); !ok {
		err := errors.New("wrong password")
		return nil, err
	}

	return &domain.UserLoginResp{
		UserID:   idNum,
		Username: username,
	}, nil
}
