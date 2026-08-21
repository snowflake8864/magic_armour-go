package enum_blackmail

const (
	TypeSuffix = 1 //类型
	TypeDir    = 2 //类型
)

var EnumMapInt = map[string]map[int]string{
	"type": {
		TypeSuffix: "文件后缀",
		TypeDir:    "文件目录",
	},
	"type_en": {
		TypeSuffix: "File Extension",
		TypeDir:    "File Directory",
	},
}
var WebMapInt = map[string][]int{
	"type": {
		TypeSuffix,
		TypeDir,
	},
}
