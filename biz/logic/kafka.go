package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"go-social-network/biz/common"
	"go-social-network/configs"
	"go-social-network/data"
	"go-social-network/pkg/kafka"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type syndicateStatusConsumerGroup struct {
	Data *data.Data
}

var group syndicateStatusConsumerGroup
var SyndicateStatusKafkaProducer sarama.SyncProducer
var SyndicateStatusKafkaConsumer sarama.ConsumerGroup

func init() {
	SyndicateStatusKafkaProducer = kafka.InitSynProducer(configs.Data().Kafka.Brokers)
	SyndicateStatusKafkaConsumer = kafka.InitConsumerGroup(configs.Data().Kafka.Brokers, configs.Data().Kafka.GroupID)
	group = syndicateStatusConsumerGroup{
		Data: data.Default(),
	}
}
func (m syndicateStatusConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (m syndicateStatusConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

type StatusReqInKafka struct {
	UserID   int64
	Score    float64
	StatusID int64
	Op       int64 // 为1表示添加 为2表示删除
}

func (m syndicateStatusConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d  value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		req := StatusReqInKafka{}
		json.Unmarshal(msg.Value, &req)
		fmt.Println(req)
		followers := m.Data.Redis.ZRangeByScoreWithScores(sess.Context(), common.FollowerZSet(req.UserID),
			&redis.ZRangeBy{Min: "0", Max: "inf"}).Val()

		pipeline := m.Data.Redis.TxPipeline()
		for _, z := range followers {
			// fmt.Println(z.Member)
			follower := z.Member.(string)
			followerID, _ := strconv.ParseInt(follower, 10, 64)
			if req.Op == 1 {
				pipeline.ZAdd(sess.Context(), common.HomeTimelineZSet(followerID), redis.Z{Score: req.Score, Member: req.StatusID})
				// todo: 把remove写成一个定时的后台协程
				pipeline.ZRemRangeByRank(sess.Context(), common.HomeTimelineZSet(followerID), 0, -common.HomeTimelineSize-1)
			} else {
				pipeline.ZRem(sess.Context(), common.HomeTimelineZSet(followerID), req.StatusID)
			}

		}
		_, err := pipeline.Exec(sess.Context())
		if err != nil {
			fmt.Printf("failed to syndicate status: %s", err)
		}
		// 标记，sarama会自动进行提交，默认间隔1秒
		sess.MarkMessage(msg, "")
	}
	return nil
}

func Consumer4SyndicateStatus() {
	for {
		fmt.Println("for loop!")
		err := SyndicateStatusKafkaConsumer.Consume(context.Background(), []string{"email"}, group)

		if err != nil {
			break
		}

	}

	_ = SyndicateStatusKafkaConsumer.Close()
}
