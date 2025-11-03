package agent_cmd

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/agent_service/agent_callback_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateAgentFactory() *AgentCmd {
	return &AgentCmd{}
}

type AgentCmd struct {
}

func (crd *AgentCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *AgentCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// Offline 离线
func (crd *AgentCmd) Offline() {
	variable.ZapLog.Info("Offline start")

	floorAt := date_helper.TimestampToDatetime(date_helper.AddTimestampByMinute(-6))
	variable.ZapLog.Info(fmt.Sprintf("Offline floorAt=%s", floorAt))

	agent_service.CreateCurdFactory().SetOffline(floorAt)

	variable.ZapLog.Info("Offline end")
}

// HandleSafe 处理安全状态
func (crd *AgentCmd) HandleSafe() {
	variable.ZapLog.Info("HandleSafe start")

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Info("HandleSafe orgModels empty")
		return
	}

	//今天
	statDay, _, _, statMonth := date_helper.CalcByDay(0)
	for _, orgModel := range orgModels {
		//获取不安全的终端
		condition := alarm_stat_models.SearchConditions{
			OrgId:     orgModel.Id,
			StatMonth: statMonth,
			StatDay:   statDay,
		}
		agentIds := alarm_stat_models.CreateStatDailyFactory().GetUnsafeAgentsByConditions(condition)

		//默认改成全部安全
		agent_service.CreateCurdFactory().SetSafeByAgents(orgModel.Id, nil, enums.IsYes)

		if len(agentIds) > 0 {
			//改成不安全的
			agent_service.CreateCurdFactory().SetSafeByAgents(orgModel.Id, agentIds, enums.IsNot)
		}
	}

	variable.ZapLog.Info("HandleSafe end")
}

// GrantExpired  授权已过期的
func (crd *AgentCmd) GrantExpired() {
	variable.ZapLog.Info("GrantExpired")
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	orgModels := crd.getOrgModels()
	if len(orgModels) > 0 {
		for _, model := range orgModels {
			_ = agent_callback_service.CreateCurdFactory().GrantAutoExpire(model.Id, yesterday)
		}
	}

	variable.ZapLog.Info("OrderExpired end dayEnd:" + yesterday)
}
