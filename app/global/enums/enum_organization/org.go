package enum_organization

const (
	TypeOrg    int = 1 //类型:企业
	TypePerson int = 2 //类型:个人

	SourceReg      int = 1 //来源:自己注册
	SourcePlatform int = 2 //来源:平台添加

	StatusVerifying int = 1 //状态:审核中
	StatusNormal    int = 2 //状态:正常
	StatusLocking   int = 3 //状态:锁定中

	OrgFieldDefault string = "id" //默认排序字段

	LogActionPlatformAdd     int = 100 //日志行为: 平台添加
	LogActionPlatformApprove int = 101 //日志行为: 通过审核
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
	"source": {
		SourceReg:      "企业注册",
		SourcePlatform: "平台添加",
	},
	"log_action": {
		LogActionPlatformAdd:     "平台添加",
		LogActionPlatformApprove: "通过审核",
	},
}
