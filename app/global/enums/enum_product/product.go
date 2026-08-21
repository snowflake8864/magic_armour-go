package enum_product

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/json_helper"
)

const (
	StatusNot int = 1 //状态：无效
	StatusYes int = 2 //状态：有效

	TypeGold   int = 1 //类型:正式
	TypeTest   int = 2 //类型:测试
	TypeSample int = 3 //类型:样本

	ModuleProcess   string = "process"   //模块:可信进程
	ModuleTamper    string = "tamper"    //模块:防篡改
	ModuleBlackmail string = "blackmail" //模块:勒索防护
	ModulePort      string = "port"      //模块:端口虚开
	ModulePlug      string = "plug"      //模块:封堵
	ModuleUsb       string = "usb"       //模块:外设
	ModuleSyslog    string = "syslog"    //模块:通信日志
	ModuleOutreach  string = "outreach"  //模块:外联
	ModuleLogin     string = "login"     //模块:登录口令破解
	ModuleThreaten  string = "threaten"  //模块:威胁情报
)

var EnumIntMap = map[string]map[int]string{
	"status": {
		StatusNot: "无效",
		StatusYes: "有效",
	},
	"status_en": {
		StatusNot: "Invalid",
		StatusYes: "Valid",
	},
	"type": {
		TypeGold:   "正式",
		TypeTest:   "测试",
		TypeSample: "样本",
	},
	"type_en": {
		TypeGold:   "Formal",
		TypeTest:   "Test",
		TypeSample: "Sample",
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
		ModuleThreaten:  "威胁情报",
	},
	"module_en": {
		ModuleProcess:   "Trusted Protection",
		ModuleTamper:    "Anti-Tampering",
		ModuleBlackmail: "Ransomware Protection",
		ModulePort:      "Port Deception",
		ModulePlug:      "Dynamic Blocking",
		ModuleUsb:       "Peripheral Management",
		ModuleSyslog:    "Communication Log",
		ModuleOutreach:  "Outreach Control",
		ModuleLogin:     "Password Cracking Protection",
		ModuleThreaten:  "Threat Intelligence",
	},
}

var WebMapInt = map[string][]int{
	"type": {
		TypeGold,
		TypeTest,
	},
}
var WebMapStr = map[string][]string{
	"module": {
		ModuleProcess,
		ModuleTamper,
		ModuleBlackmail,
		ModulePort,
		ModulePlug,
		ModuleUsb,
		ModuleSyslog,
		ModuleOutreach,
		ModuleLogin,
		//ModuleThreaten,
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
	Threaten  int `json:"threaten"`
}

// Default 默认配置
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
	if c.Threaten == 0 {
		c.Threaten = enums.IsYes
	} else {
		c.Threaten = enums.IsYes
	}
}

// SelectedTxt 选中模块文案
func (c *ProductGrantConfig) SelectedTxt() (txtList []string) {
	if c.Process != enums.IsNot {
		txtList = append(txtList, "可信防护")
	}
	if c.Tamper != enums.IsNot {
		txtList = append(txtList, "防篡改")
	}
	if c.Blackmail != enums.IsNot {
		txtList = append(txtList, "勒索防护")
	}
	if c.Port != enums.IsNot {
		txtList = append(txtList, "端口虚开")
	}
	if c.Plug != enums.IsNot {
		txtList = append(txtList, "动态封堵")
	}
	if c.Usb != enums.IsNot {
		txtList = append(txtList, "外设管理")
	}
	if c.Syslog != enums.IsNot {
		txtList = append(txtList, "通信日志")
	}
	if c.Outreach != enums.IsNot {
		txtList = append(txtList, "外联管控")
	}
	if c.Login != enums.IsNot {
		txtList = append(txtList, "口令破解")
	}
	if c.Threaten != enums.IsNot {
		//txtList = append(txtList, "威胁情报")
	}
	return txtList
}

