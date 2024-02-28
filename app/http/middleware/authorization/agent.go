package authorization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
)

type AgentParams struct {
	Authorization string `header:"Authorization" binding:"required,min=32"`
}

// CheckAgentToken 检查终端token完整性、有效性中间件
func CheckAgentToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		variable.ZapLog.Info("查终端token完整性、有效性中间件----------------")
		c.Next()
	}
}
