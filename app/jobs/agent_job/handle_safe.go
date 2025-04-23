package agent_job

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func HandleSafe(orgId int64, agentId string) {
	statDay, _, _, statMonth := date_helper.CalcByDay(0)

	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s, statDay=%s, statMonth=%d", orgId, agentId, statDay, statMonth))
	condition := alarm_stat_models.SearchConditions{
		OrgId:     orgId,
		AgentId:   agentId,
		StatMonth: statMonth,
		StatDay:   statDay,
	}
	isSafe := enums.IsNot
	if model := alarm_stat_models.CreateStatDailyFactory().GetItemByCondition(condition); model != nil && (model.Total-model.Handled) <= 0 {
		isSafe = enums.IsYes
	}
	agent_service.CreateCurdFactory().SetSafeByAgent(orgId, agentId, isSafe)
}
