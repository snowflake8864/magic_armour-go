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
	response.SuccessReturnData(ctx, "AlarmTodayStat")
}

// 终端: 今天统计
func (u *HomeController) AgentTodayStat(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AgentTodayStat")
}

// 终端: 每天在线
func (u *HomeController) AgentDailyOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AgentDailyOnline")
}

// 用户: 今天统计
func (u *HomeController) UserTodayStat(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UserTodayStat")
}

// 用户: 每天在线
func (u *HomeController) UserDailyOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UserDailyOnline")
}

// 用户: 最新日志
func (u *HomeController) UserNewestLogs(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UserNewestLogs")
}

// 授权: 使用统计
func (u *HomeController) GrantUsedStat(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GrantUsedStat")
}
