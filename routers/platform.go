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
	"goskeleton/app/utils/cur_app"
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
		vPlatform.GET("customConfig", func(context *gin.Context) {
			context.JSON(http.StatusOK, cur_app.GetPlatformConfig())
		})

		vPlatform.Use(authorization.RequestParams).Use(authorization.SavePlatformOperate())
		{
			//登录认证
			auth := vPlatform.Group("auth/")
			{
				auth.POST("checkUser", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthCheckUser"))
				auth.POST("getSliderCode", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthGetSliderCode"))
				auth.POST("sendSms", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthSendSms"))
				auth.POST("normalLogin", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AuthNormalLogin"))
				auth.POST("logout", (&platform.AuthController{}).Logout)
			}

			// 【需要token】中间件验证的路由
			vPlatform.Use(authorization.CheckPlatformPToken())
			{
				//文件上传公共路由
				uploadFiles := vPlatform.Group("upload/")
				{
					uploadFiles.POST("files", validatorFactory.Create(consts.ValidatorPrefixPlatform+"UploadFiles"))
				}
				//账户
				account := vPlatform.Group("account/")
				{
					account.POST("modifyPw", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AccountModifyPw"))
				}

				//数据字典
				dictionary := vPlatform.Group("dictionary/")
				{
					dictionary.POST("list", (&platform.DictionaryController{}).List)
					dictionary.POST("save", validatorFactory.Create(consts.ValidatorPrefixPlatform+"DictionarySave"))
				}

				//首页
				home := vPlatform.Group("home/")
				{
					home.POST("alarmTodayStat", (&platform.HomeController{}).AlarmTodayStat)
					home.POST("agentTodayStat", (&platform.HomeController{}).AgentTodayStat)
					home.POST("agentDailyOnline", validatorFactory.Create(consts.ValidatorPrefixPlatform+"HomeAgentDailyOnline"))
					home.POST("userTodayStat", (&platform.HomeController{}).UserTodayStat)
					home.POST("userDailyOnline", validatorFactory.Create(consts.ValidatorPrefixPlatform+"HomeUserDailyOnline"))
					home.POST("userNewestLogs", (&platform.HomeController{}).UserNewestLogs)
					home.POST("grantUsedStat", (&platform.HomeController{}).GrantUsedStat)
				}
				//菜单 - 隐藏页
				menu := vPlatform.Group("orgMenu/")
				{
					menu.POST("tree", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrgMenuTree"))
					menu.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrgMenuAdd"))
					menu.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrgMenuEdit"))
					menu.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrgMenuDelete"))
				}

				//企业
				organize := vPlatform.Group("organize/")
				{
					public := organize.Group("public/")
					{
						public.POST("types", (&platform.OrganizeController{}).Types)
						public.POST("detail", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizePublicDetail"))
					}
					manage := organize.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManagePage"))
						manage.POST("lock", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManageLock"))
						manage.POST("unlock", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManageUnlock"))
						manage.POST("clean", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManageClean"))
						manage.POST("setHsSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManageSetHsSwitch"))
						manage.POST("setLoginSmsSwitch", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeManageSetLoginSmsSwitch"))

					}
					verify := organize.Group("verify/")
					{
						verify.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeVerifyPage"))
						verify.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeVerifyAdd"))
						verify.POST("approve", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeVerifyApprove"))
						verify.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeVerifyRefuse"))
					}
					user := organize.Group("user/")
					{
						user.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeUserPage"))
					}
					log := organize.Group("log/")
					{
						log.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrganizeLogPage"))
					}

				} //企业

				//产品
				product := vPlatform.Group("product/")
				{
					public := product.Group("public/")
					{
						public.POST("grantTypes", (&platform.ProductController{}).GrantTypes)
						public.POST("grantModules", (&platform.ProductController{}).GrantModules)
						public.POST("allPublished", (&platform.ProductController{}).AllPublished)
					}
					manage := product.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProductManagePage"))
						manage.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProductManageAdd"))
						manage.POST("edit", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProductManageEdit"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ProductManageDelete"))

					}
				} //产品

				//订单
				order := vPlatform.Group("order/")
				{
					public := order.Group("public/")
					{
						public.POST("orderStatus", (&platform.OrderController{}).OrderStatus)
						public.POST("expireStatus", (&platform.OrderController{}).ExpireStatus)
						public.POST("grantDuration", (&platform.OrderController{}).GrantDuration)
						public.POST("validByOrg", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderPublicValidByOrg"))
					}
					grant := order.Group("grant/")
					{
						grant.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderGrantPage"))
					}
					manage := order.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManagePage"))
						manage.POST("detail", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManageDetail"))
						manage.POST("logs", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManageLogs"))
						manage.POST("decline", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManageDecline"))
						manage.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManageCancel"))
						manage.POST("paid", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManagePaid"))
						manage.POST("unpaid", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OrderManageUnpaid"))
					}
				} //订单

				//终端
				agent := vPlatform.Group("agent/")
				{
					public := agent.Group("public/")
					{
						public.POST("onlineStatus", (&platform.AgentController{}).OnlineStatus)
						public.POST("safeStatus", (&platform.AgentController{}).SafeStatus)
						public.POST("statByOrg", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentPublicStatByOrg"))
					}
					org := agent.Group("org/")
					{
						org.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentOrgPage"))
						org.POST("addGrant", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentOrgAddGrant"))
					}
					own := agent.Group("own/")
					{
						own.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentOwnPage"))
					}

				} //终端

				//软件
				software := vPlatform.Group("software/")
				{
					public := software.Group("public/")
					{
						public.POST("osTypes", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwarePublicOsTypes"))
						public.POST("agentConfigGet", (&platform.SoftwareController{}).GetAgentConfig)
						public.POST("agentConfigSet", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwarePublicAgentConfigSet"))
					}
					manage := software.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManagePage"))
						manage.POST("add", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManageAdd"))
						manage.POST("download", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManageDownload"))
						manage.POST("publish", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManagePublish"))
						manage.POST("cancel", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManageCancel"))
						manage.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SoftwareManageDelete"))
						manage.POST("clear", (&platform.SoftwareController{}).ClearForManage)
					}

				}

				//威胁情报
				threaten := vPlatform.Group("threaten/")
				{
					public := threaten.Group("public/")
					{
						public.POST("levels", (&platform.ThreatenController{}).Levels)
					}
					config := threaten.Group("config/")
					{
						config.POST("get", (&platform.ThreatenController{}).GetForConfig)
						config.POST("save", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenConfigSave"))
						config.POST("online", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenConfigOnline"))
					}
					cache := threaten.Group("cache/")
					{
						cache.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenCachePage"))
						cache.POST("getAttr", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenCacheGetAttr"))
						cache.POST("delete", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenCacheDelete"))
						cache.POST("batchWhite", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenCacheBatchWhite"))
						cache.POST("batchBlack", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenCacheBatchBlack"))
						cache.POST("allWhite", (&platform.ThreatenController{}).AllWhiteForPush)
						cache.POST("clear", (&platform.ThreatenController{}).ClearForCache)
						cache.POST("refresh", (&platform.ThreatenController{}).RefreshForCache)
					}
					push := threaten.Group("push/")
					{
						push.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenPushPage"))
						push.POST("clear", (&platform.ThreatenController{}).ClearForPush)
					}
				} //威胁情报
				//操作相关
				operation := vPlatform.Group("operation/")
				{
					operation.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"OperationPage"))
					operation.POST("modules", (&platform.OperationController{}).Modules)
				}
				//服务器相关
				system := vPlatform.Group("system/")
				{
					monitor := system.Group("monitor/")
					{
						monitor.POST("records", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemMonitorRecords"))
					}
					time := system.Group("time/")
					{
						time.POST("get", (&platform.SystemController{}).TimeGet)
						time.POST("set", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemTimeSet"))
					}
					ssh := system.Group("ssh/")
					{
						ssh.POST("get", (&platform.SystemController{}).SshGet)
						ssh.POST("set", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemSshSet"))
					}
					network := system.Group("network/")
					{
						network.POST("get", (&platform.SystemController{}).NetworkGet)
						network.POST("set", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemNetworkSet"))
					}
					agentLog := system.Group("agentLogConfig/")
					{
						agentLog.POST("get", (&platform.SystemController{}).AgentLogConfigGet)
						agentLog.POST("set", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemAgentLogConfigSet"))
					}
					huanshi := system.Group("huanshi/")
					{
						huanshi.POST("get", (&platform.SystemController{}).AgentLogConfigGet)
						huanshi.POST("online", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemHuanshiOnline"))
						huanshi.POST("set", validatorFactory.Create(consts.ValidatorPrefixPlatform+"SystemHuanshiSet"))
					}
				} //服务器相关
			}
		}
	}
	return router
}
