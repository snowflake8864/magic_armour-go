package enum_product

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
		TypeGold: "正式",
		TypeTest: "测试",
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

// ProductGrantConfig 终端授权配置
type ProductGrantConfig struct {
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

func (c *ProductGrantConfig) Default() {
	if c.Process == 0 {
		c.Process = enums.IsNot
	} else {
		c.Process = enums.IsYes
	}
	if c.Tamper == 0 {
		c.Tamper = enums.IsNot
	} else {
		c.Tamper = enums.IsYes
	}
	if c.Blackmail == 0 {
		c.Blackmail = enums.IsNot
	} else {
		c.Blackmail = enums.IsYes
	}
	if c.Port == 0 {
		c.Port = enums.IsNot
	} else {
		c.Port = enums.IsYes
	}
	if c.Plug == 0 {
		c.Plug = enums.IsNot
	} else {
		c.Plug = enums.IsYes
	}
	if c.Usb == 0 {
		c.Usb = enums.IsNot
	} else {
		c.Usb = enums.IsYes
	}
	if c.Syslog == 0 {
		c.Syslog = enums.IsNot
	} else {
		c.Syslog = enums.IsYes
	}
	if c.Outreach == 0 {
		c.Outreach = enums.IsNot
	} else {
		c.Outreach = enums.IsYes
	}
	if c.Login == 0 {
		c.Login = enums.IsNot
	} else {
		c.Login = enums.IsYes
	}
}

func (c *ProductGrantConfig) SelectedTxt() (txtList []string) {
	if c.Process == enums.IsYes {
		txtList = append(txtList, "可信防护")
	}
	if c.Tamper == enums.IsYes {
		txtList = append(txtList, "防篡改")
	}
	if c.Blackmail == enums.IsYes {
		txtList = append(txtList, "勒索防护")
	}
	if c.Port == enums.IsYes {
		txtList = append(txtList, "端口虚开")
	}
	if c.Plug == enums.IsYes {
		txtList = append(txtList, "动态封堵")
	}
	if c.Usb == enums.IsYes {
		txtList = append(txtList, "外设管理")
	}
	if c.Syslog == enums.IsYes {
		txtList = append(txtList, "通信日志")
	}
	if c.Outreach == enums.IsYes {
		txtList = append(txtList, "外联管控")
	}
	if c.Login == enums.IsYes {
		txtList = append(txtList, "口令破解")
	}
	return txtList
}
