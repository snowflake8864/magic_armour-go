package enum_antivirus_engine_task

import (
	"fmt"
)

const (
	EngineTypeClamav int = 1 //引擎类型: clamav

	TargetTypeSingleAgent int = 1 //目标类型: 单个终端
	TargetTypeAllAgent    int = 2 //目标类型: 全部终端
	TargetTypeDepartments int = 3 //目标类型: 组织架构集合

	UpgradeModeManual int = 1 //升级模式: 人工
	UpgradeModeAuto   int = 2 //升级模式: 自动

	TaskTypeManualInstantly int = 10 //任务类型: 立即
	TaskTypeManualCrontab   int = 11 //任务类型: 定时
	TaskTypeCycleWeekly     int = 20 //任务类型: 每周
	TaskTypeCycleMonth      int = 21 //任务类型: 每月

	ExecStatusWaiting  int = 1 //执行方式: 任务等待中（默认）
	ExecStatusStarting int = 2 //执行方式: 任务开始
	ExecStatusHandling int = 3 //执行方式: 任务进行中

	ProgressStatusUnfinished int = 0 //进度状态： 未完成
	ProgressStatusFinished   int = 1 //进度状态： 已完成
	ProgressStatusFail       int = 2 //进度状态： 失败

	FailReasonOffline  int = 1 //错误原因：终端离线
	FailReasonUnneeded int = 2 //错误原因: 无需升级
	FailReasonEngine   int = 3 //错误原因: 引擎问题（agent成功接到任务开始计算，包括下载失败、升级失败等）
)

var EnumMap = map[string]map[int]string{
	"engine_type": {
		EngineTypeClamav: "clamav",
	}, "upgrade_mode": {
		UpgradeModeAuto:   "自动升级",
		UpgradeModeManual: "人工升级",
	},
	"task_type": {
		TaskTypeCycleWeekly:     "每周执行",
		TaskTypeCycleMonth:      "每月执行",
		TaskTypeManualInstantly: "立即执行",
		TaskTypeManualCrontab:   "定时执行",
	},

	"exec_status": {
		ExecStatusWaiting:  "等待开始",
		ExecStatusStarting: "开始执行",
		ExecStatusHandling: "执行中",
	},
}

var WebMapInt = map[string][]int{
	"upgrade_mode": {
		UpgradeModeManual,
		UpgradeModeAuto,
	},
	"task_type": {
		TaskTypeManualInstantly,
		TaskTypeManualCrontab,
		TaskTypeCycleWeekly,
		TaskTypeCycleMonth,
	},
}

// GetCycleTimingMap 获取周期和频率的对应关系
func GetCycleTimingMap(cycle int) []int {
	var list []int

	e := 1
	if cycle == TaskTypeCycleWeekly {
		e = 7
	} else if cycle == TaskTypeCycleMonth {
		e = 25
	}
	for i := 1; i <= e; i++ {
		list = append(list, i)
	}

	return list
}

// GetTimingLabel 获取周期频率标记
func GetTimingLabel(cycle, timing int) string {
	var data string
	if cycle == TaskTypeCycleWeekly {
		var val string
		switch timing {
		case 1:
			val = "一"
		case 2:
			val = "二"
		case 3:
			val = "三"
		case 4:
			val = "四"
		case 5:
			val = "五"
		case 6:
			val = "六"
		case 7:
			val = "日"
		}
		data = fmt.Sprintf("周%v", val)
	} else if cycle == TaskTypeCycleMonth {
		data = fmt.Sprintf("%d号", timing)
	}

	return data
}
