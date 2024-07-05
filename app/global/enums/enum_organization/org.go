package enum_organization

const (
	TypeOrg    int = 1 //类型:企业
	TypePerson int = 2 //类型:个人

	StatusVerifying int = 1 //状态:审核中
	StatusNormal    int = 2 //状态:正常
	StatusLocking   int = 3 //状态:锁定中

	LogActionPlatformAdd            int = 100 //日志行为: 平台添加
	LogActionPlatformApprove        int = 101 //日志行为: 通过审核
	LogActionPlatformRefuse         int = 102 //日志行为: 未通过审核
	LogActionPlatformLock           int = 103 //日志行为: 锁定
	LogActionPlatformUnlock         int = 104 //日志行为: 锁定
	LogActionPlatformClean          int = 105 //日志行为: 清理业务数据
	LogActionPlatformHuanshiSwitch  int = 106 //日志行为: 设置幻势开关
	LogActionPlatformLoginSmsSwitch int = 107 //日志行为: 设置企业登录短信验证开关

	LogActionOrgRegister        int = 200 //日志行为: 用户注册
	LogActionOrgAddSubUser      int = 201 //日志行为: 添加子用户
	LogActionOrgDelSubUser      int = 202 //日志行为: 删除子用户
	LogActionOrgEditSubUser     int = 203 //日志行为: 编辑子用户
	LogActionOrgModifyPwSubUser int = 204 //日志行为: 子用户修改密码
	LogActionOrgSelfModifyPw    int = 210 //日志行为: 修改密码
	LogActionOrgSelfModifyEmail int = 211 //日志行为: 修改邮箱
	LogActionOrgSelfModifyPhone int = 212 //日志行为: 修改手机
	LogActionOrgSelfModifyInfo  int = 213 //日志行为: 修改基础信息
)

var EnumMap = map[string]map[int]string{
	"type": {
		TypeOrg:    "企业",
		TypePerson: "个人",
	},
	"status": {
		StatusVerifying: "审核中",
		StatusNormal:    "企业",
		StatusLocking:   "锁定",
	},
	"log_action": {
		LogActionPlatformAdd:            "平台添加",
		LogActionPlatformApprove:        "通过审核",
		LogActionPlatformRefuse:         "未通过审核",
		LogActionPlatformLock:           "企业账户锁定",
		LogActionPlatformUnlock:         "企业账户解锁",
		LogActionPlatformClean:          "清理业务数据",
		LogActionPlatformHuanshiSwitch:  "设置幻势开关",
		LogActionPlatformLoginSmsSwitch: "设置企业登录短信验证开关",
		LogActionOrgRegister:            "用户注册",
		LogActionOrgAddSubUser:          "添加子账户",
		LogActionOrgDelSubUser:          "删除子账户",
		LogActionOrgEditSubUser:         "编辑子账户",
		LogActionOrgModifyPwSubUser:     "修改子用户密码",
		LogActionOrgSelfModifyPw:        "修改密码",
		LogActionOrgSelfModifyEmail:     "修改邮箱",
		LogActionOrgSelfModifyPhone:     "修改手机",
		LogActionOrgSelfModifyInfo:      "修改基本信息",
	},
}

var WebMapInt = map[string][]int{
	"type": {
		TypeOrg,
		TypePerson,
	},
	"log_action": {
		LogActionPlatformAdd,
		LogActionPlatformApprove,
		LogActionPlatformRefuse,
		LogActionPlatformLock,
		LogActionPlatformUnlock,
		LogActionPlatformClean,
		LogActionPlatformHuanshiSwitch,
		LogActionPlatformLoginSmsSwitch,
		LogActionOrgRegister,
		LogActionOrgAddSubUser,
		LogActionOrgDelSubUser,
		LogActionOrgEditSubUser,
		LogActionOrgModifyPwSubUser,
		LogActionOrgSelfModifyPw,
		LogActionOrgSelfModifyEmail,
		LogActionOrgSelfModifyPhone,
		LogActionOrgSelfModifyInfo,
	},
}
