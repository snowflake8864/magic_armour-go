package cur_app

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/string_helper"
)

// GetAppUsageVersion 获取当前程序使用版本(离线|云版)
func GetAppUsageVersion() int {
	return variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Version")
}

// GetCustomVersionNumber 获取自定义版本号
func GetCustomVersionNumber() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.VersionNumber")
}

// GetPlatformConfig  获取平台自定义配置
func GetPlatformConfig(lang string) map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = GetAppUsageVersion()
	data["ver_number"] = variable.VersionNumber
	if ver := GetCustomVersionNumber(); !string_helper.IsEmpty(ver) {
		data["ver_number"] = ver
	}
	data["title"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.Title")
	data["iconShow"] = variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Platform.IconShow")
	data["projectName"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.ProjectName")
	data["loginTip"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.LoginTip")
	data["qrShow"] = variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Platform.QrShow")
	data["orgAddr"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.Address")
	data["website"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.Website")
	data["tel"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.Tel")
	data["copyright"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.Copyright")
	data["license"] = variable.ConfigCustomV2Yml.GetStringSlice("AppCustomConfig.Platform.License")

	if lang == consts.LanguageEn {
		data["title"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.TitleEn")
		data["projectName"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.ProjectNameEn")
		data["loginTip"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.LoginTipEn")
		data["orgAddr"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.AddressEn")
		data["copyright"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Platform.CopyrightEn")
		data["license"] = variable.ConfigCustomV2Yml.GetStringSlice("AppCustomConfig.Platform.LicenseEn")
	}

	return data
}

// GetOrganizationConfig  获取企业自定义配置
func GetOrganizationConfig(lang string) map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = GetAppUsageVersion()
	data["ver_number"] = variable.VersionNumber
	if ver := GetCustomVersionNumber(); !string_helper.IsEmpty(ver) {
		data["ver_number"] = ver
	}

	data["title"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.Title")
	data["iconShow"] = variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Organization.IconShow")
	data["projectName"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.ProjectName")
	data["loginTip"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.LoginTip")
	data["qrShow"] = variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Organization.QrShow")
	data["orgAddr"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.Address")
	data["website"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.Website")
	data["tel"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.Tel")
	data["copyright"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.Copyright")
	data["license"] = variable.ConfigCustomV2Yml.GetStringSlice("AppCustomConfig.Organization.License")

	if lang == consts.LanguageEn {
		data["title"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.TitleEn")
		data["projectName"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.ProjectNameEn")
		data["loginTip"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.LoginTipEn")
		data["orgAddr"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.AddressEn")
		data["copyright"] = variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Organization.CopyrightEn")
		data["license"] = variable.ConfigCustomV2Yml.GetStringSlice("AppCustomConfig.Organization.LicenseEn")
	}

	return data
}

// GetPlatformTokenExpire  获取平台口令过期时间
func GetPlatformTokenExpire() int {
	return variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Platform.TokenExpire")
}

// GetOrganizationTokenExpire  获取企业口令过期时间
func GetOrganizationTokenExpire() int {
	return variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Organization.TokenExpire")
}

// GetAgentTokenExpire  获取终端口令过期时间
func GetAgentTokenExpire() int {
	return variable.ConfigCustomV2Yml.GetInt("AppCustomConfig.Agent.TokenExpire")
}

// GetServiceCode 获取服务器编号：
func GetServiceCode() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.ServiceMonitor")
}

// ShellNetworkGet 网络配置脚本： 获取
func ShellNetworkGet() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.NetworkGet")
}

// ShellNetworkSet 网络配置脚本： 设置
func ShellNetworkSet() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.NetworkSet")
}

// ShellServiceMonitor 服务器监控脚本：
func ShellServiceMonitor() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.ServiceMonitor")
}

// ShellUpgrade 升级脚本：
func ShellUpgrade() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.Upgrade")
}

// ShellInit 脚本初始化命令：
func ShellInit() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.Init")
}

// ShellVedaLicence 卫达授权证书解析：
func ShellVedaLicence() string {
	return variable.ConfigCustomV2Yml.GetString("AppCustomConfig.Shell.VedaLicence")
}
