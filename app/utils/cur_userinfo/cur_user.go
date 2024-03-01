package cur_userinfo

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
)

// GetCurrentAgent 获取当前终端
// @context 请求上下文
func GetCurrentAgent(context *gin.Context) agent_service.AgentCacheLogin {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	current, _ := context.MustGet(tokenKey).(agent_service.AgentCacheLogin)
	return current
}
