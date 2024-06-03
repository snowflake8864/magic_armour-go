package agent_service

import (
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/date_helper"
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

// GetLoginSolve 获取登录破解配置
func (crd *AgentExtend) GetLoginSolve(orgId int64, agentId string) (config enum_login_solve.LoginSolveConfig) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "login_solve_config")
	if model != nil && !string_helper.IsEmpty(model.LoginSolveConfig) {
		_ = json_helper.ToStruct(model.LoginSolveConfig, &config)
	}

	config.Format()
	return config
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

// SetFlushAt 设置刷新时间
func (crd *AgentExtend) SetFlushAt(orgId int64, agentId string, process bool) bool {
	updates := map[string]interface{}{}
	if process {
		updates["process_flush_at"] = date_helper.CurDatetime()
	} else {
		updates["dll_flush_at"] = date_helper.CurDatetime()
	}

	isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates)
	return isOk
}

// BatchSetProcessFlushAt 批量设置刷新时间
func (crd *AgentExtend) BatchSetProcessFlushAt(orgId int64, agentIds []string, process bool) bool {
	updates := map[string]interface{}{}
	if process {
		updates["process_flush_at"] = date_helper.CurDatetime()
	} else {
		updates["dll_flush_at"] = date_helper.CurDatetime()
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates)
	return isOk
}
