package helpers

import (
	"goskeleton/app/global/variable"
	"time"
)

// GetDatetime 获取当前日期， 默认格式
func GetDatetime() string {
	return time.Now().Format(variable.DateFormat)
}

// GetTimestamp 获取时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// TimestampToDatetime 时间戳转日志， 默认格式
func TimestampToDatetime(timestamp int64) string {
	if timestamp == 0 {
		timestamp = GetTimestamp()
	}
	t := time.Unix(timestamp, 0).Local()

	return t.Format(variable.DateFormat)
}
