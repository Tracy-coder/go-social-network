package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-social-network/biz/common"
	"go-social-network/biz/domain"
	"go-social-network/configs"
	"go-social-network/data"
	"go-social-network/pkg/encrypt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
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
	u.syndicateStatus(userID, float64(posted), id, 1)
	// fmt.Println(u.Data.Redis.HGetAll(ctx, common.StatusInfoHashTable(id)).Val())
	return &domain.StatusInfo{
		ID:       id,
		UserID:   userID,
		Username: username,
		Message:  message,
		Posted:   uint64(posted),
	}, nil
}

func (u *User) syndicateStatus(uid int64, posted float64, statusID int64, op int64) error {
	msg := StatusReqInKafka{
		UserID:   uid,
		Score:    posted,
		StatusID: statusID,
		Op:       op,
	}
	json_msg, _ := json.Marshal(msg)
	_, _, err := SyndicateStatusKafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: configs.Data().Kafka.Topic,
		Key:   sarama.StringEncoder(json_msg),
		Value: sarama.StringEncoder(json_msg),
	})
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

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	u.syndicateStatus(userID, float64(postID), 0, 2)
	return err
}
func (u *User) GetTimeline(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*domain.StatusInfo, error) {
	statuses := u.Data.Redis.ZRevRange(ctx, common.HomeTimelineZSet(userID),
		int64((pageID-1)*pageSize), int64(pageID*pageSize)-1).Val()
	pipeline := u.Data.Redis.TxPipeline()
	for _, id := range statuses {
		pipeline.HGetAll(ctx, common.StatusInfoHashTable(id))
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

func (u *User) GetProfile(ctx context.Context, userID int64, pageID int32, pageSize int32) ([]*domain.StatusInfo, error) {
	statuses := u.Data.Redis.ZRevRange(ctx, common.UserProfileZSet(userID),
		int64((pageID-1)*pageSize), int64(pageID*pageSize)-1).Val()
	pipeline := u.Data.Redis.TxPipeline()
	for _, id := range statuses {
		pipeline.HGetAll(ctx, common.StatusInfoHashTable(id))
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
	// fmt.Println(statusAndScore)
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

	if res, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline error in unfollow action:%s %s", err, res)
	}
	// return u.refillTimeline(ctx, userID)
	return nil
}

func (u *User) CreateChat(ctx context.Context, ownerID int64, membersID []int64) (int64, error) {
	chatID := int64(u.Data.Redis.Incr(ctx, common.ChatIDCounter).Val())
	membersID = append(membersID, ownerID)
	var recipientsd []redis.Z
	for _, r := range membersID {
		temp := redis.Z{
			Score:  0,
			Member: r,
		}
		recipientsd = append(recipientsd, temp)
	}

	pipeline := u.Data.Redis.TxPipeline()
	pipeline.ZAdd(ctx, common.ChatMembersZSet(chatID), recipientsd...)
	for _, id := range membersID {
		pipeline.ZAdd(ctx, common.UserLastSeenZset(id), redis.Z{Member: chatID, Score: 0})
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		return -1, fmt.Errorf("pipeline err in CreateChat: %s", err)
	}
	return chatID, nil
}

func (u *User) PostMessage(ctx context.Context, userID int64, chatID int64, message string) (*domain.MessageInfo, error) {
	_, err := u.Data.Redis.ZScore(ctx, common.ChatMembersZSet(chatID), fmt.Sprintf("%d", userID)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("can't send messages in that group")
	}
	messageID := u.Data.Redis.Incr(ctx, common.MessageIDCounter(chatID)).Val()
	ts := time.Now().UnixNano()
	packed := domain.MessageInfo{
		ID:        messageID,
		CreatedAt: uint64(ts),
		Content:   message,
		SenderID:  userID,
		ChatID:    chatID,
	}

	jsonValue, err := json.Marshal(packed)
	if err != nil {
		log.Println("marshal err in SendMessage: ", err)
	}
	fmt.Println(jsonValue, messageID, chatID)
	u.Data.Redis.ZAdd(ctx, common.MessageInChatZset(chatID), redis.Z{Member: jsonValue, Score: float64(messageID)}).Val()
	return &domain.MessageInfo{
		ID:        messageID,
		CreatedAt: uint64(ts),
		Content:   message,
		SenderID:  userID,
		ChatID:    chatID,
	}, nil
}

func (u *User) GetPendingMessage(ctx context.Context, userID int64) (*[]domain.ChatMessageInfo, error) {
	seen := u.Data.Redis.ZRangeWithScores(ctx, common.UserLastSeenZset(userID), 0, -1).Val()
	pipeline := u.Data.Redis.TxPipeline()
	res := []*redis.StringSliceCmd{}
	length := len(seen)
	temp := make([]string, 0, length)

	for _, v := range seen {
		chatID := v.Member.(string)
		seenID := v.Score
		res = append(res, pipeline.ZRangeByScore(ctx, common.MessageInChatZset(chatID), &redis.ZRangeBy{Min: strconv.Itoa(int(seenID) + 1), Max: "inf"}))
		// fmt.Println(res)
		temp = append(temp, chatID)
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, fmt.Errorf("Pipeline error in GetPendingMessage:%s", err)
	}
	fmt.Println(res)

	result := make([]domain.ChatMessageInfo, len(temp))

	for i := 0; i < len(temp); i++ {
		if len(res[i].Val()) == 0 {
			continue
		}
		messages := []domain.MessageInfo{}
		for _, v := range res[i].Val() {
			message := domain.MessageInfo{}
			if err := json.Unmarshal([]byte(v), &message); err != nil {
				return nil, fmt.Errorf("Unmarshal error in GetPendingMessage:%s", err)
			}
			messages = append(messages, message)
		}

		chatID := temp[i]
		seenID := float64(messages[len(messages)-1].ID)
		u.Data.Redis.ZAdd(ctx, common.ChatMembersZSet(chatID), redis.Z{Member: userID, Score: seenID})

		minID := u.Data.Redis.ZRangeWithScores(ctx, common.ChatMembersZSet(chatID), 0, 0).Val()

		u.Data.Redis.ZAdd(ctx, common.UserLastSeenZset(userID), redis.Z{Member: chatID, Score: seenID})
		if minID != nil {
			u.Data.Redis.ZRemRangeByScore(ctx, common.MessageInChatZset(chatID), "0", strconv.Itoa(int(minID[0].Score)))
		}
		result[i].Info = messages
	}
	// fmt.Println("result:", result)
	return &result, nil
}

func (u *User) LeaveChat(ctx context.Context, userID, chatID int64) error {
	_, err := u.Data.Redis.ZScore(ctx, common.ChatMembersZSet(chatID), fmt.Sprintf("%d", userID)).Result()
	if err == redis.Nil {
		return fmt.Errorf("You've left the group")
	}
	pipeline := u.Data.Redis.TxPipeline()
	pipeline.ZRem(ctx, common.ChatMembersZSet(chatID), userID)
	pipeline.ZRem(ctx, common.UserLastSeenZset(userID), chatID)
	res := pipeline.ZCard(ctx, common.ChatMembersZSet(chatID)).Val()
	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline err in LeaveChat: %s", err)
	}

	if res == 0 {
		pipeline.Del(ctx, common.MessageInChatZset(chatID))
		pipeline.Del(ctx, common.MessageIDCounter(chatID))
		pipeline.Del(ctx, common.ChatMembersZSet(chatID))
		if _, err := pipeline.Exec(ctx); err != nil {
			return fmt.Errorf("pipeline err in LeaveChat: %s", err)
		}
	} else {
		oldest := u.Data.Redis.ZRangeWithScores(ctx, common.ChatMembersZSet(chatID), 0, 0).Val()[0]
		u.Data.Redis.ZRemRangeByScore(ctx, common.MessageInChatZset(chatID), "0", strconv.Itoa(int(oldest.Score)))
	}
	return nil
}
