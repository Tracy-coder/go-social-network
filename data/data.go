package data

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"go-social-network/biz/common"
	"go-social-network/configs"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	redis "github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
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
	cron  *cron.Cron
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

func (d *Data) FreshStatus() {
	users := d.Redis.HGetAll(context.Background(), common.Username2ID).Val()
	for _, v := range users {
		cnt := d.Redis.ZCard(context.Background(), common.HomeTimelineZSet(v)).Val()
		userID, _ := strconv.ParseInt(v, 10, 64)
		if cnt > common.HomeTimelineSize {
			d.Redis.ZRemRangeByRank(context.Background(), common.HomeTimelineZSet(userID), 0, -common.HomeTimelineSize-1)
		} else if cnt < common.HomeTimelineSize {
			d.refillTimeline(context.Background(), userID)
		}
	}
	cnt := d.Redis.ZCard(context.Background(), common.HotStatusZSet).Val()
	if cnt > common.HotStatusZSetSize {
		d.Redis.ZRemRangeByRank(context.Background(), common.HotStatusZSet, 0, -common.HotStatusZSetSize-1)
	}
}

func (d *Data) refillTimeline(ctx context.Context, userID int64) error {
	followings := d.Redis.ZRange(ctx, common.FollowingZSet(userID), 0, -1).Val()
	pipeline := d.Redis.TxPipeline()
	posts := make([]string, 0)

	for _, follower := range followings {
		pipeline.ZRange(ctx, common.UserProfileZSet(follower), 0, common.HomeTimelineSize)
	}
	res, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline error in refill timeline:%s %s", err, res)
	}

	for _, cmd := range res {
		tmp := cmd.(*redis.StringSliceCmd).Val()
		posts = append(posts, tmp...)
	}
	fmt.Println(posts)
	pipeline = d.Redis.TxPipeline()
	for _, id := range posts {
		pipeline.HGet(ctx, common.StatusInfoHashTable(id), "posted")
	}
	res, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline error in refill timeline:%s %s", err, res)
	}
	for i, cmd := range res {
		postedString := cmd.(*redis.StringCmd).Val()
		posted, _ := strconv.ParseUint(postedString, 10, 64)
		pipeline.ZAdd(ctx, common.HomeTimelineZSet(userID), redis.Z{Score: float64(posted), Member: posts[i]})
	}
	pipeline.ZRemRangeByRank(ctx, common.HomeTimelineZSet(userID), 0, -common.HomeTimelineSize-1)
	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline error in refill timeline:%s", err)
	}
	return nil
}

// NewData .
func NewData(config configs.Config) (data *Data, err error) {
	data = new(Data)
	rdb := newRedisDB(config)
	data.Redis = rdb
	data.cron = cron.New()
	_, err = data.cron.AddFunc(config.CronExpr, data.FreshStatus)
	if err != nil {
		return nil, err
	}
	data.cron.Start()
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
