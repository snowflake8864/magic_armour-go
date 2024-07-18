package cur_userinfo

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/admin_service"
)

// GetCurrentAgent 获取当前终端
// @context 请求上下文
func GetCurrentAgent(context *gin.Context) *agent_service.AgentCacheLogin {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	value, exists := context.Get(tokenKey)
	if !exists {
		return nil
	}
	if current, isOk := value.(agent_service.AgentCacheLogin); isOk {
		return &current
	} else {
		return nil
	}
}

// GetCurrentAdmin 获取当前管理员
// @context 请求上下文
func GetCurrentAdmin(context *gin.Context) *admin_service.AdminCacheLogin {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	value, exists := context.Get(tokenKey)
	if !exists {
		return nil
	}
	if current, isOk := value.(admin_service.AdminCacheLogin); isOk {
		return &current
	} else {
		return nil
	}
}

// GetCurrentUser 获取当前企业用户
// @context 请求上下文
func GetCurrentUser(context *gin.Context) *user_service.UserCacheLogin {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	value, exists := context.Get(tokenKey)
	if !exists {
		return nil
	}
	if current, isOk := value.(user_service.UserCacheLogin); isOk {
		return &current
	} else {
		return nil
	}

}
