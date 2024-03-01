package cur_userinfo

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
	"goskeleton/app/service/agent/agent_service"
)

// GetCurrentUserId 获取当前用户的id
// @context 请求上下文
func GetCurrentUserId(context *gin.Context) (int64, bool) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	currentUser, exist := context.MustGet(tokenKey).(my_jwt.CustomClaims)
	return currentUser.UserId, exist
}

// GetCurrentAgent 获取当前终端
// @context 请求上下文
func GetCurrentAgent(context *gin.Context) agent_service.AgentCacheLogin {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	current, _ := context.MustGet(tokenKey).(agent_service.AgentCacheLogin)
	return current
}
