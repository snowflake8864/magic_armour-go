package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/monitor_service"
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
