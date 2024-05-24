package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/cur_userinfo"
)

// ControllerBase 基础
type ControllerBase struct {
}

// getLogin 获取登录数据
func (u *ControllerBase) getLogin(ctx *gin.Context) agent_service.AgentCacheLogin {
	return cur_userinfo.GetCurrentAgent(ctx)
}
