package data

import (
	"context"
	"fmt"
	"log"

	"go-social-network/configs"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	redis "github.com/redis/go-redis/v9"
)

var data *Data

func initData() {
	var err error
	data, err = NewData(configs.Data())
	if err != nil {
		hlog.Fatal(err)
	}
}

// Default Get a default database and cache instance
func Default() *Data {
	return data
}

// Data .
type Data struct {
	Redis *redis.Client
}

// NewData .
func NewData(config configs.Config) (data *Data, err error) {
	data = new(Data)
	rdb := newRedisDB(config)
	data.Redis = rdb
	return data, nil
}

func newRedisDB(config configs.Config) (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Connect to redis client failed, err: %v\n", err)
	}
	return rdb
}
