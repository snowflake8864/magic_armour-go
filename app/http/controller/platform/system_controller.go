package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/monitor_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
)

// SystemController 系统管理
type SystemController struct {
	ControllerBase
}

// MonitorRecords 监控记录
func (u *SystemController) MonitorRecords(ctx *gin.Context) {
	date := ctx.GetString(consts.ValidatorPrefixPlatform + "date")

	//页面上传
	conditions := monitor_service.SearchConditions{}
	conditions.MonitorDay = date

	models := monitor_service.CreateCurdFactory().GetList(conditions)

	response.SuccessReturnData(ctx, (&platform_trans.Monitor{}).TransList(models))
}

// TimeGet  服务器器时间
func (u *SystemController) TimeGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// TimeSet  服务器器时间
func (u *SystemController) TimeSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// SshGet  SSH服务
func (u *SystemController) SshGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// SshSet  SSH服务
func (u *SystemController) SshSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// NetworkGet  网络服务
func (u *SystemController) NetworkGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// NetworkSet  网络服务
func (u *SystemController) NetworkSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// AgentLogGet  终端日志
func (u *SystemController) AgentLogGet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}

// AgentLogSet  终端日志
func (u *SystemController) AgentLogSet(ctx *gin.Context) {
	response.SuccessReturnData(ctx, gin.H{
		"date": helpers.DatetimeCurr(),
	})
}
