package data

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"go-social-network/configs"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	redis "github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
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

func (d *Data) AcquireLockWithTimeout(ctx context.Context, lockname string, acquireTimeout, lockTimeout float64) string {
	identifier := uuid.NewV4().String()
	lockname = "lock:" + lockname
	finalLockTimeout := math.Ceil(lockTimeout)

	end := time.Now().UnixNano() + int64(acquireTimeout*1e9)
	for time.Now().UnixNano() < end {
		if d.Redis.SetNX(ctx, lockname, identifier, 0).Val() {
			d.Redis.Expire(ctx, lockname, time.Duration(finalLockTimeout)*time.Second)
			return identifier
		} else if d.Redis.TTL(ctx, lockname).Val() < 0 {
			d.Redis.Expire(ctx, lockname, time.Duration(finalLockTimeout)*time.Second)
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ""
}

func (d *Data) ReleaseLock(ctx context.Context, lockname, identifier string) bool {
	lockname = "lock:" + lockname
	lostLock := false
	for {
		err := d.Redis.Watch(ctx, func(tx *redis.Tx) error {
			if tx.Get(ctx, lockname).Val() == identifier {
				pipe := tx.TxPipeline()
				pipe.Del(ctx, lockname)
				_, err := pipe.Exec(ctx)
				return err
			}
			// lock was grabbed by others
			lostLock = true
			return nil
		}, lockname)

		if err != nil {
			log.Println("watch failed in ReleaseLock, err is: ", err)
			continue
		}

		if lostLock {
			return true
		}
	}
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
