package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// HomeController 首页相关
type HomeController struct {
	ControllerBase
}

// 告警: 今天统计
func (u *HomeController) AlarmTodayStat(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	response.SuccessReturnData(ctx, loginInfo)
}

// 终端: 每天在线
func (u *HomeController) AgentDailyOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AgentDailyOnline")
}
