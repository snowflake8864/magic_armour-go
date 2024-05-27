package middleware_auth_org

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/response"
)

// ApiIsEnabled 接口是否启用
func ApiIsEnabled() gin.HandlerFunc {
	return func(context *gin.Context) {
		if cur_app.GetAppUsageVersion() == variable.AppUsageVersionOffline {
			response.ErrorCasbinAuthFail(context, "api disabled."+context.Request.RequestURI)
			return
		}

		context.Next()
	}
}
