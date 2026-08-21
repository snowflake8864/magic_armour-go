package enum_software

import (
	"goskeleton/app/global/enums"
	"strings"
)

const (
	PlatformTypeLinux   = 1 //平台类型: linux
	PlatformTypeWindows = 2 //平台类型: windows

	ArchAmd64       = 1 //CPU架构: amd64(即x86_64)
	ArchArm64       = 2 //CPU架构: aarch64
	ArchMips64      = 3 //CPU架构: mips64
	ArchLoongarch64 = 4 //CPU架构: loongarch64(即国产龙芯)

	OsTypeLinuxUnknown      = 100 //系统类型: linux 未知
	OsTypeLinuxUbuntu       = 101 //系统类型: linux Ubuntu
	OsTypeLinuxCentos       = 102 //系统类型: linux Centos
	OsTypeLinuxKylin        = 103 //系统类型: linux 麒麟(NEO开头的是中标麒麟)
	OsTypeLinuxUos          = 104 //系统类型: linux 统信
	OsTypeLinuxRedhat       = 105 //系统类型: linux 红帽
	OsTypeLinuxDebian       = 106 //系统类型: linux Debian
	OsTypeLinuxRocky        = 107 //系统类型: linux Rocky
	OsTypeLinuxAnolis       = 108 //系统类型: linux 阿里龙蜥
	OsTypeLinuxOpenCloudOs  = 109 //系统类型: linux 腾讯云
	OsTypeLinuxAlibabaCloud = 110 //系统类型: linux 阿里云
	OsTypeLinuxAlmaLinux    = 111 //系统类型: linux 艾玛（阿里支持）
	OsTypeLinuxCentosStream = 112 //系统类型: linux Centos Stream
	OsTypeLinuxFedora       = 113 //系统类型: linux Fedora
	OsTypeLinuxEuler        = 114 //系统类型: linux Euler
	OsTypeLinuxSuse         = 115 //系统类型: linux SUSE

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
	"arch_type": {
		ArchAmd64:  "amd64",
		ArchArm64:  "arm64",
		ArchMips64: "mips64",
	},
	"os_type": {
		OsTypeWinSeven:          "Windows",
		OsTypeWinXp:             "Windows Early",
		OsTypeLinuxUbuntu:       "Ubuntu",
		OsTypeLinuxCentos:       "CentOS",
		OsTypeLinuxKylin:        "Kylin",
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
		OsTypeLinuxEuler:        "Euler",
		OsTypeLinuxSuse:         "SUSE",
	},
	"os_type_alias": {
		OsTypeWinSeven:          "windows7 and later versions",
		OsTypeWinXp:             "windows xp-windows 2003",
		OsTypeLinuxUbuntu:       "Ubuntu 64bit",
		OsTypeLinuxCentos:       "CentOS 64bit",
		OsTypeLinuxKylin:        "kylin 64bit",
		OsTypeLinuxUos:          "Uos 64bit",
		OsTypeLinuxRedhat:       "RedHat 64bit",
		OsTypeLinuxDebian:       "Debian 64bit",
		OsTypeLinuxRocky:        "Rocky 64bit",
		OsTypeLinuxAnolis:       "Anolis 64bit",
		OsTypeLinuxOpenCloudOs:  "OpenCloudOs 64bit",
		OsTypeLinuxAlibabaCloud: "AlibabaCloud 64bit",
		OsTypeLinuxAlmaLinux:    "AlmaLinux 64bit",
		OsTypeLinuxCentosStream: "CentosStream 64bit",
		OsTypeLinuxFedora:       "Fedora 64bit",
		OsTypeLinuxEuler:        "Euler 64bit",
		OsTypeLinuxSuse:         "SUSE 64bit",
	},
	"os_ver": {
		//Ubuntu
		1010: "64bit(14.04-22.04)",
		//CentOS
		1020: "64bit(6.5 - 8.5)",
		//麒麟
		1030: "64bit NeoKylin(V7)",
		//统信
		1040: "64bit (20)",
		//RedHat
		1050: "64bit(6.5 - 7.9)",
		//Debian
		1060: "64bit(14.04 - 22.04)",
		//Rocky
		1070: "64bit(8.4 - 9.1)",
		//ANOLIS
		1080: "64bit(7.7 - 8.8)",
		//OpenCloudos
		1090: "64bit(8.6-9.0)",
		//AlibabaCloud
		1100: "64bit",
		//AlmaLinux
		1110: "64bit",
		//CentosStream
		1120: "64bit(8-9)",
		//Fedora
		1130: "64bit",
		//Euler
		1140: "64bit",
		//SUSE
		1150: "64bit",

		//windows
		2010: "windows7 and later versions",
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
			OsTypeLinuxKylin,
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
			OsTypeLinuxEuler,
			OsTypeLinuxSuse,
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
		OsTypeLinuxKylin: {
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
		OsTypeLinuxEuler: {
			1140,
		},
		OsTypeLinuxSuse: {
			1150,
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
	} else if strings.Contains(name, "kylin") || strings.Contains(name, "麒麟") {
		return OsTypeLinuxKylin
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
	} else if strings.Contains(name, "euler") || strings.Contains(name, "fusionos") {
		return OsTypeLinuxEuler
	} else if strings.Contains(name, "suse") || strings.Contains(name, "sles") {
		return OsTypeLinuxSuse
	} else {
		return OsTypeLinuxUnknown
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
