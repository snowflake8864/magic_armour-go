package enum_software

import (
	"goskeleton/app/global/enums"
	"strings"
)

const (
	PlatformTypeLinux   = 1 //平台类型: linux
	PlatformTypeWindows = 2 //平台类型: windows

	OsTypeLinuxUbuntu       = 101 //系统类型: linux Ubuntu
	OsTypeLinuxCentos       = 102 //系统类型: linux Centos
	OsTypeLinuxNeokylin     = 103 //系统类型: linux 麒麟
	OsTypeLinuxUos          = 104 //系统类型: linux 通信
	OsTypeLinuxRedhat       = 105 //系统类型: linux 红帽
	OsTypeLinuxDebian       = 106 //系统类型: linux Debian
	OsTypeLinuxRocky        = 107 //系统类型: linux Rocky
	OsTypeLinuxAnolis       = 108 //系统类型: linux 阿里龙蜥
	OsTypeLinuxOpenCloudOs  = 109 //系统类型: linux 腾讯云
	OsTypeLinuxAlibabaCloud = 110 //系统类型: linux 阿里云
	OsTypeLinuxAlmaLinux    = 111 //系统类型: linux AlmaLinux
	OsTypeLinuxCentosStream = 112 //系统类型: linux Centos Stream
	OsTypeLinuxFedora       = 113 //系统类型: linux Fedora

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
		OsTypeWinSeven:          "Windows",
		OsTypeWinXp:             "Windows早期",
		OsTypeLinuxUbuntu:       "Ubuntu",
		OsTypeLinuxCentos:       "CentOS",
		OsTypeLinuxNeokylin:     "Neokylin",
		OsTypeLinuxUos:          "Uos",
		OsTypeLinuxRedhat:       "RedHat",
		OsTypeLinuxDebian:       "Debian",
		OsTypeLinuxRocky:        "Rocky",
		OsTypeLinuxAnolis:       "Anolis",
		OsTypeLinuxOpenCloudOs:  "OpenCloudOs",
		OsTypeLinuxAlibabaCloud: "AlibabaCloud",
		OsTypeLinuxAlmaLinux:    "AlmaLinux",
		OsTypeLinuxCentosStream: "CentosStream",
		OsTypeLinuxFedora:       "Fedora",
	},
	"os_type_alias": {
		OsTypeWinSeven:          "windows7-至今",
		OsTypeWinXp:             "windows xp-windows 2003",
		OsTypeLinuxUbuntu:       "Ubuntu64位",
		OsTypeLinuxCentos:       "CentOS64位",
		OsTypeLinuxNeokylin:     "Neokylin64位",
		OsTypeLinuxUos:          "Uos64位",
		OsTypeLinuxRedhat:       "RedHat64位",
		OsTypeLinuxDebian:       "Debian64位",
		OsTypeLinuxRocky:        "Rocky64位",
		OsTypeLinuxAnolis:       "Anolis64位",
		OsTypeLinuxOpenCloudOs:  "OpenCloudOs64位",
		OsTypeLinuxAlibabaCloud: "AlibabaCloud64位",
		OsTypeLinuxAlmaLinux:    "AlmaLinux64位",
		OsTypeLinuxCentosStream: "CentosStream64位",
		OsTypeLinuxFedora:       "Fedora64位",
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
		//OpenCloudos
		1090: "64位(8.6-9.0)",
		//AlibabaCloud
		1100: "64位",
		//AlmaLinux
		1110: "64位",
		//CentosStream
		1120: "64位(8-9)",
		//Fedora
		1130: "64位",

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
			OsTypeLinuxOpenCloudOs,
			OsTypeLinuxAlibabaCloud,
			OsTypeLinuxAlmaLinux,
			OsTypeLinuxCentosStream,
			OsTypeLinuxFedora,
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
		OsTypeLinuxOpenCloudOs: {
			1090,
		},
		OsTypeLinuxAlibabaCloud: {
			1100,
		},
		OsTypeLinuxAlmaLinux: {
			1110,
		},
		OsTypeLinuxCentosStream: {
			1120,
		},
		OsTypeLinuxFedora: {
			1130,
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
			return OsTypeWinXp
		} else {
			return OsTypeWinSeven
		}
	} else if strings.Contains(name, "ubuntu") {
		return OsTypeLinuxUbuntu
	} else if strings.Contains(name, "centos") {
		if strings.Contains(name, "stream") {
			return OsTypeLinuxCentosStream
		} else {
			return OsTypeLinuxCentos
		}
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
		return OsTypeLinuxAnolis
	} else if strings.Contains(name, "opencloudos") {
		return OsTypeLinuxOpenCloudOs
	} else if strings.Contains(name, "alibaba cloud") || strings.Contains(name, "aliyun") {
		return OsTypeLinuxRedhat
	} else if strings.Contains(name, "almalinux") || strings.Contains(name, "alma") {
		return OsTypeLinuxAlmaLinux
	} else if strings.Contains(name, "fedora") {
		return OsTypeLinuxFedora
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
			return true, newestVer
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
