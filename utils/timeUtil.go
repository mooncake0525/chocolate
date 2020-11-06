package utils

import (
	"fmt"
	"time"
)

/*
@Author : VictorTu
@Software: GoLand
*/

type timeUtil struct {
}

var TimeUtil timeUtil

const (
	Layout_Date_Time9      = "2006-01-02 15:04:05.999999999"
	Layout_Date_Time6      = "2006-01-02 15:04:05.999999"
	Layout_Date_Time       = "2006-01-02 15:04:05"
	Layout_Date            = "2006-01-02"
	Layout_Date0           = "20060102"
	Time_Layout_Date_Time2 = "20060102150405"
)

func (t *timeUtil) FormatTime(layout string, d time.Time) (dString string) {
	if !d.IsZero() {
		dString = d.Format(layout)
	}
	return
}

func (t *timeUtil) GetDayBeginTime(d time.Time) (nd time.Time) {
	nd, _ = time.ParseInLocation(Layout_Date_Time, fmt.Sprintf("%s 00:00:00", d.Format(Layout_Date)), time.Local)
	return
}
