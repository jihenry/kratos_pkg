package util

import "time"

func GetTodayZeroTime() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func GetDayZeroTime(day int) time.Time {
	zero := GetTodayZeroTime()
	return zero.AddDate(0, 0, day)
}