// SelectedTxtEn 选中模块文案
func (c *ProductGrantConfig) SelectedTxtEn() (txtList []string) {
	if c.Process != enums.IsNot {
		txtList = append(txtList, "Trusted Protection")
	}
	if c.Tamper != enums.IsNot {
		txtList = append(txtList, "Anti-Tampering")
	}
	if c.Blackmail != enums.IsNot {
		txtList = append(txtList, "Blackmail Protection")
	}
	if c.Port != enums.IsNot {
		txtList = append(txtList, "Port Deception")
	}
	if c.Plug != enums.IsNot {
		txtList = append(txtList, "Dynamic Blocking")
	}
	if c.Usb != enums.IsNot {
		txtList = append(txtList, "Peripheral Management")
	}
	if c.Syslog != enums.IsNot {
		txtList = append(txtList, "Communication Log")
	}
	if c.Outreach != enums.IsNot {
		txtList = append(txtList, "Outreach Control")
	}
	if c.Login != enums.IsNot {
		txtList = append(txtList, "Password Cracking Protection")
	}
	if c.Threaten != enums.IsNot {
		//txtList = append(txtList, "Threat Intelligence")
	}
	return txtList
}

// Num2Config 数字转成标准配置
func (c *ProductGrantConfig) Num2Config() (r ProductGrantConfig) {
	if c.Process == 0 {
		r.Process = enums.IsNot
	} else {
		r.Process = enums.IsYes
	}
	if c.Tamper == 0 {
		r.Tamper = enums.IsNot
	} else {
		r.Tamper = enums.IsYes
	}
	if c.Blackmail == 0 {
		r.Blackmail = enums.IsNot
	} else {
		r.Blackmail = enums.IsYes
	}
	if c.Port == 0 {
		r.Port = enums.IsNot
	} else {
		r.Port = enums.IsYes
	}
	if c.Plug == 0 {
		r.Plug = enums.IsNot
	} else {
		r.Plug = enums.IsYes
	}
	if c.Usb == 0 {
		r.Usb = enums.IsNot
	} else {
		r.Usb = enums.IsYes
	}
	if c.Syslog == 0 {
		r.Syslog = enums.IsNot
	} else {
		r.Syslog = enums.IsYes
	}
	if c.Outreach == 0 {
		r.Outreach = enums.IsNot
	} else {
		r.Outreach = enums.IsYes
	}
	if c.Login == 0 {
		r.Login = enums.IsNot
	} else {
		r.Login = enums.IsYes
	}
	if c.Threaten == 0 {
		r.Threaten = enums.IsYes
	} else {
		r.Threaten = enums.IsYes
	}
	return r
}

// MinGreaterThanZero 如果是数字， 返回最小且大于0的数
func (c *ProductGrantConfig) MinGreaterThanZero() int {
	numbers := []int{
		c.Process,
		c.Tamper,
		c.Blackmail,
		c.Port,
		c.Plug,
		c.Usb,
		c.Syslog,
		c.Outreach,
		c.Login,
		c.Threaten,
	}

	minimum := 0
	for _, num := range numbers {
		if num > 0 && (minimum == 0 || num < minimum) {
			minimum = num
		}
	}

	return minimum
}

// Sub 如果是数字，大于0时减去值
func (c *ProductGrantConfig) Sub(x int) {
	if c.Process > 0 {
		c.Process -= x
		if c.Process < 0 {
			c.Process = 0
		}
	}
	if c.Tamper > 0 {
		c.Tamper -= x
		if c.Tamper < 0 {
			c.Tamper = 0
		}
	}
	if c.Blackmail > 0 {
		c.Blackmail -= x
		if c.Blackmail < 0 {
			c.Blackmail = 0
		}
	}
	if c.Port > 0 {
		c.Port -= x
		if c.Port < 0 {
			c.Port = 0
		}
	}
	if c.Plug > 0 {
		c.Plug -= x
		if c.Plug < 0 {
			c.Plug = 0
		}
	}
	if c.Usb > 0 {
		c.Usb -= x
		if c.Usb < 0 {
			c.Usb = 0
		}
	}
	if c.Syslog > 0 {
		c.Syslog -= x
		if c.Syslog < 0 {
			c.Syslog = 0
		}
	}
	if c.Outreach > 0 {
		c.Outreach -= x
		if c.Outreach < 0 {
			c.Outreach = 0
		}
	}
	if c.Login > 0 {
		c.Login -= x
		if c.Login < 0 {
			c.Login = 0
		}
	}
	if c.Threaten > 0 {
		c.Threaten -= x
		if c.Threaten < 0 {
			c.Threaten = 0
		}
	}
}

// Json2Config json转配置
func Json2Config(jsonStr string) (data ProductGrantConfig) {
	err := json_helper.ToStruct(jsonStr, &data)
	if err != nil {
		data = ProductGrantConfig{}
	}
	data.Default()
	return data
}
