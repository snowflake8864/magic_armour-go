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
	ValidatorNotExists       string = "ValidatorNotExists"
	ValidatorTransInitFail   string = "ValidatorTransInitFail"
	ErrorNotAllParamsIsBlank string = "ErrorNotAllParamsIsBlank"
	ValidatorBindParamsFail  string = "ValidatorBindParamsFail"

	// Token部分
	TokenInvalid      string = "TokenInvalid"
	TokenNotActiveYet string = "TokenNotActiveYet"
	TokenMalFormed    string = "TokenMalFormed"

	// Snowflake
	SnowflakeGetIdFail string = "SnowflakeGetIdFail"

	// 文件上传
	FilesUploadOpenFail            string = "FilesUploadOpenFail"
	FilesUploadReadFail            string = "FilesUploadReadFail"
	FilesUploadFailCode            int    = -400250
	FilesUploadFailMsg             string = "文件上传失败, 获取上传文件发生错误!"
	FilesUploadMoreThanMaxSizeCode int    = -400251
	FilesUploadMoreThanMaxSizeMsg  string = "长传文件超过系统设定的最大值,系统允许的最大值："
	FilesUploadMimeTypeFailCode    int    = -400252
	FilesUploadMimeTypeFailMsg     string = "文件mime类型不允许"
	FilesUploadIsEmpty             string = "不允许上传空文件"
)
