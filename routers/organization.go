package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/gin_release"
	"net/http"
)

// 该路由主要设置平台路由

func InitOrganizationRouter() *gin.Engine {
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
		context.String(http.StatusOK, "Organization接口: "+consts.VersionNumber)
	})

	//  创建一个接口路由组
	vPlatform := router.Group("/organization/")
	{
		vPlatform.Use(authorization.RequestParams).Use(authorization.SavePlatformOperate())
		{
			//登录认证
			auth := vPlatform.Group("auth/")
			{
				auth.POST("checkUser", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthCheckUser"))
				auth.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthGetSliderCode"))
				auth.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthSendSms"))
				auth.POST("normalLogin", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthNormalLogin"))
				auth.POST("logout", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthLogout"))
				auth.POST("hsLogin", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AuthHsLogin"))
			}

			// 【需要token】中间件验证的路由
			vPlatform.Use(authorization.CheckPlatformPToken())
			{
				//首页
				home := vPlatform.Group("home/")
				{
					home.POST("getAgentStatistic", validatorFactory.Create(consts.ValidatorPrefixOrganization+"HomeGetAgentStatistic"))
					home.POST("getAlarmStatistic", validatorFactory.Create(consts.ValidatorPrefixOrganization+"HomeGetAlarmStatistic"))
				}

				//终端
				agent := vPlatform.Group("agent/")
				{
					public := agent.Group("public/")
					{
						public.POST("onlineStatus", (&platform.AgentController{}).OnlineStatus)
						public.POST("safeStatus", (&platform.AgentController{}).SafeStatus)
						public.POST("typeList", (&platform.AgentController{}).TypeList)
					}
					manage := agent.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManagePage"))
						//manage.POST("import", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageImport"))
						//manage.POST("export", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageExport"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageDelete"))
						manage.POST("uninstall", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageUninstall"))
						manage.POST("upgrade", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageUpgrade"))
						manage.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageBatchDelete"))
						manage.POST("batchUninstall", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageBatchUninstall"))
						manage.POST("batchUpgrade", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageBatchUpgrade"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentManageBatchSetConfig"))
					}
					basic := agent.Group("basic/")
					{
						basic.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentBasicGet"))
						basic.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentBasicSave"))
					}
					config := agent.Group("config/")
					{
						config.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentConfigGet"))
						config.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentConfigSave"))
					}
					monitor := agent.Group("monitor/")
					{
						monitor.POST("list", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentMonitorList"))
						monitor.POST("info", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentMonitorInfo"))
					}
					errorV := agent.Group("error")
					{
						errorV.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"AgentErrorPage"))
					}
				} //终端

				//进程
				process := vPlatform.Group("process/")
				{
					public := process.Group("public/")
					{
						public.POST("searchStatus", (&organization.ProcessController{}).SearchStatus)
						public.POST("changeStatus", (&organization.ProcessController{}).ChangeStatus)
					}
					manage := process.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageBatchSetConfig"))
						manage.POST("batchSetWhite", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageBatchSetWhite"))
						manage.POST("batchRefresh", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageBatchRefresh"))
						manage.POST("batchCollect", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessManageBatchCollect"))
					}
					agentV2 := process.Group("agent/")
					{
						agentV2.POST("pageAll", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentPageAll"))
						agentV2.POST("pageRunning", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentPageRunning"))
						agentV2.POST("setStatus", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentSetStatus"))
						agentV2.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentCancel"))
						agentV2.POST("clearAll", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentClearAll"))
						agentV2.POST("clearRunning", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentClearRunning"))
						agentV2.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentRefresh"))
						agentV2.POST("collect", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentCollect"))
						agentV2.POST("extractRunningDll", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentExtractRunningDll"))
						agentV2.POST("showRunningDll", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentShowRunningDll"))
						agentV2.POST("batchWhite", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentBatchWhite"))
						agentV2.POST("repeatHash", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentRepeatPage"))

						agentV2.POST("extractSample", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentExtractSample"))
						agentV2.POST("setGlobal", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentSetGlobal"))
						agentV2.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessAgentBatchGlobalWhite"))
					}

					global := process.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessGlobalAdd"))
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessGlobalEdit"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessGlobalDelete"))
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessGlobalClear"))
					}

					search := process.Group("search/")
					{
						search.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessSearchPage"))
						search.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessSearchBatchGlobalWhite"))
						search.POST("batchGlobalBlack", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcessSearchBatchGlobalBlack"))
					}
				} //进程

				//进程和模块样本
				procSample := vPlatform.Group("procSample/")
				{
					procSample.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcSamplePage"))
					procSample.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcSampleDelete"))
					procSample.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcSampleClear"))
					procSample.POST("download", validatorFactory.Create(consts.ValidatorPrefixOrganization+"ProcSampleDownload"))
				} //进程和模块样本

				//防篡改
				tamper := vPlatform.Group("tamper/")
				{
					public := tamper.Group("public/")
					{
						public.POST("FilePower", (&organization.TamperController{}).FilePower)
						public.POST("FolderPower", (&organization.TamperController{}).FolderPower)
					}
					manage := tamper.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperManageBatchSetConfig"))
					}

					policy := tamper.Group("policy/")
					{
						policy.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyPage"))
						policy.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicySetSwitch"))
						policy.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyAddFolder"))
						policy.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyAddFile"))
						policy.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyEditFolder"))
						policy.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyEditFile"))
						policy.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperPolicyDelete"))
					}

					global := tamper.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalSetSwitch"))
						global.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalAddFolder"))
						global.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalAddFile"))
						global.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalEditFolder"))
						global.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalEditFile"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrganization+"TamperGlobalDelete"))
					}
				} //防篡改
			}
		}
	}
	return router
}
