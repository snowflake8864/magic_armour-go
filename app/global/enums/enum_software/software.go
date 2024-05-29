package enum_software

import (
	"goskeleton/app/global/enums"
	"strings"
)

const (
	PlatformTypeLinux   = 1 //平台类型: linux
	PlatformTypeWindows = 2 //平台类型: windows

	OsTypeLinuxUbuntu      = 101 //系统类型: linux Ubuntu
	OsTypeLinuxCentos      = 102 //系统类型: linux Centos
	OsTypeLinuxNeokylin    = 103 //系统类型: linux 麒麟
	OsTypeLinuxUos         = 104 //系统类型: linux 通信
	OsTypeLinuxRedhat      = 105 //系统类型: linux 红帽
	OsTypeLinuxDebian      = 106 //系统类型: linux Debian
	OsTypeLinuxRocky       = 107 //系统类型: linux Rocky
	OsTypeLinuxAnolis      = 108 //系统类型: linux 阿里龙蜥
	OsTypeLinuxOpencloudos = 109 //系统类型: linux 腾讯

	OsTypeWinSeven = 201 //系统类型: win7以上
	OsTypeWinXp    = 202 //系统类型: win xp 和 win 1003

)

var EnumMap = map[string]map[int]string{
	"platform_type": {
		PlatformTypeWindows: "Windows",
		PlatformTypeLinux:   "Linux",
	},
	"platform_type_alias": {
		PlatformTypeWindows: "WIN",
		PlatformTypeLinux:   "LINUX",
	},
	"os_type": {
		OsTypeWinSeven:         "Windows",
		OsTypeWinXp:            "Windows早期",
		OsTypeLinuxUbuntu:      "Ubuntu",
		OsTypeLinuxCentos:      "CentOS",
		OsTypeLinuxNeokylin:    "Neokylin",
		OsTypeLinuxUos:         "Uos",
		OsTypeLinuxRedhat:      "RedHat",
		OsTypeLinuxDebian:      "Debian",
		OsTypeLinuxRocky:       "Rocky",
		OsTypeLinuxAnolis:      "Anolis",
		OsTypeLinuxOpencloudos: "OpenCloudOs",
	},
	"os_type_alias": {
		OsTypeWinSeven:         "windows7-至今",
		OsTypeWinXp:            "windows xp-windows 2003",
		OsTypeLinuxUbuntu:      "Ubuntu64位",
		OsTypeLinuxCentos:      "CentOS64位",
		OsTypeLinuxNeokylin:    "Neokylin64位",
		OsTypeLinuxUos:         "Uos64位",
		OsTypeLinuxRedhat:      "RedHat64位",
		OsTypeLinuxDebian:      "Debian64位",
		OsTypeLinuxRocky:       "Rocky64位",
		OsTypeLinuxAnolis:      "Anolis64位",
		OsTypeLinuxOpencloudos: "OpenCloudOs64位",
	},
	"os_ver": {
		//Ubuntu
		1010: "64位(14.04-22.04)",
		//CentOS
		1020: "64位(6.5 - 8.5)",
		//麒麟
		1030: "64位 中标麒麟(V7)",
		//统信
		1040: "64位 (20)",
		//RedHat
		1050: "64位(6.5 - 7.9)",
		//Debian
		1060: "64位(14.04 - 22.04)",
		//Rocky
		1070: "64位(8.4 - 9.1)",
		//ANOLIS
		1080: "64位(7.7 - 8.8)",
		//OPENCLOUDOS
		1090: "64位(8.6-9.0)",

		//windows
		2010: "windows7-至今",
		//Windows其它
		2020: "windows xp-windows 2003",
	},
}
var EnumSlice = map[string]map[int][]int{
	"map_platform_os": {
		PlatformTypeWindows: {
			OsTypeWinSeven,
			OsTypeWinXp,
		},
		PlatformTypeLinux: {
			OsTypeLinuxUbuntu,
			OsTypeLinuxCentos,
			OsTypeLinuxNeokylin,
			OsTypeLinuxUos,
			OsTypeLinuxRedhat,
			OsTypeLinuxDebian,
			OsTypeLinuxRocky,
			OsTypeLinuxAnolis,
			OsTypeLinuxOpencloudos,
		},
	},
	"map_os_ver": {
		OsTypeWinSeven: {
			2010,
		},
		OsTypeWinXp: {
			2020,
		},
		OsTypeLinuxUbuntu: {
			1010,
		},
		OsTypeLinuxCentos: {
			1020,
		},
		OsTypeLinuxNeokylin: {
			1030,
		},
		OsTypeLinuxUos: {
			1040,
		},
		OsTypeLinuxRedhat: {
			1050,
		},
		OsTypeLinuxDebian: {
			1060,
		},
		OsTypeLinuxRocky: {
			1070,
		},
		OsTypeLinuxAnolis: {
			1080,
		},
		OsTypeLinuxOpencloudos: {
			1090,
		},
	},
}

var WebMapInt = map[string][]int{
	"platform_type": {
		PlatformTypeWindows,
		PlatformTypeLinux,
	},
}

// GetOsTypeByName 获取系统类型
func GetOsTypeByName(name string) int {
	name = strings.ToLower(name)
	if strings.Contains(name, "windows") {
		if strings.Contains(name, "xp") || strings.Contains(name, "2003") {
			return OsTypeWinSeven
		} else {
			return OsTypeWinXp
		}
	} else if strings.Contains(name, "ubuntu") {
		return OsTypeLinuxUbuntu
	} else if strings.Contains(name, "centos") {
		return OsTypeLinuxCentos
	} else if strings.Contains(name, "neokylin") || strings.Contains(name, "麒麟") {
		return OsTypeLinuxNeokylin
	} else if strings.Contains(name, "uos") || strings.Contains(name, "统信") {
		return OsTypeLinuxUos
	} else if strings.Contains(name, "redhat") || strings.Contains(name, "red hat") {
		return OsTypeLinuxRedhat
	} else if strings.Contains(name, "debian") {
		return OsTypeLinuxDebian
	} else if strings.Contains(name, "rocky") {
		return OsTypeLinuxRocky
	} else if strings.Contains(name, "anolis") {
		return OsTypeLinuxNeokylin
	} else if strings.Contains(name, "opencloudos") {
		return OsTypeLinuxNeokylin
	} else {
		return 0
	}
}

// CanUpgradeDirect 直接比较版本号是否能升级
// 参数: newest 新版本号
// 参数: original 原版本号
// 返回: 如果newest大于original为真
func CanUpgradeDirect(newest, original string) bool {
	if newest != original {
		return true
	}
	return false
}

// CanUpgradeIndirect 间接比较版本号是否能升级
// @return 是否能实际, 对应升级版本号
func CanUpgradeIndirect(newestMap map[int]string, agentOsType int, agentVer string) (bool, string) {
	if newestMap == nil || len(newestMap) == 0 {
		return false, ""
	}
	if newestVer, isOk := newestMap[agentOsType]; isOk {
		if CanUpgradeDirect(newestVer, agentVer) {
			return false, newestVer
		} else {
			return false, ""
		}
	} else {
		return false, ""
	}
}

// DownloadName 下载名称
func DownloadName(platformType, agentOsType int, agentVer string) string {
	name := "MA_"
	name += enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type", agentOsType, EnumMap)
	name += "("
	name += strings.ReplaceAll(agentVer, ".", "_")
	name += ").zip"
	return name
}

// ShowVer  显示版本号
func ShowVer(platformType, agentOsType int, agentVer string) string {
	name := "MA_"
	name += enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type_alias", platformType, EnumMap)
	name += "_"
	name += agentVer
	return name
}
