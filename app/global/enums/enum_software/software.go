package enum_software

import "strings"

const (
	PlatformTypeLinux   = 1 //平台类型: linux
	PlatformTypeWindows = 2 //平台类型: windows

	OsTypeWinSeven         = 101 //系统类型: win7以上
	OsTypeWinXp            = 102 //系统类型: win xp 和 win 2003
	OsTypeLinuxUbuntu      = 201 //系统类型: linux Ubuntu
	OsTypeLinuxCentos      = 202 //系统类型: linux Centos
	OsTypeLinuxNeokylin    = 203 //系统类型: linux 麒麟
	OsTypeLinuxUos         = 204 //系统类型: linux 通信
	OsTypeLinuxRedhat      = 205 //系统类型: linux 红帽
	OsTypeLinuxDebian      = 206 //系统类型: linux Debian
	OsTypeLinuxRocky       = 207 //系统类型: linux Rocky
	OsTypeLinuxAnolis      = 108 //系统类型: linux 阿里龙蜥
	OsTypeLinuxOpencloudos = 209 //系统类型: linux 腾讯
)

var EnumMap = map[string]map[int]string{
	"platform_type": {
		PlatformTypeWindows: "Windows",
		PlatformTypeLinux:   "Linux",
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
		//windows
		1010: "windows7-至今",
		//Windows其它
		1020: "windows xp-windows 2003",
		//Ubuntu
		2010: "64位(14.04-22.04)",
		//CentOS
		2020: "64位(6.5 - 7.9)",
		//麒麟
		2030: "64位 中标麒麟(V7)",
		//统信
		2040: "64位 (20)",
		//RedHat
		2050: "64位(6.5 - 7.9)",
		//Debian
		2060: "64位(14.04 - 22.04)",
		//Rocky
		2070: "64位(8.4 - 9.1)",
		//ANOLIS
		2080: "64位(7.7 - 8.8)",
		//OPENCLOUDOS
		2090: "64位(8.6-9.0)",
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
			1100,
		},
		OsTypeWinXp: {
			1200,
		},
		OsTypeLinuxUbuntu: {
			2010,
		},
		OsTypeLinuxCentos: {
			2020,
		},
		OsTypeLinuxNeokylin: {
			2030,
		},
		OsTypeLinuxUos: {
			2040,
		},
		OsTypeLinuxRedhat: {
			2050,
		},
		OsTypeLinuxDebian: {
			2060,
		},
		OsTypeLinuxRocky: {
			2070,
		},
		OsTypeLinuxAnolis: {
			2080,
		},
		OsTypeLinuxOpencloudos: {
			2090,
		},
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
func CanUpgradeIndirect(newestMap map[int]string, agentOsType int, agentVer string) bool {
	if newestMap == nil || len(newestMap) == 0 {
		return false
	}
	if newest, isOk := newestMap[agentOsType]; isOk {
		return CanUpgradeDirect(newest, agentVer)
	} else {
		return false
	}
}
