package my_errors

const (
	//系统部分
	ContainerKeyAlreadyExists string = "ContainerKeyAlreadyExists"
	PublicNotExists           string = "PublicNotExists"
	ConfigYamlNotExists       string = "ConfigYamlNotExists"
	ConfigGormNotExists       string = "ConfigGormNotExists"
	StorageLogsNotExists      string = "StorageLogsNotExists"
	ConfigInitFail            string = "ConfigInitFail"
	SoftLinkCreateFail        string = "SoftLinkCreateFail"
	SoftLinkDeleteFail        string = "SoftLinkDeleteFail"

	FuncEventAlreadyExists   string = "FuncEventAlreadyExists"
	FuncEventNotRegister     string = "FuncEventNotRegister"
	FuncEventNotCall         string = "FuncEventNotCall"
	BasePath                 string = "BasePath"
	TokenBaseInfo            string = "TokenBaseInfo"
	NoAuthorization          string = "NoAuthorization"
	RefreshTokenFail         string = "RefreshTokenFail"
	ParseTokenFail           string = "ParseTokenFail"
	GormInitFail             string = "GormInitFail"
	KafkaProducerInitFail    string = "KafkaProducerInitFail"
	CasbinNoAuthorization    string = "CasbinNoAuthorization"
	GormNotInitGlobalPointer string = "GormNotInitGlobalPointer"

	ServerOccurredErrorCode int    = -500100
	ServerOccurredErrorMsg  string = "ServerOccurredErrorMsg"
	GinSetTrustProxyError   string = "GinSetTrustProxyError"

	// 数据库部分
	DbDriverNotExists        string = "DbDriverNotExists"
	DialectorDbInitFail      string = "DialectorDbInitFail"
	GormDBCreateParamsNotPtr string = "GormDBCreateParamsNotPtr"
	GormDBUpdateParamsNotPtr string = "GormDBUpdateParamsNotPtr"

	// Redis部分
	RedisInitConnFail string = "RedisInitConnFail"
	RedisAuthFail     string = "RedisAuthFail"
	RedisGetConnFail  string = "RedisGetConnFail"

	// 表单参数验证器
	ValidatorNotExists           string = "ValidatorNotExists"
	ValidatorTransInitFail       string = "ValidatorTransInitFail"
	ErrorNotAllParamsIsBlank     string = "ErrorNotAllParamsIsBlank"
	ValidatorBindParamsFail      string = "ValidatorBindParamsFail"
	ValidatorParamsCheckFailCode int    = -400300
	ValidatorParamsCheckFailMsg  string = "ValidatorParamsCheckFailMsg"

	// Token部分
	TokenInvalid          string = "TokenInvalid"
	TokenNotActiveYet     string = "TokenNotActiveYet"
	TokenMalFormed        string = "TokenMalFormed"
	JwtTokenOK            int    = 200100                 //token有效
	JwtTokenInvalid       int    = -400100                //无效的token
	JwtTokenExpired       int    = -400101                //过期的token
	JwtTokenFormatErrCode int    = -400102                //提交的 token 格式错误
	JwtTokenFormatErrMsg  string = "JwtTokenFormatErrMsg" //提交的 token 格式错误
	JwtTokenMustValid     string = "JwtTokenMustValid"    //提交的 token 格式错误

	// Snowflake
	SnowflakeGetIdFail string = "SnowflakeGetIdFail"

	// 文件上传
	FilesUploadOpenFail            string = "FilesUploadOpenFail"
	FilesUploadReadFail            string = "FilesUploadReadFail"
	FilesUploadFailCode            int    = -400250
	FilesUploadFailMsg             string = "FilesUploadFailMsg"
	FilesUploadMoreThanMaxSizeCode int    = -400251
	FilesUploadMoreThanMaxSizeMsg  string = "FilesUploadMoreThanMaxSizeMsg"
	FilesUploadMimeTypeFailCode    int    = -400252
	FilesUploadMimeTypeFailMsg     string = "FilesUploadMimeTypeFailMsg"
	FilesUploadIsEmpty             string = "FilesUploadIsEmpty"
	FilesSaveFail                  string = "FilesSaveFail"

	// CURD 常用业务状态码
	CurdStatusOkCode         string = "000000"
	CurdStatusOkMsg          string = "OK"
	CurdBaseCode             int    = -1
	PlatformLoginExpiredCode int    = 500003    //平台:登录过期
	OrgLoginExpiredCode      int    = 500003    //平台:登录过期
	AgentParamErrorCode      string = "1000002" //终端：参数错误
	AgentINSTALLCode         string = "410"     //终端：卸载
	AgentLoginExpiredCode    string = "401"     //终端：登录过期

	CurdPleaseWaitPatiently string = " CurdPleaseWaitPatiently"
	CurdDataIsEmpty         string = " CurdDataIsEmpty"
	CurdDataCannotEmpty     string = " CurdDataCannotEmpty"
	CurdDataNotExist        string = " CurdDataNotExist"
	CurdInvalidParameter    string = " CurdInvalidParameter"
	CurdAddFailed           string = " CurdAddFailed"
	CurdSaveFailed          string = " CurdSaveFailed"
	CurdModifyFailed        string = " CurdModifyFailed"
	CurdDeleteFailed        string = " CurdDeleteFailed"
	CurdRefreshFailed       string = " CurdRefreshFailed"
	CurdClearFailed         string = " CurdClearFailed"

	CurdSliderIncorrect         string = " CurdSliderIncorrect"
	CurdTwoFactorAuthNotEnabled string = " CurdTwoFactorAuthNotEnabled"
	CurdUserNotExist            string = " CurdUserNotExist"
	CurdSmsVerifyCodeNotMatch   string = " CurdSmsVerifyCodeNotMatch"
	CurdOrgNameAlreadyExists    string = " CurdOrgNameAlreadyExists"
	CurdUsernameAlreadyExists   string = " CurdUsernameAlreadyExists"
	CurdPhoneAlreadyExists      string = " CurdPhoneAlreadyExists"

	CurdEnableProcessProtectMode string = " CurdEnableProcessProtectMode"
	CurdEnableDllProtectMode     string = " CurdEnableDllProtectMode"
	CurdProcSampleNotUpload      string = " CurdProcSampleNotUpload"
	CurdDllNotEnabled            string = " CurdDllNotEnabled"
	CurdTrustProcessDataError    string = " TrustProcessDataError"
	CurdDepartmentNotExist       string = " CurdDepartmentNotExist"
	CurdDepartmentNoToManage     string = " CurdDepartmentNoToManage"
	CurdDepartmentMustSelect     string = " CurdDepartmentMustSelect"
	CurdDepartmentOnlyChooseYour string = " CurdDepartmentOnlyChooseYour"
	CurdAgentNoToManage          string = " CurdAgentNoToManage"
	CurdAgentSwitchMustOn        string = " CurdAgentSwitchMustOn"
	CurdAgentNotExist            string = " CurdAgentNotExist"
	CurdAgentLinuxNotSupported   string = " CurdAgentLinuxNotSupported"
)
