package enum_usb

const (
	StatusDefault = 1 //设备状态
	StatusWhite   = 2 //设备状态
	StatusBlack   = 3 //设备状态

	PlugStatusOut = 1 //插拔状态
	PlugStatusIn  = 2 //插拔状态
)

var EnumMapInt = map[string]map[int]string{
	"status": {
		StatusDefault: "未确认",
		StatusWhite:   "信任",
		StatusBlack:   "禁止",
	},
	"plug_status": {
		PlugStatusOut: "拔出",
		PlugStatusIn:  "插入",
	},
}

var WebMapInt = map[string][]int{
	"edit_status": {
		StatusWhite,
		StatusBlack,
	},
}
