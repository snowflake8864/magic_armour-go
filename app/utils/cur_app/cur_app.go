package cur_app

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
)

// GetAppUsageVersion 获取当前程序使用版本
func GetAppUsageVersion() int {
	return variable.ConfigYml.GetInt("AppCustomConfig.Version")
}

// GetPlatformConfig  获取平台自定义配置
func GetPlatformConfig() map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = GetAppUsageVersion()
	data["ver_number"] = consts.VersionNumber

	data["title"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.Title")
	data["iconShow"] = variable.ConfigYml.GetInt("AppCustomConfig.Platform.IconShow")
	data["projectName"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.ProjectName")
	data["loginTip"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.LoginTip")
	data["qrShow"] = variable.ConfigYml.GetInt("AppCustomConfig.Platform.QrShow")
	data["address"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.Address")
	data["website"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.Website")
	data["tel"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.Tel")
	data["copyright"] = variable.ConfigYml.GetString("AppCustomConfig.Platform.Copyright")
	data["license"] = variable.ConfigYml.GetStringSlice("AppCustomConfig.Platform.License")

	return data
}

// GetOrganizationConfig  获取企业自定义配置
func GetOrganizationConfig() map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = GetAppUsageVersion()
	data["ver_number"] = consts.VersionNumber

	data["title"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.Title")
	data["iconShow"] = variable.ConfigYml.GetInt("AppCustomConfig.Organization.IconShow")
	data["projectName"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.ProjectName")
	data["loginTip"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.LoginTip")
	data["qrShow"] = variable.ConfigYml.GetInt("AppCustomConfig.Organization.QrShow")
	data["address"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.Address")
	data["website"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.Website")
	data["tel"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.Tel")
	data["copyright"] = variable.ConfigYml.GetString("AppCustomConfig.Organization.Copyright")
	data["license"] = variable.ConfigYml.GetStringSlice("AppCustomConfig.Organization.License")

	return data
}

// GetPlatformTokenExpire  获取平台口令过期时间
func GetPlatformTokenExpire() int {
	return variable.ConfigYml.GetInt("AppCustomConfig.Platform.TokenExpire")
}

// GetOrganizationTokenExpire  获取企业口令过期时间
func GetOrganizationTokenExpire() int {
	return variable.ConfigYml.GetInt("AppCustomConfig.Organization.TokenExpire")
}

// GetAgentTokenExpire  获取终端口令过期时间
func GetAgentTokenExpire() int {
	return variable.ConfigYml.GetInt("AppCustomConfig.Agent.TokenExpire")
}

// ShellNetworkGet 网络配置脚本： 获取
func ShellNetworkGet() string {
	return variable.ConfigYml.GetString("AppCustomConfig.Shell.NetworkGet")
}

// ShellNetworkSet 网络配置脚本： 设置
func ShellNetworkSet() string {
	return variable.ConfigYml.GetString("AppCustomConfig.Shell.NetworkSet")
}
