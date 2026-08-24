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
	"type_en": {
		TypeUnknown:  "Unknown",
		TypeFolder:   "Folder",
		TypeFile:     "File",
		TypeVolume:   "Volume",
		TypeShortcut: "Shortcut",
		TypeLinkFile: "Symbolic Link",
		TypeMount:    "Hard Link",
	},
	"power_folder": {
		PowerRead:   "读取",
		PowerWrite:  "写入",
		PowerDelete: "删除",
		PowerRename: "重命名",
		PowerCreate: "新建",
	},
	"power_folder_en": {
		PowerRead:   "Read",
		PowerWrite:  "Write",
		PowerDelete: "Delete",
		PowerRename: "Rename",
		PowerCreate: "Create",
	},
	"is_extend": {
		ExtendYes: "是",
		ExtendNot: "否",
	},
	"is_extend_en": {
		ExtendYes: "Yes",
		ExtendNot: "No",
	},
	"search_type": {
		TypeFolder: "文件夹",
		TypeFile:   "单个文件",
	},
	"search_type_en": {
		TypeFolder: "Folder",
		TypeFile:   "File",
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

func GetPowerLabels(rw int) []string {
	var labels []string

	for _, val := range WebMapInt["power_folder"] {
		if (rw & val) > 0 {
			labels = append(labels, EnumMapInt["power_folder"][val])
		}
	}

	return labels
}

func GetPowerEnLabels(rw int) []string {
	var labels []string

	for _, val := range WebMapInt["power_folder"] {
		if (rw & val) > 0 {
			labels = append(labels, EnumMapInt["power_folder_en"][val])
		}
	}

	return labels
}
