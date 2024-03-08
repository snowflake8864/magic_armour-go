package cur_userinfo

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/platform/admin_service"
)

// GetCurrentAgent 获取当前终端
// @context 请求上下文
func GetCurrentAgent(context *gin.Context) (current agent_service.AgentCacheLogin) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	current, exist := context.MustGet(tokenKey).(agent_service.AgentCacheLogin)
	if exist {
		return current
	}
	return current
}

// GetCurrentAdmin 获取当前管理员
// @context 请求上下文
func GetCurrentAdmin(context *gin.Context) (current admin_service.AdminCacheLogin) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	current, exist := context.MustGet(tokenKey).(admin_service.AdminCacheLogin)
	if exist {
		return current
	}
	return current
}
