package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/utils/response"
)

// AlarmController 告警相关
type AlarmController struct {
	ControllerBase
}

// PutMaster 上传报警信息(除端口虚开外)
func (u *AlarmController) PutMaster(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "alert")
	err := agent_post_service.CreateAlarmFactory().ToQueryBYComm(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutPort 上传端口虚开告警
func (u *AlarmController) PutPort(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "alert")
	err := agent_post_service.CreateAlarmFactory().ToQueryBYPort(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutVirus 上传病毒查杀告警
func (u *AlarmController) PutVirus(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "alert")
	err := agent_post_service.CreateAlarmFactory().ToQueryBYVirus(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
