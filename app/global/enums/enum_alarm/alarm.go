package enum_alarm

import "goskeleton/app/global/enums/enum_port"

const (
	HandleStatusNot int = 1 //处理状态: 未处理
	HandleStatusYes int = 2 //处理状态: 已处理

	LevelCommon    int = 1 //告警级别: 普通
	LevelImportant int = 2 // 告警级别: 重要
	LevelSerious   int = 3 //告警级别: 严重

	ModuleAgentError int = 3  //告警功能: 终端错误
	ModuleProcess    int = 10 //告警功能: 进程
	ModuleDll        int = 11 //告警功能: 模块
	ModuleTamper     int = 12 //告警功能: 防篡改
	ModuleBlackmail  int = 13 //告警功能: 勒索防护
	ModulePort       int = 14 //告警功能: 端口虚开
	ModuleUsb        int = 15 //告警功能: 外设
	ModulePlug       int = 16 //告警功能: 动态封堵
	ModuleOutreach   int = 17 //告警功能: 外联
	ModuleLogin      int = 18 // 告警功能: 登录破解
	ModuleSelf       int = 19 //告警功能: 自保

	/**   --------进程告警 S--------- */

	TypeProcessMonitorUnknownRun int = 1001 //不明进程启动告警
	TypeProcessMonitorBlackRun   int = 1002 //黑名单启动告警

	TypeProcessProtectUnknownRun int = 1101 //不明进程启动阻止告警
	TypeProcessProtectBlackRun   int = 1102 //黑名单进程启动阻止告警
	TypeProcessProtectUnknownEnd int = 1103 //不明进程终止告警
	TypeProcessProtectBlackEnd   int = 1104 //黑名单进程终止告警
	/**   --------进程告警 E--------- */

	/**   --------模块告警 S--------- */

	TypeDllMonitorUnknownRun int = 1201 //不明模块启动告警
	TypeDllMonitorBlackRun   int = 1202 //黑名单模块启动告警

	TypeDllProtectUnknownRun int = 1301 //不明模块启动阻止告警
	TypeDllProtectBlackRun   int = 1302 //黑名单模块启动阻止告警
	/**   --------模块告警 E--------- */

	/**   --------防篡改告警 S--------- */

	TypeTamperMonitorFolderCreate int = 2001 //文件夹新建告警
	TypeTamperMonitorFolderDelete int = 2002 //文件夹删除告警
	TypeTamperMonitorFolderUpdate int = 2003 //文件夹修改告警
	TypeTamperMonitorFolderRead   int = 2004 //文件夹读取告警
	TypeTamperMonitorFolderRename int = 2005 //文件夹重命名告警

	TypeTamperProtectFolderCreate int = 2101 //文件夹新建阻止告警
	TypeTamperProtectFolderDelete int = 2102 //文件夹删除阻止告警
	TypeTamperProtectFolderUpdate int = 2103 //文件夹修改阻止告警
	TypeTamperProtectFolderRead   int = 2104 //文件夹读取阻止告警
	TypeTamperProtectFolderRename int = 2105 //文件夹重命名阻止告警

	TypeTamperMonitorFileCreate int = 3001 //文件新建告警
	TypeTamperMonitorFileDelete int = 3002 //文件删除告警
	TypeTamperMonitorFileUpdate int = 3003 //文件修改告警
	TypeTamperMonitorFileRead   int = 3004 //文件读取告警
	TypeTamperMonitorFileRename int = 3005 //文件重命名告警

	TypeTamperProtectFileCreate int = 3101 //文件新建阻止告警
	TypeTamperProtectFileDelete int = 3102 //文件删除阻止告警
	TypeTamperProtectFileUpdate int = 3103 //文件修改阻止告警
	TypeTamperProtectFileRead   int = 3104 //文件读取阻止告警
	TypeTamperProtectFileRename int = 3105 //文件重命名阻止告警
	/**   --------防篡改告警 E--------- */

	/**   ---------勒索防护告警 S--------- */

	TypeBlackmailProtect int = 3201 //异常进程访问阻止
	TypeBlackmailMonitor int = 3202 //异常进程访问
	/**   ---------勒索防护告警 E--------- */

	/**   ---------外联告警 S--------- */

	TypeOutreachNotInWhite int = 4001 //Ip不在白名单中
	/**   ---------外设告警 E--------- */

	/**   ---------动态封堵告警 S--------- */

	TypePlugDynamicFromPort  int = 4100 //端口虚开
	TypePlugDynamicFromLogin int = 4101 //口令破解
	/**   ---------动态封堵告警 E--------- */

	/**   ---------登录口令暴力破解告警 S--------- */

	TypeLoginSolve int = 4200 //登录口令暴力破解
	/**   ---------登录口令暴力破解告警 E--------- */

	/**   ---------端口告警 S--------- */

	TypePortVirtual int = 4300 //登录口令暴力破解
	/**   ---------端口告警 E--------- */

	/**   ---------自保告警 S--------- */

	TypeSelfProtectFileCreate     int = 4700 //新建文件告警
	TypeSelfProtectFileDelete     int = 4701 //删除文件告警
	TypeSelfProtectFileUpdate     int = 4702 //修改文件止告警
	TypeSelfProtectFileRename     int = 4703 //重命名文件告警
	TypeSelfProtectFileMove       int = 4704 //移动文件告警
	TypeSelfProtectProcessKill    int = 4710 //杀死进程
	TypeSelfProtectServiceStop    int = 4720 //停止服务
	TypeSelfProtectServiceRestart int = 4721 //重启服务
	TypeSelfProtectServiceDisable int = 4722 //禁止服务自启动

	TypeSelfProtectRegeditDelete int = 4730 //删除注册表信息
	TypeSelfProtectRegeditUpdate int = 4731 //修改注册表信息
	/**   ---------自保告警 E--------- */

	/**   ---------agent错误告警 S--------- */

	TypeAgentUpgradeError int = 9001 //升级错误
	TypeAgentDestroy      int = 9002 //Agent被破坏告警
	/**   ---------agent错误告警 E--------- */

	/**   ---------外设告警 S--------- */

	TypeUsbMonitorUnknownAccess  int = 9003 //未知usb设备接入
	TypeUsbMonitorDisabledAccess int = 9004 //禁止usb设备接入
	TypeUsbProtectUnknownAccess  int = 9005 //未知usb设备接入阻断
	TypeUsbProtectDisabledAccess int = 9006 //禁止usb设备接入阻断
	TypeUsbDisabledExtract       int = 9007 //禁止usb设备拔出
	TypeUsbUnknownExtract        int = 9008 //未知usb设备拔出
	/**   ---------外设告警 E--------- */
)

