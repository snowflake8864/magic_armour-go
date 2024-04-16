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
				auth.POST("logout", (&platform.AuthController{}).Logout)
			}

			// 【需要token】中间件验证的路由
			vPlatform.Use(authorization.CheckPlatformPToken())
			{
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

				//终端
				agent := vPlatform.Group("agent/")
				{
					public := agent.Group("public/")
					{
						public.POST("onlineStatus", (&platform.AgentController{}).OnlineStatus)
						public.POST("safeStatus", (&platform.AgentController{}).SafeStatus)
						public.POST("countOnline", (&platform.AgentController{}).CountOnline)
					}
					manage := agent.Group("manage/")
					{
						manage.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"AgentManagePage"))

					}

				} //终端
				//威胁情报
				threaten := vPlatform.Group("threaten/")
				{
					public := threaten.Group("public/")
					{
						public.POST("levels", (&platform.ThreatenCacheController{}).Levels)
					}
					push := threaten.Group("push/")
					{
						push.POST("page", validatorFactory.Create(consts.ValidatorPrefixPlatform+"ThreatenPushPage"))
						push.POST("clear", (&platform.ThreatenPushController{}).Clear)
					}
				} //终端
			}
		}
	}
	return router
}
