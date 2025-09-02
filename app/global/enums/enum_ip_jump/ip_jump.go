package enum_ip_jump

import (
	"fmt"
	"goskeleton/app/global/enums"
)

const (
	CycleDay   int = 1 //频率: 天
	CycleWeek  int = 2 //频率: 星期
	CycleMonth int = 3 //频率: 月

	WayRandom   int = 1 //方式: 随机
	WayOrderAsc int = 2 //方式: 顺序-从小到大

	NewestStatusSuccess int = 1
	NewestStatusFail    int = 2
)

var EnumIntMap = map[string]map[int]string{
	"cycle": {
		CycleDay:   "每天",
		CycleWeek:  "每周",
		CycleMonth: "每月",
	},
	"way": {
		WayOrderAsc: "顺序",
		WayRandom:   "随机",
	},
}

var WebMapInt = map[string][]int{
	"cycle": {
		CycleDay,
		CycleWeek,
		CycleMonth,
	},
	"way": {
		WayRandom,
		WayOrderAsc,
	},
}

// IpJumpConfig IP跳变
type IpJumpConfig struct {
	Switch  int    `json:"switch"`
	Ip      string `json:"ip"`      //
	Mask    int    `json:"mask"`    //
	Gateway string `json:"gateway"` //
	Cycle   int    `json:"cycle"`   //频率
	Timing  int    `json:"timing"`  //时机
	Way     int    `json:"way"`     //方式
}

func (c *IpJumpConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}
	if c.Mask == 0 {
		c.Mask = 24
	}
	if c.Cycle == 0 {
		c.Cycle = CycleWeek
	}
	if c.Way == 0 {
		c.Way = WayRandom
	}
}

// GetCycleTimingMap 获取周期和频率的对应关系
func GetCycleTimingMap(cycle int) []int {
	var list []int
	e := 1
	if cycle == CycleWeek {
		e = 7
	} else if cycle == CycleMonth {
		e = 30
	}
	for i := 1; i <= e; i++ {
		list = append(list, i)
	}
	return list
}

// GetTimingLabel 获取周频率标记
func GetTimingLabel(cycle, timing int) string {
	var data string
	if cycle == CycleWeek {
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
		data = fmt.Sprintf("周%v凌晨", val)
	} else if cycle == CycleMonth {
		data = fmt.Sprintf("%d号", timing)
	} else {
		data = "当天凌晨"
	}

	return data
}

// IpJumpNewest 跳变最新结果
type IpJumpNewest struct {
	SourceIp string `json:"source_ip"`
	TargetIp string `json:"target_ip"`
	Gateway  string `json:"gateway"`
	AgentIp  string `json:"agent_ip"`

	Status int    `json:"status"` //最近一次修改状态(0无修改 1成功  2失败)
	Reason string `json:"reason"` //最近一次修改失败原因
}

func (c *IpJumpNewest) Format() {

}
