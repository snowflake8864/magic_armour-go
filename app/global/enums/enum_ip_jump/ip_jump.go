package enum_ip_jump

import (
	"goskeleton/app/global/enums"
)

const (
	CycleDay   int = 1 //周期: 天
	CycleWeek  int = 2 //周期: 星期
	CycleMonth int = 3 //周期: 月
)

var EnumIntMap = map[string]map[int]string{
	"cycle": {
		CycleDay:   "每天",
		CycleWeek:  "每周",
		CycleMonth: "每月",
	},
}

var WebMapInt = map[string][]int{
	"cycle": {
		CycleDay,
		CycleWeek,
		CycleMonth,
	},
}

// PwJumpConfig 密码跳变
type PwJumpConfig struct {
	Switch  int    `json:"switch"`
	Cycle   int    `json:"cycle"`   //
	Ip      string `json:"ip"`      //
	Mask    int    `json:"mask"`    //
	Gateway string `json:"gateway"` //
}

func (c *PwJumpConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}

	if c.Cycle == 0 {
		c.Cycle = CycleWeek
	}
	if c.Mask == 0 {
		c.Mask = 24
	}
}
