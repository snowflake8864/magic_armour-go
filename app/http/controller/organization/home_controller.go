package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/response"
)

// HomeController 首页相关
type HomeController struct {
	ControllerBase
}

// AgentStat 终端：实时统计
func (u *HomeController) AgentStat(ctx *gin.Context) {

	loginInfo := u.getLogin(ctx)
	conditions := agent_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.IsOnline = enums.AllDefaultInt
	conditions.IsSafe = enums.AllDefaultInt
	err, agentIds := u.getCurrUserAgentIds(u.getLogin(ctx))

	if err != nil {
		//不是管理员 & 无终端
		response.SuccessReturnData(ctx, agent_models.AgentCalc{})
		return
	} else {
		//不是管理员有终端
		if agentIds != nil {
			conditions.AgentIds = agentIds
		}
	}
	//统计
	agentFact := agent_models.CreateAgentFactory()
	result := agentFact.CalcByConditions(conditions, agentFact.ColumnCalcAgents(false, true, true, true, true))

	response.SuccessReturnData(ctx, agentFact.FillCalcData(result))
}

// AlarmAgentTops 告警：终端排行
func (u *HomeController) AlarmAgentTops(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	response.SuccessReturnData(ctx, loginInfo)
}

// AlarmTodayStat 告警： 今天统计
func (u *HomeController) AlarmTodayStat(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	response.SuccessReturnData(ctx, loginInfo)
}

// AlarmDailyStat 告警： 每天统计
func (u *HomeController) AlarmDailyStat(ctx *gin.Context) {
	month := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "month"))

	response.SuccessReturnData(ctx, month)
}

// GrantUsedStat 授权： 使用统计
func (u *HomeController) GrantUsedStat(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetAlarmStatistic")
}
