package enum_agent

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/json_helper"
)

// 终端

const (
	ModeMonitor int = 1 //模式：监控
	ModeProject int = 2 //模式：保护

	UpgradeStatusDefault     int = 1 //升级状态：默认
	upgradeStatusIssue       int = 2 //升级状态：任务进行中
	upgradeStatusDownloading int = 3 //升级状态：下载中
	upgradeStatusSuccess     int = 4 //升级状态:成功
	upgradeStatusFail        int = 5 //升级状态：失败

	HeartbeatCommonTime int = 60 //心跳:普通时间（秒）
	HeartbeatFastTime   int = 15 //心跳:快速时间（秒）

	HardwareTime int = 300 //上传硬件资源时间间隔（秒）

	LogProtoUdp int = 1 //日志协议:TDP
	LogProtoTcp int = 2 //日志协议:TCP

	AgentFieldDefault string = "id" //默认排序字段
)

var EnumIntMap = map[string]map[int]string{
	"is_online": {
		enums.IsYes: "在线",
		enums.IsNot: "离线",
	},
	"is_safe": {
		enums.IsYes: "正常",
		enums.IsNot: "异常",
	},
	"mode": {
		ModeMonitor: "监控",
		ModeProject: "保护",
	},
}

var WebMapInt = map[string][]int{
	"is_online": {
		enums.IsYes,
		enums.IsNot,
	},
	"is_safe": {
		enums.IsYes,
		enums.IsNot,
	},
}

// AgentTypeDefault 终端类型： 默认
func AgentTypeDefault() map[int]string {
	return map[int]string{
		1: "服务器",
		2: "业务系统",
		3: "办公台式机",
		4: "办公笔记本",
		5: "其它",
	}
}

// GetTypeMap 终端类型： 获取
func GetTypeMap(jsonStr string) map[int]string {
	types := make(map[int]string)

	if err := json_helper.ToAny(jsonStr, &types); err != nil {
		return nil
	} else {
		return types
	}
}

// GetTypeLabel 终端类型： 获取标签
func GetTypeLabel(types map[int]string, val int) string {
	if types == nil {
		return ""
	} else {
		if label, isOK := types[val]; isOK {
			return label
		} else {
			return ""
		}
	}
}
