package utils

import "time"

/*
@Author : VictorTu
@Software: GoLand
*/

type timeUtil struct {
}

var TimeUtil timeUtil

func (t *timeUtil) FormatTime(layout string, d time.Time) (dString string) {
	if !d.IsZero() {
		dString = d.Format(layout)
	}
	return
}
