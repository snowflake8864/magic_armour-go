package agent_cmd

import (
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/organization/organization_models"
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

// Offline 离线
func (crd *AgentCmd) Offline() {
	variable.ZapLog.Sugar().Debug("Offline")
}

// HandleSafe 处理安全状态
func (crd *AgentCmd) HandleSafe() {
	variable.ZapLog.Sugar().Debug("HandleSafe")
}

// GrantExpired  授权已过期的
func (crd *AgentCmd) GrantExpired() {
	variable.ZapLog.Sugar().Debug("GrantExpired")
	_, _, dayEnd, _ := date_helper.CalcByDay(-1)

	orgModels := crd.getOrgModels()
	if len(orgModels) > 0 {
		for _, model := range orgModels {
			_ = agent_callback_service.CreateCurdFactory().GrantAutoExpire(model.Id, dayEnd)
		}
	}

	variable.ZapLog.Sugar().Debug("UnpaidExpired end dayEnd:" + dayEnd)
}
