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
	"goskeleton/app/http/middleware/rate_limit"
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

	//限流
	if isOpen := variable.ConfigYml.GetInt("RateLimiting.IsOpen"); isOpen == 1 {
		var rate float64 = 1.0
		var limit int64 = 100
		if val := variable.ConfigYml.GetFloat64("RateLimiting.Common.Rate"); val > 0 {
			rate = val
		}
		if val := variable.ConfigYml.GetInt64("RateLimiting.Common.Limit"); val > 0 {
			limit = val
		}
		router.Use(rate_limit.LimitHandler(rate, limit))
	}

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
				//账户
				account := vOrg.Group("account/")
				{
					account.POST("getInfo", (&organization.AccountController{}).GetUser)
					account.POST("saveInfo", validatorFactory.Create(consts.ValidatorPrefixOrg+"AccountSaveInfo"))
					account.POST("getMenuTree", (&organization.AccountController{}).GetMenuTree)
					account.POST("modifyPw", validatorFactory.Create(consts.ValidatorPrefixOrg+"AccountModifyPw"))
					account.POST("modifyPhone", validatorFactory.Create(consts.ValidatorPrefixOrg+"AccountModifyPhone"))
					account.POST("modifyEmail", validatorFactory.Create(consts.ValidatorPrefixOrg+"AccountModifyEmail"))
				}

				//组织架构
				department := vOrg.Group("department/")
				{
					department.POST("tree", (&organization.DepartmentController{}).Tree)
					department.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentAdd"))
					department.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentEdit"))
					department.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"DepartmentDelete"))

				}

				//角色
				role := vOrg.Group("role/")
				{
					role.POST("getAll", (&organization.RoleController{}).AllForManage)
					role.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"RoleAdd"))
					role.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"RoleEdit"))
					role.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"RoleDelete"))
					role.POST("getMenuTree", validatorFactory.Create(consts.ValidatorPrefixOrg+"RoleGetMenuTree"))
					role.POST("saveMenus", validatorFactory.Create(consts.ValidatorPrefixOrg+"RoleSaveMenus"))
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
						public.POST("grantStatus", (&organization.AgentController{}).GrantStatus)
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
						basic.POST("resetUid", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentBasicResetUid"))
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
					onlineLog := agent.Group("onlineLog/")
					{
						onlineLog.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentOnlineLogPage"))
					}

				} //终端

				//信任目录
				trustDir := vOrg.Group("trustDir/")
				{
					global := trustDir.Group("global/")
					{
						global.POST("getTypes", validatorFactory.Create(consts.ValidatorPrefixOrg+"TrustDirGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"TrustDirGlobalAdd"))
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"TrustDirGlobalEdit"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"TrustDirGlobalSetSwitch"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"TrustDirGlobalDelete"))
					}
				} //信任目录

				//文件目录
				fileDir := vOrg.Group("fileDir/")
				{
					fileDir.POST("getTypes", (&organization.FileDirController{}).GetTypes)
					fileDir.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"FileDirPage"))
					fileDir.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"FileDirRefresh"))
					fileDir.POST("hasUpload", validatorFactory.Create(consts.ValidatorPrefixOrg+"FileDirHasUpload"))
				} //文件目录

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
						agentV2.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessAgentBatchDelete"))

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
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"ProcessGlobalClear")) //仅管理员
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
						agentV2.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentBatchDelete"))

						agentV2.POST("repeatStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentRepeatStat"))
						agentV2.POST("repeatPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentRepeatPage"))

						agentV2.POST("sampleExtract", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentSampleExtract"))

						agentV2.POST("setGlobal", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentSetGlobal"))
						agentV2.POST("batchGlobalWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllAgentBatchGlobalWhite"))
					}

					global := dll.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalPage"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalAdd"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalSetSwitch"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalDelete"))
						global.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"DllGlobalClear")) //仅管理员
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
						public.POST("filePower", (&organization.TamperController{}).FilePower)
						public.POST("folderPower", (&organization.TamperController{}).FolderPower)
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

				//勒索防护
				blackmail := vOrg.Group("blackmail/")
				{
					public := blackmail.Group("public/")
					{
						public.POST("types", (&organization.BlackmailController{}).Types)
					}
					group := blackmail.Group("group/")
					{
						group.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGroupPage"))
						group.POST("getAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGroupGetAll"))
						group.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGroupAdd"))
						group.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGroupEdit"))
						group.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGroupDelete"))
					}
					manage := blackmail.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailManageBatchSetConfig"))
					}
					policy := blackmail.Group("policy/")
					{
						policy.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailPolicyPage"))
						policy.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailPolicySetSwitch"))
						policy.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailPolicyAdd"))
						policy.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailPolicyEdit"))
						policy.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailPolicyDelete"))
					}
					global := blackmail.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGlobalSetSwitch"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGlobalAdd"))
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGlobalEdit"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"BlackmailGlobalDelete"))
					}
				} //勒索防护

				//端口虚开
				port := vOrg.Group("port/")
				{
					public := port.Group("public/")
					{
						public.POST("protocols", (&organization.PortController{}).Protocols)
						public.POST("ipTypes", (&organization.PortController{}).IpTypes)
						public.POST("alarmLevels", (&organization.PortController{}).AlarmLevels)
						public.POST("forwardTypes", (&organization.PortController{}).ForwardTypes)
					}
					working := port.Group("working/")
					{
						working.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortWorkingPage"))
						working.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortWorkingRefresh"))
					}
					manage := port.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortManageBatchSetConfig"))
					}
					policy := port.Group("policy/")
					{
						policy.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortPolicyPage"))
						policy.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortPolicySetSwitch"))
						policy.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortPolicyAdd"))
						policy.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortPolicyEdit"))
						policy.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortPolicyDelete"))
					}
					global := port.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortGlobalSetSwitch"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortGlobalAdd"))
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortGlobalEdit"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PortGlobalDelete"))
					}
				} //端口虚开

				//动态封堵
				plug := vOrg.Group("plug/")
				{
					public := plug.Group("public/")
					{
						public.POST("plugTypes", (&organization.PlugController{}).PlugTypes)
						public.POST("effectModes", (&organization.PlugController{}).EffectModes)
					}

					manage := plug.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugManageBatchSetConfig"))
					}

					//策略
					policy := plug.Group("policy/")
					{
						configPort := policy.Group("configPort/")
						{
							configPort.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyConfigPortGet"))
							configPort.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyConfigPortSet"))
						}
						configLogin := policy.Group("configLogin/")
						{
							configLogin.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyConfigLoginGet"))
							configLogin.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyConfigLoginSet"))
						}

						ipWhite := policy.Group("ipWhite/")
						{
							ipWhite.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhitePage"))
							ipWhite.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhiteSetSwitch"))
							ipWhite.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhiteAdd"))
							ipWhite.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhiteEdit"))
							ipWhite.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhiteDelete"))
							ipWhite.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpWhiteClear"))
						}
						ipBlack := policy.Group("ipBlack/")
						{
							ipBlack.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackPage"))
							ipBlack.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackSetSwitch"))
							ipBlack.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackAdd"))
							ipBlack.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackEdit"))
							ipBlack.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackDelete"))
							ipBlack.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyIpBlackClear"))
						}
						dynamic := policy.Group("dynamic/")
						{
							dynamic.POST("enabledPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicEnabledPage"))
							dynamic.POST("disabledPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicDisabledPage"))
							dynamic.POST("disable", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicDisable"))
							dynamic.POST("disableBatch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicDisableBatch"))
							dynamic.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicDelete"))
							dynamic.POST("deleteBatch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugPolicyDynamicDeleteBatch"))
						}
					}

					//全局
					global := plug.Group("global/")
					{
						configPort := global.Group("configPort/")
						{
							configPort.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalConfigPortGet"))
							configPort.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalConfigPortSet"))
						}
						configLogin := global.Group("configLogin/")
						{
							configLogin.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalConfigLoginGet"))
							configLogin.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalConfigLoginSet"))
						}

						ipWhite := global.Group("ipWhite/")
						{
							ipWhite.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhitePage"))
							ipWhite.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhiteSetSwitch"))
							ipWhite.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhiteAdd"))
							ipWhite.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhiteEdit"))
							ipWhite.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhiteDelete"))
							ipWhite.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpWhiteClear")) //仅管理员
						}
						ipBlack := global.Group("ipBlack/")
						{
							ipBlack.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackPage"))
							ipBlack.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackSetSwitch"))
							ipBlack.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackAdd"))
							ipBlack.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackEdit"))
							ipBlack.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackDelete"))
							ipBlack.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalIpBlackClear")) //仅管理员
						}
						dynamic := global.Group("dynamic/")
						{

							dynamic.POST("enabledPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicEnabledPage"))
							dynamic.POST("disabledPage", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicDisabledPage"))
							dynamic.POST("disable", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicDisable"))
							dynamic.POST("disableBatch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicDisableBatch"))
							dynamic.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicDelete"))
							dynamic.POST("deleteBatch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PlugGlobalDynamicDeleteBatch"))
						}
					}
				} //动态封堵

				//外接设备
				usb := vOrg.Group("usb/")
				{
					public := usb.Group("public/")
					{
						public.POST("status", (&organization.UsbController{}).EditStatus)
					}
					manage := usb.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbManageSetSwitch"))
						manage.POST("setMode", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbManageSetMode"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbManageBatchSetConfig"))
					}
					policy := usb.Group("policy/")
					{
						policy.POST("pageRunning", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyPageRunning"))
						policy.POST("pageHistory", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyPageHistory"))
						policy.POST("pageLog", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyPageLog"))
						policy.POST("setWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicySetWhite"))
						policy.POST("setBlack", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicySetBlack"))
						policy.POST("setAlias", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicySetAlias"))
						policy.POST("refresh", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyRefresh"))
						policy.POST("hasUpload", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyHasUpload"))
						policy.POST("addGlobal", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbPolicyAddGlobal"))
					}
					global := usb.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbGlobalPage"))
						global.POST("setWhite", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbGlobalSetWhite"))
						global.POST("setBlack", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbGlobalSetBlack"))
						global.POST("setAlias", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbGlobalSetAlias"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"UsbGlobalDelete"))
					}
				} //外接设备

				//外联管控
				outreach := vOrg.Group("outreach/")
				{
					manage := outreach.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachManageBatchSetConfig"))
					}
					global := outreach.Group("global/")
					{
						global.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachGlobalPage"))
						global.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachGlobalSetSwitch"))
						global.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachGlobalAdd"))
						global.POST("fastAdd", (&organization.OutreachController{}).FastAddForGlobal)
						global.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachGlobalEdit"))
						global.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachGlobalDelete"))
					}
					detect := outreach.Group("detect/")
					{
						detect.POST("list", (&organization.OutreachController{}).ListForDetect)
						detect.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachDetectAdd"))
						detect.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachDetectSetSwitch"))
						detect.POST("edit", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachDetectEdit"))
						detect.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"OutreachDetectDelete"))
					}
				} //外联管控

				//口令破解
				loginSolve := vOrg.Group("loginSolve/")
				{
					manage := loginSolve.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolveManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolveManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolveManageBatchSetConfig"))
					}
					policyConfig := loginSolve.Group("policyConfig/")
					{
						policyConfig.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolvePolicyConfigGet"))
						policyConfig.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolvePolicyConfigSet"))
					}
					globalConfig := loginSolve.Group("globalConfig/")
					{
						globalConfig.POST("get", (&organization.LoginSolveController{}).ConfigGetForGlobal)
						globalConfig.POST("set", validatorFactory.Create(consts.ValidatorPrefixOrg+"LoginSolveGlobalConfigSet"))

					}
				} //口令破解

				//密码跳变
				pwJump := vOrg.Group("pwJump/")
				{
					public := pwJump.Group("public/")
					{
						public.POST("combines", (&organization.PwJumpController{}).Combines)
						public.POST("cycles", (&organization.PwJumpController{}).Cycles)
					}
					manage := pwJump.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpManageBatchSetConfig"))
						manage.POST("showPw", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpManageShowPw"))
					}
					agentConfig := pwJump.Group("agentConfig/")
					{
						agentConfig.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpAgentConfigGet"))
						agentConfig.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpAgentConfigSave"))
						agentConfig.POST("execute", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpAgentConfigExecute"))
					}
					globalConfig := pwJump.Group("globalConfig/")
					{
						globalConfig.POST("get", (&organization.PwJumpController{}).GetForGlobalConfig)
						globalConfig.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpGlobalConfigSave"))
						globalConfig.POST("execute", validatorFactory.Create(consts.ValidatorPrefixOrg+"PwJumpGlobalConfigExecute"))
					}
				} //密码跳变

				//ip跳变
				ipJump := vOrg.Group("ipJump/")
				{
					public := ipJump.Group("public/")
					{
						public.POST("cycles", (&organization.IpJumpController{}).Cycles)
					}
					manage := ipJump.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpManagePage"))
						manage.POST("setSwitch", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpManageSetSwitch"))
						manage.POST("batchSetConfig", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpManageBatchSetConfig"))
					}
					agentConfig := ipJump.Group("agentConfig/")
					{
						agentConfig.POST("get", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpAgentConfigGet"))
						agentConfig.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpAgentConfigSave"))
						agentConfig.POST("execute", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpAgentConfigExecute"))
					}
					globalConfig := ipJump.Group("globalConfig/")
					{
						globalConfig.POST("get", (&organization.IpJumpController{}).GetForGlobalConfig)
						globalConfig.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpGlobalConfigSave"))
						globalConfig.POST("execute", validatorFactory.Create(consts.ValidatorPrefixOrg+"IpJumpGlobalConfigExecute"))
					}
				} //ip跳变

				//备份回滚
				backup := vOrg.Group("backup")
				{
					manage := backup.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BackupManagePage"))
					}
					history := backup.Group("history/")
					{
						history.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"BackupHistoryPage"))
						history.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"BackupHistoryAdd"))
						history.POST("rollback", validatorFactory.Create(consts.ValidatorPrefixOrg+"BackupHistoryRollback"))
					}
				} //备份回滚

				//告警
				alarm := vOrg.Group("alarm/")
				{
					public := alarm.Group("public/")
					{
						public.POST("handleStatus", (&organization.AlarmController{}).HandleStatus)
						public.POST("levels", (&organization.AlarmController{}).Levels)
						public.POST("typeProcess", (&organization.AlarmController{}).TypeProcess)
						public.POST("typeDll", (&organization.AlarmController{}).TypeDll)
						public.POST("typeBlackmail", (&organization.AlarmController{}).TypeBlackmail)
						public.POST("typeTamper", (&organization.AlarmController{}).TypeTamper)
						public.POST("typeUsb", (&organization.AlarmController{}).TypeUsb)
						public.POST("typeSelf", (&organization.AlarmController{}).TypeSelf)
					}
					manage := alarm.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmManagePage"))
						manage.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmManageClear"))
						manage.POST("batchClear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmManageBatchClear"))
						manage.POST("agentStat", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmManageAgentStat"))
					}

					process1 := alarm.Group("process/")
					{
						process1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmProcessPage"))
						process1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmProcessHandle"))
						process1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmProcessHandleAll"))
						process1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmProcessClear"))
					}
					dll1 := alarm.Group("dll/")
					{
						dll1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmDllPage"))
						dll1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmDllHandle"))
						dll1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmDllHandleAll"))
						dll1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmDllClear"))
					}
					tamper1 := alarm.Group("tamper/")
					{
						tamper1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmTamperPage"))
						tamper1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmTamperHandle"))
						tamper1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmTamperHandleAll"))
						tamper1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmTamperClear"))
					}
					blackmail1 := alarm.Group("blackmail/")
					{
						blackmail1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmBlackmailPage"))
						blackmail1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmBlackmailHandle"))
						blackmail1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmBlackmailHandleAll"))
						blackmail1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmBlackmailClear"))
					}
					port1 := alarm.Group("port/")
					{
						port1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPortPage"))
						port1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPortHandle"))
						port1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPortHandleAll"))
						port1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPortClear"))
					}
					plug1 := alarm.Group("plug/")
					{
						plug1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPlugPage"))
						plug1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPlugHandle"))
						plug1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPlugHandleAll"))
						plug1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmPlugClear"))
					}
					usb1 := alarm.Group("usb/")
					{
						usb1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmUsbPage"))
						usb1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmUsbHandle"))
						usb1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmUsbHandleAll"))
						usb1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmUsbClear"))
					}
					outreach1 := alarm.Group("outreach/")
					{
						outreach1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmOutreachPage"))
						outreach1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmOutreachHandle"))
						outreach1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmOutreachHandleAll"))
						outreach1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmOutreachClear"))
					}
					login1 := alarm.Group("login/")
					{
						login1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmLoginPage"))
						login1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmLoginHandle"))
						login1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmLoginHandleAll"))
						login1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmLoginClear"))
					}
					self1 := alarm.Group("self/")
					{
						self1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmSelfPage"))
						self1.POST("handle", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmSelfHandle"))
						self1.POST("handleAll", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmSelfHandleAll"))
						self1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"AlarmSelfClear"))
					}
				} //告警

				//通信日志
				syslog := vOrg.Group("syslog/")
				{
					public := syslog.Group("public/")
					{
						public.POST("loginStatus", (&organization.SyslogController{}).LoginStatus)
					}
					dns1 := syslog.Group("dns/")
					{
						dns1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogDnsPage"))
						dns1.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogDnsBatchDelete"))
						dns1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogDnsClear"))
					}
					inner1 := syslog.Group("inner/")
					{
						inner1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogInnerPage"))
						inner1.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogInnerBatchDelete"))
						inner1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogInnerClear"))
					}
					outer1 := syslog.Group("outer/")
					{
						outer1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogOuterPage"))
						outer1.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogOuterBatchDelete"))
						outer1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogOuterClear"))
					}
					process1 := syslog.Group("process/")
					{
						process1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogProcessPage"))
						process1.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogProcessBatchDelete"))
						process1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogProcessClear"))
					}
					login1 := syslog.Group("login/")
					{
						login1.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogLoginPage"))
						login1.POST("batchDelete", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogLoginBatchDelete"))
						login1.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"SyslogLoginClear"))
					}
				} //通信日志

				//订单授权
				order := vOrg.Group("order/")
				{
					public := order.Group("public/")
					{
						public.POST("orderStatus", (&organization.OrderController{}).OrderStatus)
						public.POST("expireStatus", (&organization.OrderController{}).ExpireStatus)
						public.POST("grantDuration", (&organization.OrderController{}).GrantDuration)
						public.POST("payMethod", (&organization.OrderController{}).PayMethod)
						public.POST("payBankcard", (&organization.OrderController{}).PayBankcard)
						public.POST("expiredRemind", (&organization.OrderController{}).ExpiredRemind)
					}
					product := order.Group("product/")
					{
						product.POST("published", (&organization.OrderController{}).PublishedForProduct)
					}
					purchase := order.Group("purchase/")
					{
						purchase.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchasePage"))
						purchase.POST("info", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchaseInfo"))
						purchase.POST("logs", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchaseLogs"))
						purchase.POST("add", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchaseAdd"))
						purchase.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchaseCancel"))
						purchase.POST("delete", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchaseDelete"))
						purchase.POST("payByBankcard", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderPurchasePayByBankcard"))
					}
					grant := order.Group("grant/")
					{
						grant.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderGrantPage"))
						grant.POST("validList", (&organization.OrderController{}).ValidList)
						grant.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderGrantSave"))
					}
					used := order.Group("used/")
					{
						used.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderUsedPage"))
					}
					license := order.Group("license/")
					{
						license.POST("import", validatorFactory.Create(consts.ValidatorPrefixOrg+"OrderLicenseImport"))
					}
				} //订单授权

				//任务相关
				task := vOrg.Group("task/")
				{
					public := task.Group("public/")
					{
						public.POST("taskTypes", (&organization.TaskController{}).TaskTypes)
						public.POST("taskStatus", (&organization.TaskController{}).TaskStatus)
					}
					history := task.Group("history/")
					{
						history.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"TaskHistoryPage"))
						history.POST("clear", validatorFactory.Create(consts.ValidatorPrefixOrg+"TaskHistoryClear"))
					}
				} //任务相关

				//操作日志
				operationLog := vOrg.Group("operationLog/")
				{
					operationLog.POST("modules", (&organization.OperationController{}).Modules)
					operationLog.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"OperationLogPage"))
				} //操作日志

				//终端访问日志
				agentAccessLog := vOrg.Group("agentAccessLog/")
				{
					agentAccessLog.POST("status", (&organization.AgentAccessController{}).Status)
					agentAccessLog.POST("urls", (&organization.AgentAccessController{}).Urls)
					agentAccessLog.POST("page", validatorFactory.Create(consts.ValidatorPrefixOrg+"AgentAccessLogPage"))
				} //操作日志

				//软件相关
				software := vOrg.Group("software/")
				{
					software.POST("platformTypes", (&organization.SoftwareController{}).PlatformTypes)
					software.POST("publishedList", (&organization.SoftwareController{}).PublishedList)
					software.POST("offlineDownload", validatorFactory.Create(consts.ValidatorPrefixOrg+"SoftwareOfflineDownload"))
					software.POST("onlineDownload", validatorFactory.Create(consts.ValidatorPrefixOrg+"SoftwareOnlineDownload"))
				} //软件相关

				//安全报告
				safetyReport := vOrg.Group("safetyReport/")
				{
					alarmNotice := safetyReport.Group("alarmNotice/")
					{
						alarmNotice.POST("get", (&organization.SafetyReportController{}).GetForAlarmNotice)
						alarmNotice.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"SafetyReportAlarmNoticeSave"))
					}
					hardwareNotice := safetyReport.Group("hardwareNotice/")
					{
						hardwareNotice.POST("get", (&organization.SafetyReportController{}).GetForHardwareNotice)
						hardwareNotice.POST("save", validatorFactory.Create(consts.ValidatorPrefixOrg+"SafetyReportHardwareNoticeSave"))
					}
				} //安全报告
			}

		}
	}
	return router
}
