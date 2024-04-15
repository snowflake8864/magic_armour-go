package enum_safety_report

import "goskeleton/app/global/enums"

const (
	AlarmNoticeCycleDay   int = 1  //告警周期间隔：每天
	AlarmNoticeCycleWeek  int = 7  //告警周期间隔：每星期
	AlarmNoticeCycleMonth int = 30 //告警周期间隔：每月

)

// AlarmNoticeConfig 告警通知配置
type AlarmNoticeConfig struct {
	EmailSwitch   int    `json:"email_switch"`
	CycleInterval int    `json:"cycle_interval"`
	CycleMod      string `json:"cycle_mod"`
	TimelySwitch  int    `json:"timely_switch"`
	TimelyMod     string `json:"timely_mod"`
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
}

// HardwareNoticeConfig 硬件资源通知配置
type HardwareNoticeConfig struct {
	EmailSwitch  int `json:"email_switch"`
	CpuSwitch    int `json:"cpu_switch"`
	CpuThreshold int `json:"cpu_threshold"`
	MemSwitch    int `json:"mem_switch"`
	MemThreshold int `json:"timely_mod"`
	HdSwitch     int `json:"hd_switch"`
	HdThreshold  int `json:"hd_threshold"`
	TimelySwitch int `json:"timely_switch"`
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
