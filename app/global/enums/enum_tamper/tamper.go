package enum_tamper

const (
	RootDir   string = "/" //根目录
	Separator string = "|" //分割字符

	TypeUnknown  = 0 //文件类型
	TypeFolder   = 1 //文件类型
	TypeFile     = 2 //文件类型
	TypeVolume   = 3 //文件类型
	TypeShortcut = 4 //文件类型
	TypeLinkFile = 5 //文件类型
	TypeMount    = 6 //文件类型

	PowerRead   = 1  //权限
	PowerWrite  = 2  //权限
	PowerDelete = 4  //权限
	PowerRename = 8  //权限
	PowerCreate = 16 //权限

	ExtendYes = 1 //继承关系
	ExtendNot = 0 //继承关系
)

var EnumMapInt = map[string]map[int]string{
	"type": {
		TypeUnknown:  "未知",
		TypeFolder:   "文件夹",
		TypeFile:     "单个文件",
		TypeVolume:   "券",
		TypeShortcut: "快捷方式",
		TypeLinkFile: "软连接",
		TypeMount:    "硬链接",
	},
	"power_folder": {
		PowerRead:   "读取",
		PowerWrite:  "写入",
		PowerDelete: "删除",
		PowerRename: "重命名",
		PowerCreate: "新建",
	},
	"power_file": {
		PowerRead:   "读取",
		PowerWrite:  "写入",
		PowerDelete: "删除",
		PowerRename: "重命名",
	},
	"is_extend": {
		ExtendYes: "是",
		ExtendNot: "否",
	},
	"search_type": {
		TypeFolder: "文件夹",
		TypeFile:   "单个文件",
	},
}

var WebMapInt = map[string][]int{
	"search_type": {
		TypeFolder,
		TypeFile,
	},
	"is_extend": {
		ExtendYes,
		ExtendNot,
	},
	"power_folder": {
		PowerRead,
		PowerWrite,
		PowerDelete,
		PowerRename,
		PowerCreate,
	},
	"power_file": {
		PowerRead,
		PowerWrite,
		PowerDelete,
		PowerRename,
	},
}
