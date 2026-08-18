package my_errors

import "fmt"

var Messages = map[string]map[string]string{
	"zh": {
		ContainerKeyAlreadyExists: "该键已经注册在容器中了",
		PublicNotExists:           "public 目录不存在",
		ConfigYamlNotExists:       "config.yml 配置文件不存在",
		ConfigGormNotExists:       "gorm_v2.yml 配置文件不存在",
		StorageLogsNotExists:      "storage/logs 目录不存在",
		ConfigInitFail:            "初始化配置文件发生错误",
		SoftLinkCreateFail:        "自动创建软连接失败,请以管理员身份运行客户端(开发环境为goland等，生产环境检查命令执行者权限), 最后一个可能：如果您是360用户，请退出360相关软件，才能保证go语言创建软连接函数： os.Symlink() 正常运行",
		SoftLinkDeleteFail:        "删除软软连接失败",

		FuncEventAlreadyExists:   "注册函数类事件失败，键名已经被注册",
		FuncEventNotRegister:     "没有找到键名对应的函数",
		FuncEventNotCall:         "注册的函数无法正确执行",
		BasePath:                 "初始化项目根目录失败",
		TokenBaseInfo:            "token最基本的格式错误,请提供一个有效的token!",
		NoAuthorization:          "token鉴权未通过，请通过token授权接口重新获取token",
		RefreshTokenFail:         "token不符合刷新条件,请通过登陆接口重新获取token!",
		ParseTokenFail:           "解析token失败",
		GormInitFail:             "Gorm 数据库驱动、连接初始化失败",
		KafkaProducerInitFail:    "Kafka生产者驱动、连接初始化失败",
		CasbinNoAuthorization:    "Casbin 鉴权未通过，请在后台检查 casbin 设置参数",
		GormNotInitGlobalPointer: "%s 数据库全局变量指针没有初始化，请在配置文件 config/gorm_v2.yml 设置 Gormv2.%s.IsInitGlobalGormMysql = 1, 并且保证数据库配置正确",

		ServerOccurredErrorMsg: "服务器内部发生代码执行错误, ",
		GinSetTrustProxyError:  "Gin 设置信任代理服务器出错",

		DbDriverNotExists:        "数据库驱动类型不存在,目前支持的数据库类型：mysql、sqlserver、postgresql，您提交数据库类型：",
		DialectorDbInitFail:      "gorm dialector 初始化失败,dbType:",
		GormDBCreateParamsNotPtr: "gorm Create 函数的参数必须是一个指针",
		GormDBUpdateParamsNotPtr: "gorm 的 Update、Save 函数的参数必须是一个指针(GinSkeleton ≥ v1.5.29 版本新增验证，为了完美支持 gorm 的所有回调函数,请在参数前面添加 & )",

		RedisInitConnFail: "初始化redis连接池失败",
		RedisAuthFail:     "Redis Auth 鉴权失败，密码错误",
		RedisGetConnFail:  "Redis 从连接池获取一个连接失败，超过最大重试次数",

		ValidatorNotExists:          "不存在的验证器",
		ValidatorTransInitFail:      "validator的翻译器初始化错误",
		ErrorNotAllParamsIsBlank:    "该接口不允许所有参数都为空,请按照接口要求提交必填参数",
		ValidatorBindParamsFail:     "验证器绑定参数失败",
		ValidatorParamsCheckFailMsg: "参数校验失败",

		TokenInvalid:         "无效的token",
		TokenNotActiveYet:    "token 尚未激活",
		TokenMalFormed:       "token 格式不正确",
		JwtTokenFormatErrMsg: "提交的 token 格式错误",
		JwtTokenMustValid:    "token为必填项,请在请求header部分提交!",

		SnowflakeGetIdFail: "获取snowflake唯一ID过程发生错误",

		FilesUploadOpenFail:           "打开文件失败，详情：",
		FilesUploadReadFail:           "读取文件32字节失败，详情：",
		FilesUploadFailMsg:            "文件上传失败, 获取上传文件发生错误!",
		FilesUploadMoreThanMaxSizeMsg: "长传文件超过系统设定的最大值,系统允许的最大值：",
		FilesUploadMimeTypeFailMsg:    "文件mime类型不允许",
		FilesUploadIsEmpty:            "不允许上传空文件",
		FilesSaveFail:                 "保存文件失败",

		CurdPleaseWaitPatiently: "请耐心等待",
		CurdDataIsEmpty:         "数据为空",
		CurdDataCannotEmpty:     "数据不能为空",
		CurdDataNotExist:        "数据不存在",
		CurdInvalidParameter:    "参数不正确",
		CurdSaveFailed:          "保存失败",
		CurdModifyFailed:        "修改失败",
		CurdDeleteFailed:        "删除失败",
		CurdRefreshFailed:       "刷新失败",
		CurdClearFailed:         "清空失败",

		CurdSliderIncorrect:         "滑块不正确",
		CurdTwoFactorAuthNotEnabled: "没有开启双因子验证",
		CurdUserNotExist:            "用户不存在",
		CurdSmsVerifyCodeNotMatch:   "短信验证码不一致",
		CurdOrgNameAlreadyExists:    "企业名称已存在",
		CurdUsernameAlreadyExists:   "用户名已存在",
		CurdPhoneAlreadyExists:      "手机号码已存在",
		CurdPasswordFormatIncorrect: "密码格式不正确",
		CurdOrgNameFormatIncorrect:  "企业名称格式不正确",
		CurdUsernameFormatIncorrect: "用户名格式不正确",
		CurdPhoneFormatIncorrect:    "手机号码格式不正确",
		CurdEmailFormatIncorrect:    "邮箱格式不正确",
		CurdIpFormatIncorrect:       "IP地址格式不正确",

		CurdEnableProcessProtectMode: "开启进程保护模式, 不能存在未设置的进程",
		CurdEnableDllProtectMode:     "开启模块保护模式, 不能存在未设置的模块",
		CurdProcSampleNotUpload:      "样本没有上传",
		CurdDllNotEnabled:            "模块功能没有开启",
		CurdTrustProcessDataError:    "信任进程数据数据错误",
		CurdDepartmentNotExist:       "组织架构不存在",
		CurdDepartmentNoToManage:     "没有需要管理的组织架构",
		CurdDepartmentMustSelect:     "必须选择一个组织架构",
		CurdDepartmentOnlyChooseYour: "只能选择属于自己的组织架构",
		CurdAgentNoToManage:          "没有需要管理的终端",
		CurdAgentSwitchMustOn:        "终端开关需要全部开启",
		CurdAgentNotExist:            "终端不存在",
		CurdAgentLinuxNotSupported:   "Linux暂不支持功能",
		CurdDirPermissionCannotEmpty: "禁止目录权限不能为空",
	},
	"en": {},
}

// GetMessage 获取指定语言的提示信息
func GetMessage(lang, key string) string {
	if msgs, ok := Messages[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// 降级到中文
	if msgs, ok := Messages["zh"]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}

// GetMessagef 支持格式化的消息获取
func GetMessagef(lang, key string, args ...interface{}) string {
	msg := GetMessage(lang, key)
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}
