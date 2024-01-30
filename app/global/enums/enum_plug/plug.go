package enum_plug

import "goskeleton/app/global/enums"

// 数据字典

const (
	IpTypeWhite int = 1 //IP策略类型:白名单
	IpTypeBlack int = 2 //IP策略类型:黑名单

	PlugTypeTemporary int = 1 //封堵类型:临时
	PlugTypePermanent int = 2 //封堵类型:永久

	SourcePortVirtual int = 1 //来源:端口虚开
	SourceLoginSolve  int = 3 //来源:口令破解

	DirectionIn  int = 1 //方向:IN，二进制表示
	DirectionOut int = 2 //方向:OUT，二进制表示
)

var EnumMap = map[string]map[int]string{
	"ip_type": {
		IpTypeWhite: "白名单",
		IpTypeBlack: "黑名单",
	},
	"plug_type": {
		PlugTypeTemporary: "临时",
		PlugTypePermanent: "永久",
	},
	"source": {
		SourcePortVirtual: "端口虚开",
		SourceLoginSolve:  "口令暴力破解",
	},
}

// PlugPortConfig 封堵：端口虚开配置
type PlugPortConfig struct {
	Switch    int `json:"switch"`
	Weight    int `json:"weight"`
	Type      int `json:"type"`
	CountTime int `json:"count_time"`
	PlugTime  int `json:"plug_time"`
}

func (c *PlugPortConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}
	if c.Weight == 0 {
		c.Weight = 10
	}
	if c.Type == 0 {
		c.Type = PlugTypeTemporary
	}
	if c.CountTime == 0 {
		c.CountTime = 2
	}
	if c.PlugTime == 0 {
		c.PlugTime = 30
	}
	if c.Type == PlugTypePermanent {
		c.PlugTime = 0
	}
}

// PlugLoginConfig 封堵：口令破解配置
type PlugLoginConfig struct {
	Switch   int `json:"switch"`
	Type     int `json:"type"`
	PlugTime int `json:"plug_time"`
}

func (c *PlugLoginConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}
	if c.Type == 0 {
		c.Type = PlugTypeTemporary
	}
	if c.PlugTime == 0 {
		c.PlugTime = 30
	}
	if c.Type == PlugTypePermanent {
		c.PlugTime = 0
	}
}
