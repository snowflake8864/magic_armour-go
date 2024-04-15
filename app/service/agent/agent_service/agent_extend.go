package agent_service

import (
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers"
)

func CreateExtendFactory() *AgentExtend {
	return &AgentExtend{masterModel: model_agent.CreateAgentExtendFactory()}
}

type AgentExtend struct {
	AgentBase
	masterModel *model_agent.AgentExtendModel
}

// AddDefault 存在登录，不存在注册
func (crd *AgentExtend) AddDefault(orgId int64, agentId string) bool {
	orgExtendFact := organize_service.CreateOrgExtendFactory()
	orgConfig := orgExtendFact.GetItemByOrg(orgId)

	var param = model_agent.AgentExtendModel{}
	param.OrgId = orgId
	param.AgentId = agentId
	param.PlugPortConfig = helpers.JsonEncode(orgExtendFact.FormatPlugPortConfig(orgConfig.PlugPortConfig))
	param.PlugLoginConfig = helpers.JsonEncode(orgExtendFact.FormatPlugLoginConfig(orgConfig.PlugLoginConfig))
	param.LoginSolveConfig = helpers.JsonEncode(orgExtendFact.FormatLoginSolveConfig(orgConfig.LoginSolveConfig))
	isOk := crd.masterModel.InsertItem(&param)
	if !isOk {
		return false
	}
	return true
}
