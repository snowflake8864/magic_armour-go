package consts

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
const (
	// 进程被结束
	ProcessKilled string = "收到信号，进程被结束"
	// 表单验证器前缀
	ValidatorPrefix         string = "Form_Validator_"
	ValidatorPrefixPlatform string = "Platform_Validator_"
	ValidatorPrefixAgent    string = "Agent_Validator_"
	ValidatorPrefixApi      string = "Api_Validator_"
	ValidatorPrefixOrg      string = "Org_Validator_"

	//SnowFlake 雪花算法
	StartTimeStamp = int64(1483228800000) //开始时间截 (2017-01-01)
	MachineIdBits  = uint(10)             //机器id所占的位数
	SequenceBits   = uint(12)             //序列所占的位数
	//MachineIdMax   = int64(-1 ^ (-1 << MachineIdBits)) //支持的最大机器id数量
	SequenceMask   = int64(-1 ^ (-1 << SequenceBits)) //
	MachineIdShift = SequenceBits                     //机器id左移位数
	TimestampShift = SequenceBits + MachineIdBits     //时间戳左移位数

	LanguageField string = "X-Language" //支持语言
	LanguageZH    string = "zh"         //支持语言
	LanguageEn    string = "en"         //支持语言
)
