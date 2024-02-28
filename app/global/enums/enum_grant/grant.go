package enum_grant

import "goskeleton/app/global/enums"

const (
	StatusNot int = 1 //状态：无效
	StatusYes int = 2 //状态：有效

	TypeGold   int = 0 //类型:正式
	TypeTest   int = 1 //类型:测试
	TypeSample int = 2 //类型:样本

	ModuleProcess   string = "process"   //模块:可信进程
	ModuleTamper    string = "tamper"    //模块:防篡改
	ModuleBlackmail string = "blackmail" //模块:勒索防护
	ModulePort      string = "port"      //模块:端口虚开
	ModulePlug      string = "plug"      //模块:封堵
	ModuleUsb       string = "usb"       //模块:外设
	ModuleSyslog    string = "syslog"    //模块:通信日志
	ModuleOutreach  string = "outreach"  //模块:外联
	ModuleLogin     string = "login"     //模块:登录口令破解
)

var EnumIntMap = map[string]map[int]string{
	"status": {
		StatusNot: "无效",
		StatusYes: "有效",
	},
	"type": {
		TypeGold:   "正式",
		TypeTest:   "测试",
		TypeSample: "样本",
	},
}

var EnumStrMap = map[string]map[string]string{
	"module": {
		ModuleProcess:   "可信防护",
		ModuleTamper:    "防篡改",
		ModuleBlackmail: "勒索防护",
		ModulePort:      "端口虚开",
		ModulePlug:      "动态封堵",
		ModuleUsb:       "外设管理",
		ModuleSyslog:    "通信日志",
		ModuleOutreach:  "外联管控",
		ModuleLogin:     "口令破解",
	},
}

var EnumStrSlice = map[string][]string{
	//@todo 顺序不能改变（导入授权码的顺序）
	"module": {
		ModuleProcess,
		ModuleBlackmail,
		ModuleTamper,
		ModulePort,
		ModulePlug,
		ModuleUsb,
	},
}

// AgentGrantConfig 终端授权配置
type AgentGrantConfig struct {
	Process   int `json:"process"`
	Tamper    int `json:"tamper"`
	Blackmail int `json:"blackmail"`
	Port      int `json:"port"`
	Plug      int `json:"plug"`
	Usb       int `json:"usb"`
	Syslog    int `json:"syslog"`
	Outreach  int `json:"outreach"`
	Login     int `json:"login"`
}

func (c *AgentGrantConfig) Default() {
	c.Process = enums.IsNot
	c.Tamper = enums.IsNot
	c.Blackmail = enums.IsNot
	c.Port = enums.IsNot
	c.Plug = enums.IsNot
	c.Usb = enums.IsNot
	c.Syslog = enums.IsYes
	c.Outreach = enums.IsYes
	c.Login = enums.IsYes
}
