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

	"github.com/redis/go-redis/v9"
)

type User struct {
	Data *data.Data
}

func NewUser(data *data.Data) domain.User {
	return &User{Data: data}
}

func (u *User) Reset(ctx context.Context) {
	u.Data.Redis.FlushDB(ctx)
}

func (u *User) Register(ctx context.Context, req domain.UserRegisterReq) error {
	if u.Data.Redis.HGet(ctx, common.Username2ID, req.Username).Val() != "" {
		return errors.New("Username already exists")
	}
	lock := u.Data.AcquireLockWithTimeout(ctx, "user:"+req.Username, 10, 10)

	if lock == "" {
		return errors.New("concurrency errors")
	}
	defer u.Data.ReleaseLock(ctx, "user:"+req.Username, lock)
	id := u.Data.Redis.Incr(ctx, common.UserIDCounter).Val()
	password, _ := encrypt.BcryptEncrypt(req.Password)
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.HSet(ctx, common.Username2ID, req.Username, id)
	pipeline.HMSet(ctx, common.UserInfoHashTable(id), "username", req.Username,
		"id", id, "email", req.Email, "followers", 0, "following", 0, "posts", 0, "password", password,
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

func (u *User) UserInfo(ctx context.Context, userID int64) (*domain.UserInfoResp, error) {
	info := u.Data.Redis.HGetAll(ctx, common.UserInfoHashTable(userID)).Val()
	followers, _ := strconv.ParseInt(info["followers"], 10, 64)
	following, _ := strconv.ParseInt(info["following"], 10, 64)
	id, _ := strconv.ParseInt(info["id"], 10, 64)
	posts, _ := strconv.ParseInt(info["posts"], 10, 64)
	signup, _ := strconv.ParseUint(info["signup"], 10, 64)
	if info != nil {
		return &domain.UserInfoResp{
			ID:        id,
			Username:  info["username"],
			Email:     info["email"],
			Followers: int32(followers),
			Following: int32(following),
			Posts:     int32(posts),
			Signup:    signup,
		}, nil
	}
	return nil, errors.New("get user info error")
}

func (u *User) PostStatus(ctx context.Context, userID int64, message string) (*domain.StatusInfo, error) {
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.HGet(ctx, common.UserInfoHashTable(userID), "username")
	pipeline.Incr(ctx, common.StatusIDCounter)
	res, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	username, id := res[0].(*redis.StringCmd).Val(), res[1].(*redis.IntCmd).Val()
	posted := time.Now().UnixNano()
	data := make(map[string]interface{})
	data["message"] = message
	data["posted"] = posted
	data["id"] = id
	data["userID"] = userID
	data["username"] = username

	pipeline.HMSet(ctx, common.StatusInfoHashTable(id), data)
	pipeline.HIncrBy(ctx, common.UserInfoHashTable(userID), "posts", 1)
	pipeline.ZAdd(ctx, common.UserProfileZSet(userID), redis.Z{Member: id, Score: float64(posted)})
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, err
	}
	u.syndicateStatus(ctx, userID, redis.Z{Score: float64(posted), Member: id})
	// fmt.Println(u.Data.Redis.HGetAll(ctx, common.StatusInfoHashTable(id)).Val())
	return &domain.StatusInfo{
		ID:       id,
		UserID:   userID,
		Username: username,
		Message:  message,
		Posted:   uint64(posted),
	}, nil
}

func (u *User) syndicateStatus(ctx context.Context, uid int64, post redis.Z) error {
	followers := u.Data.Redis.ZRangeByScoreWithScores(ctx, common.FollowerZSet(uid),
		&redis.ZRangeBy{Min: "0", Max: "inf"}).Val()

	pipeline := u.Data.Redis.TxPipeline()
	for _, z := range followers {
		fmt.Println(z.Member)
		follower := z.Member.(string)
		followerID, _ := strconv.ParseInt(follower, 10, 64)
		pipeline.ZAdd(ctx, common.HomeTimelineZSet(followerID), post)
		pipeline.ZRemRangeByRank(ctx, common.HomeTimelineZSet(followerID), 0, -common.HomeTimelineSize-1)
	}
	_, err := pipeline.Exec(ctx)
	return err

}

func (u *User) DeleteStatus(ctx context.Context, userID int64, postID int64) error {
	key := common.StatusInfoHashTable(postID)
	lock := u.Data.AcquireLockWithTimeout(ctx, key, 1, 100)
	if lock == "" {
		return errors.New("concurrency errors")
	}
	defer u.Data.ReleaseLock(ctx, key, lock)
	ownerID := u.Data.Redis.HGet(ctx, key, "userID").Val()

	// fmt.Println(ownerID)
	ownerIDNum, _ := strconv.ParseInt(ownerID, 10, 64)
	if ownerIDNum != userID {
		return errors.New("can't delete someone else's status")
	}
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.Del(ctx, key)
	pipeline.ZRem(ctx, common.UserProfileZSet(userID), postID)

	pipeline.HIncrBy(ctx, common.UserInfoHashTable(userID), "posts", -1)
	pipeline.ZRangeByScoreWithScores(ctx, common.FollowerZSet(userID),
		&redis.ZRangeBy{Min: "0", Max: "inf"})
	res, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	followers := res[3].(*redis.ZSliceCmd).Val()
	for _, z := range followers {
		fmt.Println(z.Member)
		follower := z.Member.(string)
		followerID, _ := strconv.ParseInt(follower, 10, 64)
		pipeline.ZRem(ctx, common.HomeTimelineZSet(followerID), postID)
	}
	pipeline.ZRem(ctx, common.HomeTimelineZSet(userID), postID)
	_, err = pipeline.Exec(ctx)
	return err
}
func (u *User) GetTimeline(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*domain.StatusInfo, error) {
	statuses := u.Data.Redis.ZRevRange(ctx, common.HomeTimelineZSet(userID),
		int64((pageID-1)*pageSize), int64(pageID*pageSize)-1).Val()
	pipeline := u.Data.Redis.TxPipeline()
	for _, id := range statuses {
		pipeline.HGetAll(ctx, common.StatusInfoStringHashTable(id))
	}
	res, err := pipeline.Exec(ctx)
	// fmt.Println(res)
	if err != nil {
		return nil, err
	}
	info := make([]*domain.StatusInfo, len(res))

	for i, val := range res {
		temp := val.(*redis.MapStringStringCmd).Val()
		id, _ := strconv.ParseInt(temp["id"], 10, 64)
		userID, _ := strconv.ParseInt(temp["userID"], 10, 64)
		posted, _ := strconv.ParseUint(temp["posted"], 10, 64)
		info[i] = &domain.StatusInfo{
			ID:       id,
			UserID:   userID,
			Username: temp["username"],
			Message:  temp["message"],
			Posted:   posted,
		}
	}
	return info, nil
}

func (u *User) FollowAction(ctx context.Context, userID int64, otherID int64) error {
	if u.Data.Redis.ZScore(ctx, common.FollowingZSet(userID), fmt.Sprintf("%d", otherID)).Val() != 0 {
		return errors.New("duplicate follow")
	}
	now := time.Now().UnixNano()

	pipeline := u.Data.Redis.TxPipeline()
	pipeline.ZAdd(ctx, common.FollowingZSet(userID), redis.Z{Member: otherID, Score: float64(now)})
	pipeline.ZAdd(ctx, common.FollowerZSet(otherID), redis.Z{Member: userID, Score: float64(now)})
	pipeline.ZRevRangeWithScores(ctx, common.UserProfileZSet(otherID), 0, common.HomeTimelineSize-1)
	res, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline error in follow action:%s %s", err, res)
	}
	following, followers, statusAndScore :=
		res[0].(*redis.IntCmd).Val(), res[1].(*redis.IntCmd).Val(), res[2].(*redis.ZSliceCmd).Val()
	fmt.Println(statusAndScore)
	pipeline.HIncrBy(ctx, common.UserInfoHashTable(userID), "following", following)
	pipeline.HIncrBy(ctx, common.UserInfoHashTable(otherID), "followers", followers)
	if len(statusAndScore) != 0 {
		pipeline.ZAdd(ctx, common.HomeTimelineZSet(userID), statusAndScore...)
	}
	pipeline.ZRemRangeByRank(ctx, common.HomeTimelineZSet(userID), 0, -common.HomeTimelineSize-1)

	if res, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline error in follow action:%s %s", err, res)
	}
	return nil
}