var EnumMapInt = map[string]map[int]string{
	"handle_status": {
		HandleStatusNot: "未处理",
		HandleStatusYes: "已处理",
	},
	"level": {
		LevelCommon:    "普通",
		LevelImportant: "重要",
		LevelSerious:   "严重",
	},
	"module": {
		ModuleProcess:   "进程异常告警",
		ModuleDll:       "模块异常告警",
		ModuleTamper:    "防篡改告警",
		ModuleBlackmail: "勒索告警",
		ModulePort:      "端口虚开告警",
		ModulePlug:      "动态封堵告警",
		ModuleUsb:       "外接设备告警",
		ModuleOutreach:  "外联告警",
		ModuleLogin:     "口令破解告警",
		ModuleSelf:      "自保告警",
	},
	"type": {
		TypeProcessMonitorUnknownRun: "不明进程启动告警",
		TypeProcessMonitorBlackRun:   "黑名单进程启动告警",
		TypeProcessProtectUnknownRun: "不明进程启动阻止告警",
		TypeProcessProtectBlackRun:   "黑名单进程启动阻止告警",
		TypeProcessProtectUnknownEnd: "不明进程终止告警",
		TypeProcessProtectBlackEnd:   "黑名单进程终止告警",

		TypeDllMonitorUnknownRun: "不明模块启动告警",
		TypeDllMonitorBlackRun:   "黑名单模块启动告警",
		TypeDllProtectUnknownRun: "不明模块启动阻止告警",
		TypeDllProtectBlackRun:   "黑名单模块启动阻止告警",

		TypeTamperMonitorFolderCreate: "文件夹新建告警",
		TypeTamperMonitorFolderDelete: "文件夹删除告警",
		TypeTamperMonitorFolderUpdate: "文件夹写入告警",
		TypeTamperMonitorFolderRead:   "文件夹读取告警",
		TypeTamperMonitorFolderRename: "文件夹重命名告警",
		TypeTamperProtectFolderCreate: "文件夹新建阻止告警",
		TypeTamperProtectFolderDelete: "文件夹删除阻止告警",
		TypeTamperProtectFolderUpdate: "文件夹写入阻止告警",
		TypeTamperProtectFolderRead:   "文件夹读取阻止告警",
		TypeTamperProtectFolderRename: "文件夹重命名阻止告警",
		TypeTamperMonitorFileCreate:   "文件新建告警",
		TypeTamperMonitorFileDelete:   "文件删除告警",
		TypeTamperMonitorFileUpdate:   "文件写入告警",
		TypeTamperMonitorFileRead:     "文件读取告警",
		TypeTamperMonitorFileRename:   "文件重命名告警",
		TypeTamperProtectFileCreate:   "文件新建阻止告警",
		TypeTamperProtectFileDelete:   "文件删除阻止告警",
		TypeTamperProtectFileUpdate:   "文件写入阻止告警",
		TypeTamperProtectFileRead:     "文件读取阻止告警",
		TypeTamperProtectFileRename:   "文件重命名阻止告警",

		TypeAgentUpgradeError: "升级错误",
		TypeAgentDestroy:      "Agent被破坏告警",

		TypeBlackmailProtect: "异常进程访问阻止",
		TypeBlackmailMonitor: "异常进程访问",

		TypeOutreachNotInWhite: "外联告警",

		TypePlugDynamicFromPort:  "端口虚开",
		TypePlugDynamicFromLogin: "口令破解",

		TypeLoginSolve: "登录口令暴力破解",

		TypePortVirtual: "端口虚开",

		TypeSelfProtectFileCreate:     "新建文件告警",
		TypeSelfProtectFileDelete:     "删除文件告警",
		TypeSelfProtectFileUpdate:     "修改文件止告警",
		TypeSelfProtectFileRename:     "重命名文件告警",
		TypeSelfProtectFileMove:       "移动文件告警",
		TypeSelfProtectProcessKill:    "杀死进程",
		TypeSelfProtectServiceStop:    "停止服务",
		TypeSelfProtectServiceRestart: "重启服务",
		TypeSelfProtectServiceDisable: "禁止服务自启动",
		TypeSelfProtectRegeditDelete:  "删除注册表信息",
		TypeSelfProtectRegeditUpdate:  "修改注册表信息",

		TypeUsbMonitorUnknownAccess:  "未知usb设备接入",
		TypeUsbMonitorDisabledAccess: "禁止usb设备接入",
		TypeUsbProtectUnknownAccess:  "未知usb设备接入阻断",
		TypeUsbProtectDisabledAccess: "禁止usb设备接入阻断",
		TypeUsbDisabledExtract:       "禁止usb设备拔出",
		TypeUsbUnknownExtract:        "未知usb设备拔出",
	},
}

