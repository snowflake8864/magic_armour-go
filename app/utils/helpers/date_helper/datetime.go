package date_helper

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"time"
)

// CurDatetime 获取当前日期， 默认格式
func CurDatetime() string {
	return time.Now().Format(variable.DateFormat)
}

// SpecialDatetime 获取当前日期，特殊格式
func SpecialDatetime(format string) string {
	return time.Now().Format(format)
}

// DayInCurrMonth 现在是当前月的第几天
func DayInCurrMonth() int {
	return time.Now().Day()
}

// WeekdayCurr 现在是周几
func WeekdayCurr() string {
	return time.Now().Weekday().String()
}

// CalcByDay 按天计算日期
func CalcByDay(d int) (day, dayStart, dayEnd string, month int) {
	t := time.Now().AddDate(0, 0, d)
	day = t.Format(variable.DayFormat)
	dayStart = day + " 00:00:00"
	dayEnd = day + " 23:59:59"
	month = string_helper.ToInt(t.Format("200601"))
	return day, dayStart, dayEnd, month
}

// CalcByMonth 按月计算日期
func CalcByMonth(m int) (month int, monthFirstDay, monthLastDay string) {
	t := time.Now().AddDate(0, m, 0)
	month = string_helper.ToInt(t.Format("200601"))
	ms := t.Format("2006-01")
	monthFirstDay = ms + " -01"
	monthLastDay = ms + " -" + int_helper.IntToStr(DaysInMonth(t.Year(), t.Month()))
	return month, monthFirstDay, monthLastDay
}

// DaysInMonth 获取指定年月一共有多少天
func DaysInMonth(year int, month time.Month) int {
	// 创建一个表示该年指定月份的第一天的时间
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 下个月的第一天减去一天就是这个月的最后一天
	return t.AddDate(0, 1, -1).Day()
}

// CurTimestamp 获取当前时间戳
func CurTimestamp() int64 {
	return time.Now().Unix()
}

// AddTimestampByDays 时间戳按照年月日时增加
func AddTimestampByDays(years, months, days int) int64 {
	now := time.Now()
	return now.AddDate(years, months, days).Unix()
}

// AddTimestampByHour 时间戳按照小时增加
func AddTimestampByHour(h int64) int64 {
	now := time.Now()
	return now.Add(time.Duration(h) * time.Hour).Unix()
}

// AddTimestampByMinute 时间戳按照分钟增加
func AddTimestampByMinute(m int64) int64 {
	now := time.Now()
	return now.Add(time.Duration(m) * time.Minute).Unix()
}

// CalcDatetimeByMinute 按分钟把给定时间从新计算
func CalcDatetimeByMinute(t string, m int64) string {
	timestamp := DatetimeToTimestamp(t, variable.DateFormat)
	timestamp = AddTimestampByMinute(m)
	return TimestampToDatetime(timestamp)
}

// TimestampToDatetime 时间戳转成日期， 默认格式
func TimestampToDatetime(timestamp int64) string {
	if timestamp == 0 {
		timestamp = CurTimestamp()
	}
	t := time.Unix(timestamp, 0).Local()

	return t.Format(variable.DateFormat)
}

// TimestampToDatetimeByFormat 时间戳转成日期， 指定格式
func TimestampToDatetimeByFormat(timestamp int64, format string) string {
	if timestamp == 0 {
		timestamp = CurTimestamp()
	}
	t := time.Unix(timestamp, 0).Local()

	return t.Format(format)
}

// DatetimeToTimestamp 日期转成时间戳
func DatetimeToTimestamp(datetime, layout string) int64 {
	// 解析日期字符串为time.Time类型
	t, err := time.Parse(layout, datetime)
	if err != nil {
		return 0
	}
	return t.Local().Add(-8 * time.Hour).Unix()
}

// DaysInMonthByTimestamp 时间戳现在是月的第几天
func DaysInMonthByTimestamp(timestamp int64) int {
	return time.Unix(timestamp, 0).Local().Day()
}
