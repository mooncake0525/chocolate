package utils

import (
	"github.com/bwmarrin/snowflake"
)

/*
@Author : VictorTu
@Software: GoLand
*/

var SnowFlakeNode *snowflake.Node

func InitSnowFlake(nodeId int64) (err error) {
	if SnowFlakeNode, err = snowflake.NewNode(nodeId); err != nil {
		return err
	}
	return nil
}
