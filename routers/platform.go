package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/http/middleware/cors"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"goskeleton/app/utils/gin_release"
	"net/http"
)

// 该路由主要设置平台路由

func InitPlatformRouter() *gin.Engine {
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
		context.String(http.StatusOK, "Platform接口: "+consts.VersionNumber)
	})

	//  创建一个接口路由组
	vPlatform := router.Group("/platform/")
	{
		vPlatform.Use(authorization.RequestParams).Use(authorization.SavePlatformOperate())
		{
			//登录认证
			auth := vPlatform.Group("auth/")
			{
				auth.POST("checkUser", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthCheckUser"))
				auth.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthGetSliderCode"))
				auth.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthSendSms"))
				auth.POST("normalLogin", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthNormalLogin"))
				auth.POST("logout", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthLogout"))
				auth.POST("hsLogin", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthHsLogin"))
			}

			// 【需要token】中间件验证的路由
			vPlatform.Use(authorization.CheckPlatformPToken())
			{
				//首页
				home := vPlatform.Group("home/")
				{
					home.POST("getAgentStatistic", validatorFactory.Create(consts.ValidatorPrefixPlatform+"HomeGetAgentStatistic"))
					home.POST("getAlarmStatistic", validatorFactory.Create(consts.ValidatorPrefixPlatform+"HomeGetAlarmStatistic"))
				}
				//菜单
				menu := vPlatform.Group("menu/")
				{
					menu.POST("tree", validatorFactory.Create(consts.ValidatorPrefixPlatform+"MenuTree"))
					menu.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"MenuAdd"))
					menu.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"MenuEdit"))
					menu.POST("sort", validatorFactory.Create(consts.ValidatorPrefixPlatform+"MenuSort"))
					menu.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"MenuDelete"))
				}
				//角色
				role := vPlatform.Group("role/")
				{
					manage := role.Group("manage/")
					{
						manage.POST("all", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RoleManageAll"))
						manage.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RoleManageAdd"))
						manage.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RoleManageEdit"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RoleManageDelete"))
					}
					power := role.Group("power/")
					{
						power.POST("loginMenuTree", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RolePowerLoginMenuTree"))       //用户登录后的菜单树
						power.POST("editableMenuTree", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RolePowerEditableMenuTree")) //可编辑的菜单树
						power.POST("getSelectedMenus", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RolePowerGetSelectedMenus")) //获取已选中的菜单集合
						power.POST("setSelectedMenus", validatorFactory.Create(consts.ValidatorPrefixPlatform+"RolePowerSetSelectedMenus")) //保存已选中的菜单集合
					}
				} //角色
				//组织架构
				department := vPlatform.Group("department/")
				{
					department.POST("tree", validatorFactory.Create(consts.ValidatorPrefixPlatform+"DepartmentTree"))
					department.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"DepartmentAdd"))
					department.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"DepartmentEdit"))
					department.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"DepartmentDelete"))
				}

				//终端
				agent := vPlatform.Group("agent/")
				{
					public := agent.Group("public/")
					{
						public.POST("onlineStatus", (&platform.AgentController{}).OnlineStatus)
						public.POST("safeStatus", (&platform.AgentController{}).SafeStatus)
						public.POST("countOnline", (&platform.AgentController{}).CountOnline)
						public.POST("typeList", (&platform.AgentController{}).TypeList)
					}
					manage := agent.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManagePage"))
						//manage.POST("import", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageImport"))
						//manage.POST("export", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageExport"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageDelete"))
						manage.POST("uninstall", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageUninstall"))
						manage.POST("upgrade", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageUpgrade"))
						manage.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageBatchDelete"))
						manage.POST("batchUninstall", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageBatchUninstall"))
						manage.POST("batchUpgrade", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageBatchUpgrade"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManageBatchSetConfig"))
					}
					basic := agent.Group("basic/")
					{
						basic.POST("get", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentBasicGet"))
						basic.POST("save", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentBasicSave"))
					}
					config := agent.Group("config/")
					{
						config.POST("get", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentConfigGet"))
						config.POST("save", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentConfigSave"))
					}
					monitor := agent.Group("monitor/")
					{
						monitor.POST("list", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentMonitorList"))
						monitor.POST("info", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentMonitorInfo"))
					}
					errorV := agent.Group("error")
					{
						errorV.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentErrorPage"))
					}
				} //终端

				//进程
				process := vPlatform.Group("process/")
				{
					public := process.Group("public/")
					{
						public.POST("searchStatus", (&platform.ProcessController{}).SearchStatus)
						public.POST("changeStatus", (&platform.ProcessController{}).ChangeStatus)
					}
					manage := process.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchSetConfig"))
						manage.POST("batchSetWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchSetWhite"))
						manage.POST("batchRefresh", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchRefresh"))
						manage.POST("batchCollect", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchCollect"))
						manage.POST("batchCollect", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchCollect"))
						manage.POST("batchThreatenWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessManageBatchThreatenWhite"))
					}
					agentV2 := process.Group("agent/")
					{
						agentV2.POST("pageAll", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentPageAll"))
						agentV2.POST("pageRunning", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentPageRunning"))
						agentV2.POST("setStatus", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentSetStatus"))
						agentV2.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentCancel"))
						agentV2.POST("clearAll", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentClearAll"))
						agentV2.POST("clearRunning", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentClearRunning"))
						agentV2.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentRefresh"))
						agentV2.POST("collect", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentCollect"))
						agentV2.POST("extractRunningDll", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentExtractRunningDll"))
						agentV2.POST("showRunningDll", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentShowRunningDll"))
						agentV2.POST("allWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentAllWhite"))
						agentV2.POST("repeatPage", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessAgentRepeatPage"))

						agentV2.POST("extractSample", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))
						agentV2.POST("setGlobal", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalAdd"))
						agentV2.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))
						manage.POST("batchThreatenWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))
						manage.POST("batchThreatenBlack", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))

					}

					global := process.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalAdd"))
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalEdit"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalDelete"))
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessGlobalClear"))
					}

					search := process.Group("search/")
					{
						search.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcessSearchPage"))
						search.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))
						search.POST("batchGlobalBlack", validatorFactory.Create(consts.ValidatorPrefixPlatform+""))
					}
				} //进程

				//进程和模块样本
				procSample := vPlatform.Group("procSample/")
				{
					procSample.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcSamplePage"))
					procSample.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcSampleDelete"))
					procSample.POST("clear", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcSampleClear"))
					procSample.POST("download", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProcSampleDownload"))
				} //进程和模块样本

				//防篡改
				tamper := vPlatform.Group("tamper/")
				{
					public := tamper.Group("public/")
					{
						public.POST("FilePower", (&platform.TamperController{}).FilePower)
						public.POST("FolderPower", (&platform.TamperController{}).FolderPower)
					}
					manage := tamper.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperManageBatchSetConfig"))
					}

					policy := tamper.Group("policy/")
					{
						policy.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyPage"))
						policy.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicySetSwitch"))
						policy.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyAddFolder"))
						policy.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyAddFile"))
						policy.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyEditFolder"))
						policy.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyEditFile"))
						policy.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperPolicyDelete"))
					}

					global := tamper.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalSetSwitch"))
						global.POST("addFolder", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalAddFolder"))
						global.POST("addFile", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalAddFile"))
						global.POST("editFolder", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalEditFolder"))
						global.POST("editFile", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalEditFile"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"TamperGlobalDelete"))
					}
				} //防篡改
			}
		}
	}
	return router
}
