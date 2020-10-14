package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"testing"
)

/*
@Author : VictorTu
@Software: GoLand
*/

func TestInitSnowFlake(t *testing.T) {
	if err := InitSnowFlake(); err != nil {
		t.Error(err.Error())
		return
	}
	m := make(map[snowflake.ID]snowflake.ID)
	for i := 0; i < 1000000; i++ {
		id := SnowFlakeNode.Generate()
		if _, ok := m[id]; ok {
			fmt.Println("repeated ", id)
		} else {
			m[id] = id
		}
	}
}
