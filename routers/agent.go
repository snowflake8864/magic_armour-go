package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
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
	vAgent := router.Group("/v1/")
	{
		vAgent.POST("auth", validatorFactory.Create(consts.ValidatorPrefix+"AuthLogin"))

		// 【需要token】中间件验证的路由
		vAgent.Use(authorization.CheckAgentToken())
		{
			//终端相关
			vAgent.Any("getconf", validatorFactory.Create(consts.ValidatorPrefix+"AgentConfigGet"))           // 获取配置信息
			vAgent.POST("getinfo", validatorFactory.Create(consts.ValidatorPrefix+"AgentInfoGet"))            //获取终端信息(调试用)
			vAgent.Any("download", validatorFactory.Create(consts.ValidatorPrefix+"AgentUpgrade"))            //升级
			vAgent.Any("uninstall", validatorFactory.Create(consts.ValidatorPrefix+"AgentUninstall"))         //卸载
			vAgent.Any("gettrustdir", validatorFactory.Create(consts.ValidatorPrefix+"TrustDirGet"))          //获取信任目录
			vAgent.POST("error/upload", validatorFactory.Create(consts.ValidatorPrefix+"ErrorFileUpload"))    //上传转储(错误) 文件
			vAgent.POST("puthardwareinfo", validatorFactory.Create(consts.ValidatorPrefix+"HardwareMonitor")) //硬件资源监控

			//任务相关
			vAgent.Any("gettask", validatorFactory.Create(consts.ValidatorPrefix+"TaskGet"))
			vAgent.Any("closetask", validatorFactory.Create(consts.ValidatorPrefix+"TaskClose"))

			//进程相关
			vAgent.POST("upload/suffix/exe", validatorFactory.Create(consts.ValidatorPrefix+"ProcessAllPut"))     //上传全部进程
			vAgent.POST("uploadproc", validatorFactory.Create(consts.ValidatorPrefix+"ProcessRunningPut"))        //上传运行中进程
			vAgent.POST("autouploadprocess", validatorFactory.Create(consts.ValidatorPrefix+"ProcessAutoPut"))    //主动上传进程
			vAgent.Any("getprocwl", validatorFactory.Create(consts.ValidatorPrefix+"ProcessWhiteGet"))            //下载进程白名单
			vAgent.Any("getprocbl", validatorFactory.Create(consts.ValidatorPrefix+"ProcessBlackGet"))            //下载进程黑名单
			vAgent.Any("module/process", validatorFactory.Create(consts.ValidatorPrefix+"ProcessAssignGet"))      //指定进程:下载信息
			vAgent.POST("upload/proc/module", validatorFactory.Create(consts.ValidatorPrefix+"ProcessAssignPut")) //指定进程:上传对应模块

			//模块相关
			vAgent.POST("upload/suffix/dll", validatorFactory.Create(consts.ValidatorPrefix+"DllAllPut"))         //上传全部模块
			vAgent.POST("upload/allproc/module", validatorFactory.Create(consts.ValidatorPrefix+"DllRunningPut")) //上传运行中模块
			vAgent.POST("auto/upload/module", validatorFactory.Create(consts.ValidatorPrefix+"DllAutoPut"))       //主动上传模块
			vAgent.Any("module/white/list", validatorFactory.Create(consts.ValidatorPrefix+"DllWhiteGet"))        //下载模块白名单
			vAgent.Any("module/black/list", validatorFactory.Create(consts.ValidatorPrefix+"DllBlackGet"))        //下载模块黑名单

			//进程或模块样本
			vAgent.Any("getdraw", validatorFactory.Create(consts.ValidatorPrefix+"ProcSampleGet"))        //进程样本:下载信息
			vAgent.POST("uploaddraw", validatorFactory.Create(consts.ValidatorPrefix+"ProcSampleUpload")) //指定进程:上传对应模块

			//防篡改相关
			vAgent.Any("getdirinfo", validatorFactory.Create(consts.ValidatorPrefix+"TamperDirGet"))      //获取查询目录
			vAgent.POST("putdir", validatorFactory.Create(consts.ValidatorPrefix+"TamperDirPut"))         //上传目录信息
			vAgent.Any("getdirpolicy", validatorFactory.Create(consts.ValidatorPrefix+"TamperPolicyGet")) //下载防篡改保护策略

			//勒索相关
			vAgent.Any("getprotect", validatorFactory.Create(consts.ValidatorPrefix+"BlackmailPolicyGet")) //下载防篡改保护策略

			//端口相关
			vAgent.POST("upserviceport", validatorFactory.Create(consts.ValidatorPrefix+"PortWorkingPut")) //上传业务端口列表
			vAgent.Any("getvirtualport", validatorFactory.Create(consts.ValidatorPrefix+"PortVirtualGet")) //下载虚开端口

			// 封堵相关
			vAgent.Any("getPlugging", validatorFactory.Create(consts.ValidatorPrefix+"PlugDynamicGet"))    // 下载动态封堵
			vAgent.Any("getipblacklist", validatorFactory.Create(consts.ValidatorPrefix+"PlugIpBlackGet")) // 下载封堵ip黑名单
			vAgent.Any("getipwhitelist", validatorFactory.Create(consts.ValidatorPrefix+"PlugIpWhiteGet")) // 下载封堵ip白名单(不在使用)
			vAgent.POST("upplugging", validatorFactory.Create(consts.ValidatorPrefix+"PlugDynamicPut"))    // 上传动态封堵(调试用)

			//外设相关
			vAgent.Any("addperipherals", validatorFactory.Create(consts.ValidatorPrefix+"UsbPut"))           // 外接设备添加
			vAgent.Any("getwhiteperipherals", validatorFactory.Create(consts.ValidatorPrefix+"UsbWhiteGet")) // 获取外设信任列表
			vAgent.Any("getblackperipherals", validatorFactory.Create(consts.ValidatorPrefix+"UsbBlackGet")) // 获取外设禁止列表

			//告警相关
			vAgent.POST("alertupload", validatorFactory.Create(consts.ValidatorPrefix+"AlarmMasterPut")) //上传报警信息(除端口虚开外)
			vAgent.POST("upOpenPort", validatorFactory.Create(consts.ValidatorPrefix+"AlarmPortPut"))    //端口虚开告警上传

			//通讯日志相关
			vAgent.POST("putsyslog", validatorFactory.Create(consts.ValidatorPrefix+"SyslogPut")) //端口虚开告警上传
		}
	}
	return router
}
