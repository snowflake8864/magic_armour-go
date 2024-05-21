package agent_service

import (
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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

// GetPlugOfPort 获取封堵配置： 端口虚开
func (crd *AgentExtend) GetPlugOfPort(orgId int64, agentId string) (config enum_plug.PlugPortConfig) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "plug_port_config")
	if model != nil && !string_helper.IsEmpty(model.PlugPortConfig) {
		_ = json_helper.ToStruct(model.PlugPortConfig, &config)
	}

	config.Format()
	return config
}
