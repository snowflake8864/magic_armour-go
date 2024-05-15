package agent_service

import (
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func CreateExtendFactory() *AgentExtend {
	return &AgentExtend{masterModel: agent_models.CreateAgentExtendFactory()}
}

type AgentExtend struct {
	AgentBase
	masterModel *agent_models.AgentExtendModel
}

// AddDefault 存在登录，不存在注册
func (crd *AgentExtend) AddDefault(orgId int64, agentId string) bool {
	orgExtendFact := organize_service.CreateOrgExtendFactory()
	orgConfig := orgExtendFact.GetItemByOrg(orgId)

	var param = agent_models.AgentExtendModel{}
	param.OrgId = orgId
	param.AgentId = agentId
	param.PlugPortConfig = json_helper.Encode(orgExtendFact.FormatPlugPortConfig(orgConfig.PlugPortConfig))
	param.PlugLoginConfig = json_helper.Encode(orgExtendFact.FormatPlugLoginConfig(orgConfig.PlugLoginConfig))
	param.LoginSolveConfig = json_helper.Encode(orgExtendFact.FormatLoginSolveConfig(orgConfig.LoginSolveConfig))
	isOk := crd.masterModel.InsertItem(&param)
	if !isOk {
		return false
	}
	return true
}
