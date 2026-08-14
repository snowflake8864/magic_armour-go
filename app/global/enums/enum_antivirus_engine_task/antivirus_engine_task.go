package enum_antivirus_engine_task

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/date_helper"
	"time"
)

const (
	EngineTypeClamav int = 1 //引擎类型: clamav

	TargetTypeSingleAgent int = 1 //目标类型: 单个终端
	TargetTypeAllAgent    int = 2 //目标类型: 全部终端
	TargetTypeDepartments int = 3 //目标类型: 组织架构集合

	UpgradeModeManual int = 1 //升级模式: 人工
	UpgradeModeAuto   int = 2 //升级模式: 自动

	TaskTypeManualInstantly int = 10 //任务类型: 立即
	TaskTypeManualFixed     int = 11 //任务类型: 指定
	TaskTypeCycleWeekly     int = 20 //任务类型: 每周
	TaskTypeCycleMonth      int = 21 //任务类型: 每月

	ExecStatusWaiting  int = 1 //执行方式: 任务等待中（默认）
	ExecStatusStarting int = 2 //执行方式: 任务开始
	ExecStatusHandling int = 3 //执行方式: 任务进行中

	ProgressStatusUnfinished int = 1 //进度状态： 未完成
	ProgressStatusFinished   int = 2 //进度状态： 已完成
	ProgressStatusFail       int = 3 //进度状态： 失败

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
		TaskTypeManualFixed:     "定时执行",
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
		TaskTypeManualFixed,
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

// GetExecNextTs 计算下一次执行时间戳
func GetExecNextTs(taskType int, execDate, execHis string, execWeekly, execDaily int) (execNextTs int64) {
	if taskType == TaskTypeManualInstantly {
		return execNextTs
	}
	//标准的时分秒
	var standardHis string
	if c := len(execHis); c == 5 {
		standardHis = fmt.Sprintf("%s:00", execHis)
	} else if c == 8 {
		standardHis = execHis
	}

	if taskType == TaskTypeManualFixed {
		execNextTs = date_helper.DatetimeToTimestamp(fmt.Sprintf("%s %s", execDate, standardHis), variable.DateFormatNow)
	} else if taskType == TaskTypeCycleWeekly {
		//映射自定义当前周几
		var customizeCurrWeek int
		if currWeekday := date_helper.WeekdayCurr(); currWeekday == time.Sunday {
			customizeCurrWeek = 7
		} else {
			customizeCurrWeek = int(currWeekday)
		}

		//相同是当天
		if customizeCurrWeek == execWeekly {
			//现在时间戳
			nowTs := date_helper.CurTimestamp()
			//今天日期
			todayAt := date_helper.SpecialDatetime(variable.DateFormatToday) + " " + standardHis
			//今天时间戳
			todayTs := date_helper.DatetimeToTimestamp(todayAt, variable.DateFormatNow)
			//还没到今天的执行时间
			if todayTs > nowTs {
				execNextTs = todayTs
			} else {
				execNextTs = todayTs + (7 * 86400)
			}
		} else {
			var diff int
			//本周不到
			if execWeekly > customizeCurrWeek {
				diff = execWeekly - customizeCurrWeek
			} else {
				diff = 7 - (customizeCurrWeek - execWeekly)
			}
			//今天日期
			todayAt := date_helper.SpecialDatetime(variable.DateFormatToday) + " " + standardHis
			//今天时间戳
			todayTs := date_helper.DatetimeToTimestamp(todayAt, variable.DateFormatNow)

			execNextTs = todayTs + (int64(diff) * 86400)
		}
	} else if taskType == TaskTypeCycleMonth {
		//当前月的第几天
		currDaily := date_helper.DayInCurrMonth()

		//相同是当天
		if currDaily == execDaily {
			//现在时间戳
			nowTs := date_helper.CurTimestamp()
			//今天日期
			todayAt := date_helper.SpecialDatetime(variable.DateFormatToday) + " " + standardHis
			//今天时间戳
			todayTs := date_helper.DatetimeToTimestamp(todayAt, variable.DateFormatNow)
			//还没到今天的执行时间
			if todayTs > nowTs {
				execNextTs = todayTs
			} else {
				execNextTs = date_helper.AddTimestampByTimestamp(todayTs, 0, 1, 0)
			}
		} else {
			//今天日期
			todayAt := date_helper.SpecialDatetime(variable.DateFormatToday) + " " + standardHis
			//今天时间戳
			todayTs := date_helper.DatetimeToTimestamp(todayAt, variable.DateFormatNow)

			var diff int
			//本月不到
			if execDaily > currDaily {
				execNextTs = date_helper.AddTimestampByTimestamp(todayTs, 0, 0, execDaily-currDaily)
			} else {
				diff = date_helper.LastDayInMonthByTimestamp(todayTs) - currDaily + execDaily
			}
			execNextTs = date_helper.AddTimestampByTimestamp(todayTs, 0, 0, diff)
		}
	}

	return execNextTs
}
