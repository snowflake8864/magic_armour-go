package helpers

import (
	"goskeleton/app/global/variable"
	"time"
)

// GetNowDate 获取当前日期， 默认格式
func GetNowDate() string {
	return time.Now().Format(variable.DateFormat)
}
