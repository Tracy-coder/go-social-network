package common

import "fmt"

const (
	UserIDCounter   string = "user:id:"
	Username2ID     string = "users:"
	StatusIDCounter string = "status:id:"
	ChatIDCounter   string = "chat:id:"
)

const (
	HomeTimelineSize = 100
)

func NewKeyGenerator(prefix string) func(interface{}) string {
	return func(id interface{}) string {
		switch v := id.(type) {
		case int64:
			return fmt.Sprintf("%s:%d", prefix, v)
		case string:
			return fmt.Sprintf("%s:%s", prefix, v)
		default:
			panic(fmt.Sprintf("Unsupported type for id: %T", id))
		}
	}
}

var UserInfoHashTable func(interface{}) string = NewKeyGenerator("user")
var StatusInfoHashTable func(interface{}) string = NewKeyGenerator("status")
var UserProfileZSet func(interface{}) string = NewKeyGenerator("profile")
var HomeTimelineZSet func(interface{}) string = NewKeyGenerator("home")
var FollowerZSet func(interface{}) string = NewKeyGenerator("followers")
var FollowingZSet func(interface{}) string = NewKeyGenerator("following")
var ChatMembersZSet func(interface{}) string = NewKeyGenerator("chat")
var UserLastSeenZset func(interface{}) string = NewKeyGenerator("seen")
var MessageIDCounter func(interface{}) string = NewKeyGenerator("ids")
var MessageInChatZset func(interface{}) string = NewKeyGenerator("msgs")
