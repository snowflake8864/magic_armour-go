package plug_cmd

import (
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/plug_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateCmdFactory() *PlugCmd {
	return &PlugCmd{
		plugService: plug_service.CreatePolicyDynamicCurdFactory(),
	}

}

type PlugCmd struct {
	plugService *plug_service.PolicyDynamicCurd
}

func (crd *PlugCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *PlugCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// DynamicExpired  动态封堵已过期的
func (crd *PlugCmd) DynamicExpired() {
	variable.ZapLog.Info("DynamicExpired start")

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Info("DynamicExpired orgModels empty")
		return
	}
	now := date_helper.CurDatetime()

	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			if err := crd.plugService.AutoDisable(orgModel.Id, agentModel.AgentId, now); err != nil {
				variable.ZapLog.Info("DynamicExpired end")
			}
		}

	} //

	variable.ZapLog.Info("DynamicExpired end")
}
