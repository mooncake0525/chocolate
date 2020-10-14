package utils

import (
	"stream/models/constants"
	"time"
)

/*
@Author : VictorTu
@Software: GoLand
*/

type timeUtil struct {
}

var TimeUtil timeUtil

func (this *timeUtil) GetNowTimeStamp() string {
	return time.Now().Format(constants.Layout_Date_Time)
}
