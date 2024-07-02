package enum_safety_report

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
)

const (
	AlarmNoticeCycleDay   int = 1  //告警周期间隔：每天
	AlarmNoticeCycleWeek  int = 7  //告警周期间隔：每星期
	AlarmNoticeCycleMonth int = 30 //告警周期间隔：每月

)

var EnumIntMap = map[string]map[int]string{
	"notice_cycle": {
		AlarmNoticeCycleDay:   "每天",
		AlarmNoticeCycleWeek:  "每周",
		AlarmNoticeCycleMonth: "每月",
	},
}
var WebMapInt = map[string][]int{
	"notice_cycle": {
		AlarmNoticeCycleDay,
		AlarmNoticeCycleWeek,
		AlarmNoticeCycleMonth,
	},
	"alarm_cycle_mod": {
		enum_alarm.ModuleProcess,
		enum_alarm.ModuleDll,
		enum_alarm.ModuleTamper,
		enum_alarm.ModuleBlackmail,
		enum_alarm.ModulePort,
		enum_alarm.ModuleUsb,
		enum_alarm.ModulePlug,
		enum_alarm.ModuleOutreach,
		enum_alarm.ModuleLogin,
		enum_alarm.ModuleSelf,
	},
	"alarm_timely_mod": {
		enum_alarm.ModuleLogin,
	},
}

// AlarmNoticeConfig 告警通知配置
type AlarmNoticeConfig struct {
	EmailSwitch   int    `json:"email_switch"`
	Email         string `json:"email"`
	CycleInterval int    `json:"cycle_interval"`
	CycleMod      []int  `json:"cycle_mod"`
	TimelySwitch  int    `json:"timely_switch"`
	TimelyMod     []int  `json:"timely_mod"`
}

func (c *AlarmNoticeConfig) Format() {
	if c.EmailSwitch == 0 {
		c.EmailSwitch = enums.SwitchOff
	}
	if c.CycleInterval == 0 {
		c.CycleInterval = AlarmNoticeCycleDay
	}
	if c.TimelySwitch == 0 {
		c.TimelySwitch = enums.SwitchOff
	}
	if c.CycleMod == nil || len(c.CycleMod) == 0 {
		c.CycleMod = []int{}
	}
	if c.TimelyMod == nil || len(c.TimelyMod) == 0 {
		c.TimelyMod = []int{}
	}
}

// HardwareNoticeConfig 硬件资源通知配置
type HardwareNoticeConfig struct {
	EmailSwitch  int    `json:"email_switch"`
	Email        string `json:"email"`
	CpuSwitch    int    `json:"cpu_switch"`
	CpuThreshold int    `json:"cpu_threshold"`
	MemSwitch    int    `json:"mem_switch"`
	MemThreshold int    `json:"timely_mod"`
	HdSwitch     int    `json:"hd_switch"`
	HdThreshold  int    `json:"hd_threshold"`
	TimelySwitch int    `json:"timely_switch"`
}

func (c *HardwareNoticeConfig) Format() {
	if c.EmailSwitch == 0 {
		c.EmailSwitch = enums.SwitchOff
	}
	if c.CpuSwitch == 0 {
		c.CpuSwitch = enums.SwitchOff
	}
	if c.MemSwitch == 0 {
		c.MemSwitch = enums.SwitchOff
	}
	if c.HdSwitch == 0 {
		c.HdSwitch = enums.SwitchOff
	}
	if c.TimelySwitch == 0 {
		c.TimelySwitch = enums.SwitchOff
	}
}
