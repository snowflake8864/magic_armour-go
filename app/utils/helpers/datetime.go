package helpers

import (
	"goskeleton/app/global/variable"
	"time"
)

// GetDatetime 获取当前日期， 默认格式
func GetDatetime() string {
	return time.Now().Format(variable.DateFormat)
}

// GetTimestamp 获取当前时间戳
func GetTimestamp() int64 {
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
		timestamp = GetTimestamp()
	}
	t := time.Unix(timestamp, 0).Local()

	return t.Format(variable.DateFormat)
}
