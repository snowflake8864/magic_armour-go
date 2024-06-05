package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/authorization/middleware_auth_org"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/cur_app"
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
		context.String(http.StatusOK, "Organization接口: "+variable.VersionNumber)
	})

	//  创建一个接口路由组
	vOrg := router.Group("/v5/organization/")
	{
		vOrg.GET("customConfig", func(context *gin.Context) {
			context.JSON(http.StatusOK, cur_app.GetOrganizationConfig())
		})

		vOrg.Use(authorization.RequestParams).Use(middleware_auth_org.SaveOperate())
		{
			//登录认证
			auth := vOrg.Group("auth/")
			{
				auth.POST("checkUser", validatorFactory.Create(consts.ValidatorPrefixOrg+"AuthCheckUser"))
				auth.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixOrg+"AuthGetSliderCode"))
				auth.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixOrg+"AuthSendSms"))
				auth.POST("normalLogin", validatorFactory.Create(consts.ValidatorPrefixOrg+"AuthNormalLogin"))
				auth.POST("logout", (&organization.AuthController{}).Logout)
				auth.POST("hsLogin", validatorFactory.Create(consts.ValidatorPrefixOrg+"AuthHsLogin"))
			}

			registry := vOrg.Group("registry/").Use(middleware_auth_org.ApiIsEnabled())
			{
				registry.POST("types", (&organization.RegistryController{}).Types)
				registry.POST("checkOrgName", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistryCheckOrgName"))
				registry.POST("checkUserName", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistryCheckUserName"))
				registry.POST("checkUserPhone", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistryCheckUserPhone"))
				registry.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistrySendSms"))
				registry.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistryGetSliderCode"))
				registry.POST("signup", validatorFactory.Create(consts.ValidatorPrefixOrg+"RegistrySignup"))
			}
			findPwd := vOrg.Group("findPw/").Use(middleware_auth_org.ApiIsEnabled())
			{
				findPwd.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixOrg+"FindPwGetSliderCode"))
				findPwd.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixOrg+"FindPwSendSms"))
				findPwd.POST("checkSms", validatorFactory.Create(consts.ValidatorPrefixOrg+"FindPwCheckSms"))
				findPwd.POST("modifyPw", validatorFactory.Create(consts.ValidatorPrefixOrg+"FindPwModifyPw"))
			}

			// 【需要token】中间件验证的路由
			vOrg.Use(middleware_auth_org.CheckToken())
			{
				//组织架构
				department := vOrg.Group("department/")
				{
					department.POST("tree", (&organization.DepartmentController{}).Tree)
					department.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentAdd"))
					department.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentEdit"))
					department.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentDelete"))

				}

				//用户列表
				user := vOrg.Group("user/")
				{
					user.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"UserPage"))
					user.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"UserAdd"))
					user.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"UserEdit"))
					user.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"UserDelete"))
					user.POST("modifyPw", validatorFactory.Create(consts.ValidatorPrefixOrg+"UserModifyPw"))
				}

				//首页
				home := vOrg.Group("home/")
				{
					home.POST("agentStat", (&organization.HomeController{}).AgentStat)
					home.POST("alarmAgentTops", (&organization.HomeController{}).AlarmAgentTops)
					home.POST("alarmTodayStat", (&organization.HomeController{}).AlarmTodayStat)
					home.POST("alarmDailyStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"HomeAlarmDailyStat"))
					home.POST("grantUsedStat", (&organization.HomeController{}).GrantUsedStat)
				}

				//终端
				agent := vOrg.Group("agent/")
				{
					public := agent.Group("public/")
					{
						public.POST("onlineStatus", (&organization.AgentController{}).OnlineStatus)
						public.POST("safeStatus", (&organization.AgentController{}).SafeStatus)
						public.POST("types", (&organization.AgentController{}).TypeList)
					}
					manage := agent.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManagePage"))
						//manage.POST("import", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageImport"))
						//manage.POST("export", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageExport"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageDelete"))
						manage.POST("uninstall", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageUninstall"))
						manage.POST("upgrade", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageUpgrade"))
						manage.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageBatchDelete"))
						manage.POST("batchUninstall", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageBatchUninstall"))
						manage.POST("batchUpgrade", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageBatchUpgrade"))
						manage.POST("batchModify", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentManageBatchModify"))
					}
					basic := agent.Group("basic/")
					{
						basic.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentBasicGet"))
						basic.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentBasicSave"))
					}
					config := agent.Group("config/")
					{
						config.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentConfigGet"))
						config.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentConfigSave"))
					}
					monitor := agent.Group("monitor/")
					{
						monitor.POST("list", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentMonitorList"))
						monitor.POST("info", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentMonitorInfo"))
					}

				} //终端

				//进程
				process := vOrg.Group("process/")
				{
					public := process.Group("public/")
					{
						public.POST("searchStatus", (&organization.ProcessController{}).SearchStatus)
						public.POST("changeStatus", (&organization.ProcessController{}).ChangeStatus)
						public.POST("searchWhile", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentSearchWhile"))
					}
					manage := process.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManagePage"))
						manage.POST("agentProcStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageAgentProcStat"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageBatchSetConfig"))
						manage.POST("batchSetWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageBatchSetWhite"))
						manage.POST("batchRefresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageBatchRefresh"))
						manage.POST("batchCollect", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessManageBatchCollect"))
					}

					agentV2 := process.Group("agent/")
					{

						agentV2.POST("pageAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentPageAll"))
						agentV2.POST("pageRunning", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentPageRunning"))
						agentV2.POST("setStatus", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentSetStatus"))
						agentV2.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentCancel"))
						agentV2.POST("clearAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentClearAll"))
						agentV2.POST("clearRunning", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentClearRunning"))
						agentV2.POST("batchWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentBatchWhite"))
						agentV2.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentRefresh"))
						agentV2.POST("collect", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentCollect"))

						agentV2.POST("runningDllExtract", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentRunningDllExtract"))
						agentV2.POST("runningDllShow", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentRunningDllShow"))

						agentV2.POST("repeatStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentRepeatStat"))
						agentV2.POST("repeatPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentRepeatPage"))

						agentV2.POST("sampleExtract", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentSampleExtract"))

						agentV2.POST("setGlobal", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentSetGlobal"))
						agentV2.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentBatchGlobalWhite"))
					}
					global := process.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalAdd"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalSetSwitch"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalDelete"))
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalClear"))
					}

					/*
						//@todo 目的为了解决全局防篡改、防篡改、勒索防护的信任进程问题，新增获取进程白名单改成查询方式，不在使用
						search := process.Group("search/")
						{
							search.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessSearchPage"))
							search.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessSearchBatchGlobalWhite"))
							search.POST("batchGlobalBlack", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessSearchBatchGlobalBlack"))
						}
					*/
				} //进程

				//模块
				dll := vOrg.Group("dll/")
				{
					manage := dll.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManagePage"))
						manage.POST("agentProcStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageAgentProcStat"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageBatchSetConfig"))
						manage.POST("batchSetWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageBatchSetWhite"))
						manage.POST("batchRefresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageBatchRefresh"))
						manage.POST("batchCollect", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllManageBatchCollect"))
					}
					agentV2 := dll.Group("agent/")
					{
						agentV2.POST("pageAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentPageAll"))
						agentV2.POST("pageRunning", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentPageRunning"))
						agentV2.POST("setStatus", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentSetStatus"))
						agentV2.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentCancel"))
						agentV2.POST("clearAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentClearAll"))
						agentV2.POST("clearRunning", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentClearRunning"))
						agentV2.POST("batchWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentBatchWhite"))
						agentV2.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentRefresh"))
						agentV2.POST("collect", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentCollect"))

						agentV2.POST("repeatStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentRepeatStat"))
						agentV2.POST("repeatPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentRepeatPage"))

						agentV2.POST("setGlobal", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentSetGlobal"))
						agentV2.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentBatchGlobalWhite"))
					}

					global := dll.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalAdd"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalSetSwitch"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalDelete"))
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalClear"))
					}
				} //模块

				//进程和模块样本
				procSample := vOrg.Group("procSample/")
				{
					procSample.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcSamplePage"))
					procSample.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcSampleDelete"))
					procSample.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcSampleClear"))
					procSample.POST("download", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcSampleDownload"))
				} //进程和模块样本

				//威胁情报
				threaten := vOrg.Group("threaten/")
				{
					threaten.POST("getLevels", validatorFactory.Create(consts.ValidatorPrefixOrg+"ThreatenGetLevels"))
					threaten.POST("getAttr", validatorFactory.Create(consts.ValidatorPrefixOrg+"ThreatenGetAttr"))
				}

				//防篡改
				tamper := vOrg.Group("tamper/")
				{
					public := tamper.Group("public/")
					{
						public.POST("FilePower", (&organization.TamperController{}).FilePower)
						public.POST("FolderPower", (&organization.TamperController{}).FolderPower)
					}
					manage := tamper.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperManageBatchSetConfig"))
					}

					policy := tamper.Group("policy/")
					{
						policy.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyPage"))
						policy.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicySetSwitch"))
						policy.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyAddFolder"))
						policy.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyAddFile"))
						policy.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyEditFolder"))
						policy.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyEditFile"))
						policy.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperPolicyDelete"))
					}

					global := tamper.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalSetSwitch"))
						global.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalAddFolder"))
						global.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalAddFile"))
						global.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalEditFolder"))
						global.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalEditFile"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"TamperGlobalDelete"))
					}
				} //防篡改
			}

		}
	}
	return router
}
