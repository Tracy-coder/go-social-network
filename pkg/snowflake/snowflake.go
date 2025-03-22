package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var defaultNode *snowflake.Node

func InitSnowFlake() {
	var err error
	defaultNode, err = snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
}
func GenerateID() string {
	return defaultNode.Generate().String()
}
