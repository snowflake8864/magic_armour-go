package enum_threaten

const (
	CenterSupportNumber int = 150 //威胁情报中心支持数量:最大300,建议200

	LevelUnknown      int = 0  //威胁情报中心等级：未知
	LevelWhite        int = 1  //威胁情报中心等级：白名单
	LevelDefaultOne   int = 2  //威胁情报中心等级：1级
	LevelDefaultTwo   int = 3  //威胁情报中心等级：2级
	LevelDefaultThree int = 4  //威胁情报中心等级：3级
	LevelDefaultFour  int = 5  //威胁情报中心等级：4级
	LevelDefaultFive  int = 6  //威胁情报中心等级：5级
	LevelDefaultSix   int = 7  //威胁情报中心等级：6级
	LevelDefaultSeven int = 8  //威胁情报中心等级：7级
	LevelDefaultEight int = 9  //威胁情报中心等级：8级
	LevelDefaultNine  int = 10 //威胁情报中心等级：9级
	LevelBlack        int = 11 //威胁情报中心等级：黑名单
	LevelDefault      int = 99 //威胁情报中心等级： 未同步

	SourceThreatenCenter int = 1 //等级来源:威胁情报中心
	SourceImport         int = 2 //等级来源:人工导入
	SourceAgentUpload    int = 3 //等级来源:内部生成

	PushStatusNot int = 1 //推送状态：未推送
	PushStatusYes int = 2 //推送状态：已推送

	DataFromUnknown int = 1 //数据来源：未知
	DataFromProcess int = 2 //数据来源：进程
	DataFromDll     int = 3 //数据来源：模块

	CheckSourceVeda     int = 1 //验证结果数据来源：威胁情报
	CheckSourceInternet int = 2 //验证结果数据来源：互联网

)

var EnumMap = map[string]map[int]string{
	"level": {
		LevelUnknown:      "未知",
		LevelWhite:        "白名单",
		LevelDefaultOne:   "1级",
		LevelDefaultTwo:   "2级",
		LevelDefaultThree: "3级",
		LevelDefaultFour:  "4级",
		LevelDefaultFive:  "5级",
		LevelDefaultSix:   "6级",
		LevelDefaultSeven: "7级",
		LevelDefaultEight: "8级",
		LevelDefaultNine:  "9级",
		LevelBlack:        "黑名单",
		LevelDefault:      "未同步",
	},
	"source": {
		SourceThreatenCenter: "威胁情报",
		SourceImport:         "人工导入",
		SourceAgentUpload:    "终端上传",
	},
	"push_status": {
		PushStatusNot: "未推送",
		PushStatusYes: "已推送",
	},
	"check_source": {
		CheckSourceVeda:     "威胁情报",
		CheckSourceInternet: "云查杀",
	},
}

var WebMapInt = map[string][]int{
	"level": {
		LevelUnknown,
		LevelWhite,
		LevelDefaultOne,
		LevelDefaultTwo,
		LevelDefaultThree,
		LevelDefaultFour,
		LevelDefaultFive,
		LevelDefaultSix,
		LevelDefaultSeven,
		LevelDefaultEight,
		LevelDefaultNine,
		LevelBlack,
		LevelDefault,
	},
}

// IsWhite 是否白名单
func IsWhite(level int) bool {
	if level == LevelWhite {
		return true
	}
	return false
}

// IsBlack 是否黑名单
func IsBlack(level int) bool {
	if level == LevelBlack {
		return true
	}
	return false
}

// HasLevel 是否有威胁等级
func HasLevel(level int) bool {
	if level == LevelUnknown || level == LevelDefault {
		return false
	}
	return true
}

// HasEffect 是否会产生影响
func HasEffect(level int) bool {
	if level == LevelWhite || level == LevelBlack {
		return false
	}
	return true
}
