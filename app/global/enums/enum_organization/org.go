package enum_organization

const (
	TypeOrg    int = 1 //类型:企业
	TypePerson int = 2 //类型:个人

	StatusVerifying int = 1 //状态:审核中
	StatusNormal    int = 2 //状态:正常
	StatusLocking   int = 3 //状态:锁定中

	OrgFieldDefault string = "id" //默认排序字段

	LogActionPlatformAdd            int = 100 //日志行为: 平台添加
	LogActionPlatformApprove        int = 101 //日志行为: 通过审核
	LogActionPlatformRefuse         int = 102 //日志行为: 未通过审核
	LogActionPlatformLock           int = 103 //日志行为: 锁定
	LogActionPlatformUnlock         int = 104 //日志行为: 锁定
	LogActionPlatformClean          int = 105 //日志行为: 清理业务数据
	LogActionPlatformHuanshiSwitch  int = 106 //日志行为: 设置幻势开关
	LogActionPlatformLoginSmsSwitch int = 107 //日志行为: 设置企业登录短信验证开关

	LogActionOrgRegister   int = 200 //日志行为: 用户注册
	LogActionOrgAddSubUser int = 201 //日志行为: 添加子用户
	LogActionOrgDelSubUser int = 202 //日志行为: 添加子用户
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
	},
}
