package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/authorization/middleware_auth_agent"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/gin_release"
	"net/http"
)

// 该路由主要设置终端路由

func InitAgentRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)

		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	// 设置可信任的代理服务器列表,gin (2021-11-24发布的v1.7.7版本之后出的新功能)
	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Agent接口: "+consts.VersionNumber)
	})

	//  创建一个接口路由组
	vAgent := router.Group("/v1/").Use(authorization.RequestParams).Use(middleware_auth_agent.SaveAccess())
	{
		vAgent.POST("auth", validatorFactory.Create(consts.ValidatorPrefixAgent+"AuthLogin"))
		vAgent.POST("getuserid", validatorFactory.Create(consts.ValidatorPrefixAgent+"AuthMakeUserid"))

		// 【需要token】中间件验证的路由
		vAgent.Use(middleware_auth_agent.CheckToken())
		{
			//终端相关
			vAgent.Any("getconf", (&agent.AgentController{}).GetConfig)                                            // 获取配置信息
			vAgent.POST("getinfo", (&agent.AgentController{}).GetInfo)                                             //获取终端信息(调试用)
			vAgent.POST("getSyslogConf", (&agent.AgentController{}).GetConfig)                                     //获取syslog配置(不在使用)
			vAgent.Any("download", (&agent.AgentController{}).Upgrade)                                             //升级
			vAgent.Any("uninstall", (&agent.AgentController{}).Uninstall)                                          //卸载
			vAgent.Any("gettrustdir", (&agent.AgentController{}).GetTrustDir)                                      //获取信任目录
			vAgent.POST("error/upload", validatorFactory.Create(consts.ValidatorPrefixAgent+"ErrorFileUpload"))    //上传转储(错误) 文件
			vAgent.POST("puthardwareinfo", validatorFactory.Create(consts.ValidatorPrefixAgent+"HardwareMonitor")) //硬件资源监控

			//任务相关
			vAgent.Any("gettask", (&agent.TaskController{}).Get)
			vAgent.Any("closetask", validatorFactory.Create(consts.ValidatorPrefixAgent+"TaskClose"))

			//进程相关
			vAgent.POST("upload/suffix/exe", validatorFactory.Create(consts.ValidatorPrefixAgent+"ProcessAllPut"))     //上传全部进程
			vAgent.POST("uploadproc", validatorFactory.Create(consts.ValidatorPrefixAgent+"ProcessRunningPut"))        //上传运行中进程
			vAgent.POST("autouploadprocess", validatorFactory.Create(consts.ValidatorPrefixAgent+"ProcessAutoPut"))    //主动上传进程
			vAgent.Any("getprocwl", (&agent.ProcessController{}).GetWhite)                                             //下载进程白名单
			vAgent.Any("getprocbl", (&agent.ProcessController{}).GetBlack)                                             //下载进程黑名单
			vAgent.Any("module/process", (&agent.ProcessController{}).GetAssign)                                       //指定进程:下载信息
			vAgent.POST("upload/proc/module", validatorFactory.Create(consts.ValidatorPrefixAgent+"ProcessAssignPut")) //指定进程:上传对应模块

			//模块相关
			vAgent.POST("upload/suffix/dll", validatorFactory.Create(consts.ValidatorPrefixAgent+"DllAllPut"))         //上传全部模块
			vAgent.POST("upload/allproc/module", validatorFactory.Create(consts.ValidatorPrefixAgent+"DllRunningPut")) //上传运行中模块
			vAgent.POST("auto/upload/module", validatorFactory.Create(consts.ValidatorPrefixAgent+"DllAutoPut"))       //主动上传模块
			vAgent.Any("module/white/list", (&agent.DllController{}).GetWhite)                                         //下载模块白名单
			vAgent.Any("module/black/list", (&agent.DllController{}).GetBlack)                                         //下载模块黑名单

			//进程或模块样本
			vAgent.Any("getdraw", (&agent.ProcSampleController{}).Get)                                         //样本:获取信息
			vAgent.POST("uploaddraw", validatorFactory.Create(consts.ValidatorPrefixAgent+"ProcSampleUpload")) //样本:上传

			//防篡改相关
			vAgent.Any("getdirinfo", (&agent.TamperController{}).GetDir)                               //获取查询目录
			vAgent.POST("putdir", validatorFactory.Create(consts.ValidatorPrefixAgent+"TamperDirPut")) //上传目录信息
			vAgent.Any("getdirpolicy", (&agent.TamperController{}).GetPolicy)                          //下载防篡改策略

			//勒索相关

			//端口相关
			vAgent.POST("upserviceport", validatorFactory.Create(consts.ValidatorPrefixAgent+"PortWorkingPut")) //上传业务端口列表
			vAgent.Any("getvirtualport", (&agent.PortController{}).GetVirtual)                                  //下载虚开端口

			// 封堵相关
			vAgent.Any("getPlugging", (&agent.PlugController{}).GetDynamic)                                  // 下载动态封堵
			vAgent.POST("upplugging", validatorFactory.Create(consts.ValidatorPrefixAgent+"PlugDynamicPut")) // 上传动态封堵(调试用)
			vAgent.Any("getipblacklist", (&agent.PlugController{}).GetIpBlack)                               // 下载封堵ip黑名单
			vAgent.Any("getipwhitelist", (&agent.PlugController{}).GetIpWhite)                               // 下载封堵ip白名单(不在使用)

			//外设相关
			vAgent.Any("addperipherals", validatorFactory.Create(consts.ValidatorPrefixAgent+"UsbPut")) // 外接设备添加
			vAgent.Any("getwhiteperipherals", (&agent.UsbController{}).GetWhite)                        // 获取外设信任列表
			vAgent.Any("getblackperipherals", (&agent.UsbController{}).GetBlack)                        // 获取外设禁止列表

			//告警相关
			vAgent.POST("alertupload", validatorFactory.Create(consts.ValidatorPrefixAgent+"AlarmMasterPut")) //上传报警信息(除端口虚开外)
			vAgent.POST("upOpenPort", validatorFactory.Create(consts.ValidatorPrefixAgent+"AlarmPortPut"))    //上传端口虚开告警

			//通讯日志相关
			vAgent.POST("putsyslog", validatorFactory.Create(consts.ValidatorPrefixAgent+"SyslogPut")) //上传通讯日志
		}
	}
	return router
}
