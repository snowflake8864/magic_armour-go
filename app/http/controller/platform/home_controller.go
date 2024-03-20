package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// HomeController 首页相关
type HomeController struct {
	ControllerBase
}

func (u *HomeController) GetAgentStatistic(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	response.SuccessReturnData(ctx, loginInfo)
}

func (u *HomeController) GetAlarmStatistic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetAlarmStatistic")
}