var WebMapInt = map[string][]int{
	"handle_status": {
		HandleStatusNot,
		HandleStatusYes,
	},
	"level": {
		LevelCommon,
		LevelImportant,
		LevelSerious,
	},
	"module": {
		ModuleProcess,
		ModuleDll,
		ModuleTamper,
		ModuleBlackmail,
		ModulePort,
		ModulePlug,
		ModuleUsb,
		ModuleOutreach,
		ModuleLogin,
		ModuleSelf,
	},
	"type_process": {
		TypeProcessMonitorUnknownRun,
		TypeProcessMonitorBlackRun,
		TypeProcessProtectUnknownRun,
		TypeProcessProtectBlackRun,
		TypeProcessProtectUnknownEnd,
		TypeProcessProtectBlackEnd,
	},
	"type_dll": {
		TypeDllMonitorUnknownRun,
		TypeDllMonitorBlackRun,
		TypeDllProtectUnknownRun,
		TypeDllProtectBlackRun,
	},
	"type_tamper": {
		TypeTamperMonitorFolderCreate,
		TypeTamperMonitorFolderDelete,
		TypeTamperMonitorFolderUpdate,
		TypeTamperMonitorFolderRead,
		TypeTamperMonitorFolderRename,
		TypeTamperProtectFolderCreate,
		TypeTamperProtectFolderDelete,
		TypeTamperProtectFolderUpdate,
		TypeTamperProtectFolderRead,
		TypeTamperProtectFolderRename,
		TypeTamperMonitorFileCreate,
		TypeTamperMonitorFileDelete,
		TypeTamperMonitorFileUpdate,
		TypeTamperMonitorFileRead,
		TypeTamperMonitorFileRename,
		TypeTamperProtectFileCreate,
		TypeTamperProtectFileDelete,
		TypeTamperProtectFileUpdate,
		TypeTamperProtectFileRead,
		TypeTamperProtectFileRename,
	},
	"type_blackmail": {
		TypeBlackmailMonitor,
		TypeBlackmailProtect,
	},
	"type_usb": {
		TypeUsbMonitorUnknownAccess,
		TypeUsbMonitorDisabledAccess,
		TypeUsbProtectUnknownAccess,
		TypeUsbProtectDisabledAccess,
		TypeUsbDisabledExtract,
		TypeUsbUnknownExtract,
	},
	"type_self": {
		TypeSelfProtectFileCreate,
		TypeSelfProtectFileDelete,
		TypeSelfProtectFileUpdate,
		TypeSelfProtectFileRename,
		TypeSelfProtectFileMove,
		TypeSelfProtectProcessKill,
		TypeSelfProtectServiceStop,
		TypeSelfProtectServiceRestart,
		TypeSelfProtectServiceDisable,
		TypeSelfProtectRegeditDelete,
		TypeSelfProtectRegeditUpdate,
	},
}