func (u *User) UnFollowAction(ctx context.Context, userID int64, otherID int64) error {
	if u.Data.Redis.ZScore(ctx, common.FollowingZSet(userID), fmt.Sprintf("%d", otherID)).Val() == 0 {
		return errors.New("duplicate unfollow")
	}

	pipeline := u.Data.Redis.TxPipeline()
	pipeline.ZRem(ctx, common.FollowingZSet(userID), fmt.Sprintf("%d", otherID))
	pipeline.ZRem(ctx, common.FollowerZSet(otherID), fmt.Sprintf("%d", userID))
	pipeline.ZRevRange(ctx, common.UserProfileZSet(otherID), 0, common.HomeTimelineSize-1)
	res, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline error in unfollow action:%s %s", err, res)
	}
	following, followers, status :=
		res[0].(*redis.IntCmd).Val(), res[1].(*redis.IntCmd).Val(), res[2].(*redis.StringSliceCmd).Val()

	pipeline.HIncrBy(ctx, common.UserInfoHashTable(userID), "following", -following)
	pipeline.HIncrBy(ctx, common.UserInfoHashTable(otherID), "followers", -followers)
	if len(status) != 0 {
		pipeline.ZRem(ctx, common.HomeTimelineZSet(userID), status)

	}
	//TODO:填充时间线

	if res, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline error in unfollow action:%s %s", err, res)
	}
	return nil
}
