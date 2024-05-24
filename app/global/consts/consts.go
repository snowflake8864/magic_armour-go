package consts

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
const (
	VersionNumber string = "MA_6.0.0.0"
	// 进程被结束
	ProcessKilled string = "收到信号，进程被结束"
	// 表单验证器前缀
	ValidatorPrefix              string = "Form_Validator_"
	ValidatorPrefixPlatform      string = "Platform_Validator_"
	ValidatorPrefixAgent         string = "Agent_Validator_"
	ValidatorPrefixApi           string = "Api_Validator_"
	ValidatorPrefixOrg           string = "Org_Validator_"
	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "参数校验失败"

	//服务器代码发生错误
	ServerOccurredErrorCode int    = -500100
	ServerOccurredErrorMsg  string = "服务器内部发生代码执行错误, "
	GinSetTrustProxyError   string = "Gin 设置信任代理服务器出错"

	// token相关
	JwtTokenOK            int    = 200100                      //token有效
	JwtTokenInvalid       int    = -400100                     //无效的token
	JwtTokenExpired       int    = -400101                     //过期的token
	JwtTokenFormatErrCode int    = -400102                     //提交的 token 格式错误
	JwtTokenFormatErrMsg  string = "提交的 token 格式错误"            //提交的 token 格式错误
	JwtTokenMustValid     string = "token为必填项,请在请求header部分提交!" //提交的 token 格式错误

	//SnowFlake 雪花算法
	StartTimeStamp = int64(1483228800000) //开始时间截 (2017-01-01)
	MachineIdBits  = uint(10)             //机器id所占的位数
	SequenceBits   = uint(12)             //序列所占的位数
	//MachineIdMax   = int64(-1 ^ (-1 << MachineIdBits)) //支持的最大机器id数量
	SequenceMask   = int64(-1 ^ (-1 << SequenceBits)) //
	MachineIdShift = SequenceBits                     //机器id左移位数
	TimestampShift = SequenceBits + MachineIdBits     //时间戳左移位数

	// CURD 常用业务状态码
	CurdStatusOkCode string = "000000"
	CurdStatusOkMsg  string = "Success"

	CurdCreatFailCode        int    = -400200
	CurdCreatFailMsg         string = "新增失败"
	CurdUpdateFailCode       int    = -400201
	CurdUpdateFailMsg        string = "更新失败"
	CurdDeleteFailCode       int    = -400202
	CurdDeleteFailMsg        string = "删除失败"
	CurdSelectFailCode       int    = -400203
	CurdSelectFailMsg        string = "查询无数据"
	CurdRegisterFailCode     int    = -400204
	CurdRegisterFailMsg      string = "注册失败"
	CurdLoginFailCode        int    = -400205
	CurdLoginFailMsg         string = "登录失败"
	CurdRefreshTokenFailCode int    = -400206
	CurdRefreshTokenFailMsg  string = "刷新Token失败"

	CurdBaseCode             int    = -1
	PlatformLoginExpiredCode int    = 500003  //平台:登录过期
	OrgLoginExpiredCode      int    = 500003  //平台:登录过期
	AgentParamErrorCode      int    = 1000002 //终端：参数错误
	AgentINSTALLCode         string = "410"   //终端：卸载
	AgentLoginExpiredCode    string = "401"   //终端：登录过期

	//文件上传
	FilesUploadFailCode            int    = -400250
	FilesUploadFailMsg             string = "文件上传失败, 获取上传文件发生错误!"
	FilesUploadMoreThanMaxSizeCode int    = -400251
	FilesUploadMoreThanMaxSizeMsg  string = "长传文件超过系统设定的最大值,系统允许的最大值："
	FilesUploadMimeTypeFailCode    int    = -400252
	FilesUploadMimeTypeFailMsg     string = "文件mime类型不允许"
	FilesUploadIsEmpty             string = "不允许上传空文件"
)