var MapTypeModule = map[int]int{
	TypeProcessMonitorUnknownRun: ModuleProcess,
	TypeProcessMonitorBlackRun:   ModuleProcess,
	TypeProcessProtectUnknownRun: ModuleProcess,
	TypeProcessProtectBlackRun:   ModuleProcess,
	TypeProcessProtectUnknownEnd: ModuleProcess,
	TypeProcessProtectBlackEnd:   ModuleProcess,

	TypeDllMonitorUnknownRun: ModuleDll,
	TypeDllMonitorBlackRun:   ModuleDll,
	TypeDllProtectUnknownRun: ModuleDll,
	TypeDllProtectBlackRun:   ModuleDll,

	TypeTamperMonitorFolderCreate: ModuleTamper,
	TypeTamperMonitorFolderDelete: ModuleTamper,
	TypeTamperMonitorFolderUpdate: ModuleTamper,
	TypeTamperMonitorFolderRead:   ModuleTamper,
	TypeTamperMonitorFolderRename: ModuleTamper,
	TypeTamperProtectFolderCreate: ModuleTamper,
	TypeTamperProtectFolderDelete: ModuleTamper,
	TypeTamperProtectFolderUpdate: ModuleTamper,
	TypeTamperProtectFolderRead:   ModuleTamper,
	TypeTamperProtectFolderRename: ModuleTamper,
	TypeTamperMonitorFileCreate:   ModuleTamper,
	TypeTamperMonitorFileDelete:   ModuleTamper,
	TypeTamperMonitorFileUpdate:   ModuleTamper,
	TypeTamperMonitorFileRead:     ModuleTamper,
	TypeTamperMonitorFileRename:   ModuleTamper,
	TypeTamperProtectFileCreate:   ModuleTamper,
	TypeTamperProtectFileDelete:   ModuleTamper,
	TypeTamperProtectFileUpdate:   ModuleTamper,
	TypeTamperProtectFileRead:     ModuleTamper,
	TypeTamperProtectFileRename:   ModuleTamper,

	TypeBlackmailMonitor: ModuleBlackmail,
	TypeBlackmailProtect: ModuleBlackmail,

	TypeUsbMonitorUnknownAccess:  ModuleUsb,
	TypeUsbMonitorDisabledAccess: ModuleUsb,
	TypeUsbProtectUnknownAccess:  ModuleUsb,
	TypeUsbProtectDisabledAccess: ModuleUsb,
	TypeUsbDisabledExtract:       ModuleUsb,
	TypeUsbUnknownExtract:        ModuleUsb,

	TypeOutreachNotInWhite: ModuleOutreach,

	TypeSelfProtectFileCreate:     ModuleSelf,
	TypeSelfProtectFileDelete:     ModuleSelf,
	TypeSelfProtectFileUpdate:     ModuleSelf,
	TypeSelfProtectFileRename:     ModuleSelf,
	TypeSelfProtectFileMove:       ModuleSelf,
	TypeSelfProtectProcessKill:    ModuleSelf,
	TypeSelfProtectServiceStop:    ModuleSelf,
	TypeSelfProtectServiceRestart: ModuleSelf,
	TypeSelfProtectServiceDisable: ModuleSelf,
	TypeSelfProtectRegeditDelete:  ModuleSelf,
	TypeSelfProtectRegeditUpdate:  ModuleSelf,

	TypeAgentUpgradeError: ModuleAgentError,
}

// GetModuleByType 通过类型获取对应模块
func GetModuleByType(t int) int {
	if val, isOk := MapTypeModule[t]; isOk {
		return val
	} else {
		return 0
	}
}

// GetLevelByWeight 通过类型获取对应模块
func GetLevelByWeight(w int) int {
	if w == enum_port.AlarmLevel1 || w == enum_port.AlarmLevel2 {
		return LevelCommon
	} else if w == enum_port.AlarmLevel3 || w == enum_port.AlarmLevel4 {
		return LevelImportant
	} else if w == enum_port.AlarmLevel5 {
		return LevelSerious
	} else {
		return 0
	}
}
