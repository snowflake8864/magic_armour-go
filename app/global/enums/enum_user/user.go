package enum_user

const (
	StatusVerifying int = 1 //状态:审核中
	StatusNormal    int = 2 //状态:正常
	StatusLocking   int = 3 //状态:锁定中

	UserFieldDefault string = "id" //默认排序字段

)

var EnumMap = map[string]map[int]string{

	"status": {
		StatusVerifying: "审核中",
		StatusNormal:    "企业",
		StatusLocking:   "锁定",
	},
}
