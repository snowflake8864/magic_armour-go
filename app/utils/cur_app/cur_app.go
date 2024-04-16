package cur_app

import "goskeleton/app/global/variable"

// GetAppUsageVersion 获取当前程序使用版本
func GetAppUsageVersion() int {
	return variable.ConfigYml.GetInt("AppUsageVersion")
}
