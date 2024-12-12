package common

import "fmt"

const (
	UserIDCounter   string = "user:id:"
	Username2ID     string = "users:"
	StatusIDCounter string = "status:id:"
)

const (
	HomeTimelineSize = 100
)

func UserInfoHashTable(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func StatusInfoHashTable(id int64) string {
	return fmt.Sprintf("status:%d", id)
}

func StatusInfoStringHashTable(id string) string {
	return fmt.Sprintf("status:%s", id)
}

func UserProfileZSet(id int64) string {
	return fmt.Sprintf("profile:%d", id)
}

func UserProfileStringZSet(id string) string {
	return fmt.Sprintf("profile:%s", id)
}

func HomeTimelineZSet(id int64) string {
	return fmt.Sprintf("home:%d", id)
}

func HomeTimelineStringZSet(id string) string {
	return fmt.Sprintf("home:%s", id)
}

func FollowerZSet(id int64) string {
	return fmt.Sprintf("followers:%d", id)
}

func FollowingZSet(id int64) string {
	return fmt.Sprintf("following:%d", id)
}
