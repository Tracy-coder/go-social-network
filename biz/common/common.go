package common

import "fmt"

const (
	UserIDCounter string = "user:id:"
	Username2ID   string = "users:"
)

func UserInfoHashTable(id int64) string {
	return fmt.Sprintf("user:%d", id)
}
