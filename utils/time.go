package utils

import (
	"math"
	"strings"
	"time"
)

// GetZeroTime 获取当天0点时间
func GetZeroTime(dayOffset int) time.Time {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day() + dayOffset
	zeroTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return zeroTime
}

// GetZeroTimeOf 获取指定日期的0点时间
func GetZeroTimeOf(t time.Time) time.Time {
	year := t.Year()
	month := t.Month()
	day := t.Day()
	zeroTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return zeroTime
}

// TimeFormat 日期格式化: yyyy-MM-dd HH:mm:ss
func TimeFormat(layout string, dateTime time.Time) string {
	converts := map[string]string{
		"yyyy": "2006",
		"MM":   "01",
		"dd":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
	}
	for k, v := range converts {
		layout = strings.Replace(layout, k, v, 1)
	}
	return dateTime.Format(layout)
}

// ToTime 时间字符串转 Time
func ToTime(layout, value string) *time.Time {
	converts := map[string]string{
		"yyyy": "2006",
		"MM":   "01",
		"dd":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
	}
	for k, v := range converts {
		layout = strings.Replace(layout, k, v, 1)
	}

	dateTime, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		panic(err)
	}
	return &dateTime
}

// GetWeekInYear 判断时间是当年的第几周
func GetWeekInYear(t time.Time) int {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}

	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = int(math.Ceil(float64(yearDay-firstWeekDays)/7) + 1)
	}
	return week
}
