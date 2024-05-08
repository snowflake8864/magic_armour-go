package helpers

import (
	"goskeleton/app/global/variable"
	"time"
)

// DatetimeCurr 获取当前日期， 默认格式
func DatetimeCurr() string {
	return time.Now().Format(variable.DateFormat)
}

// DatetimeSpecial 获取当前日期，特殊格式
func DatetimeSpecial(format string) string {
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

// DateCalcByDay 按天计算日期
func DateCalcByDay(d int) (day, dayStart, dayEnd string, month int) {
	t := time.Now().AddDate(0, 0, d)
	day = t.Format(variable.DayFormat)
	dayStart = day + " 00:00:00"
	dayEnd = day + " 23:59:59"
	month = String2Int(t.Format("200601"))
	return day, dayStart, dayEnd, month
}

// DateCalcByMonth 按月计算日期
func DateCalcByMonth(m int) (month int, monthFirstDay, monthLastDay string) {
	t := time.Now().AddDate(0, m, 0)
	month = String2Int(t.Format("200601"))
	ms := t.Format("2006-01")
	monthFirstDay = ms + " -01"
	monthLastDay = ms + " -" + IntToStr(DaysInMonth(t.Year(), t.Month()))
	return month, monthFirstDay, monthLastDay
}

// DaysInMonth 获取指定年月一共有多少天
func DaysInMonth(year int, month time.Month) int {
	// 创建一个表示该年指定月份的第一天的时间
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 下个月的第一天减去一天就是这个月的最后一天
	return t.AddDate(0, 1, -1).Day()
}

// TimestampCurr 获取当前时间戳
func TimestampCurr() int64 {
	return time.Now().Unix()
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

// TimestampToDatetime 时间戳转成日期， 默认格式
func TimestampToDatetime(timestamp int64) string {
	if timestamp == 0 {
		timestamp = TimestampCurr()
	}
	t := time.Unix(timestamp, 0).Local()

	return t.Format(variable.DateFormat)
}
