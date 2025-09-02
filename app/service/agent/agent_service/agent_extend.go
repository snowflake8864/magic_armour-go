package agent_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_ip_jump"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateExtendFactory() *AgentExtend {
	return &AgentExtend{masterModel: agent_models.CreateExtendFactory()}
}

type AgentExtend struct {
	AgentBase
	masterModel *agent_models.ExtendModel
}

// AddDefault 存在登录，不存在注册
func (crd *AgentExtend) AddDefault(orgId int64, agentId string) bool {
	orgExtendFact := organize_service.CreateOrgExtendFactory()
	orgConfig := orgExtendFact.GetItemByOrg(orgId)

	var param = agent_models.ExtendModel{}
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

// SetLoginSolve 设置登录破解配置
func (crd *AgentExtend) SetLoginSolve(orgId int64, agentId string, config enum_login_solve.LoginSolveConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["login_solve_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// BatchSetLoginSolve 批量设置登录破解配置
func (crd *AgentExtend) BatchSetLoginSolve(orgId int64, agentIds []string, config enum_login_solve.LoginSolveConfig) error {
	conditions := agent_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}

	config.Format()
	updates := make(map[string]interface{})
	updates["login_solve_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByConditions(conditions, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
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

// GetPlugOfLogin 获取封堵配置： 口令破解
func (crd *AgentExtend) GetPlugOfLogin(orgId int64, agentId string) (config enum_plug.PlugLoginConfig) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "plug_login_config")
	if model != nil && !string_helper.IsEmpty(model.PlugLoginConfig) {
		_ = json_helper.ToStruct(model.PlugLoginConfig, &config)
	}

	config.Format()
	return config
}

// SetPlugOfPort 设置封堵配置： 端口虚开
func (crd *AgentExtend) SetPlugOfPort(orgId int64, agentId string, config enum_plug.PlugPortConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["plug_port_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// BatchSetPlugOfPort 批量设置封堵配置： 端口虚开
func (crd *AgentExtend) BatchSetPlugOfPort(orgId int64, agentIds []string, config enum_plug.PlugPortConfig) error {
	conditions := agent_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}

	config.Format()
	updates := make(map[string]interface{})
	updates["plug_port_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByConditions(conditions, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetPlugOfLogin 设置封堵配置： 口令破解
func (crd *AgentExtend) SetPlugOfLogin(orgId int64, agentId string, config enum_plug.PlugLoginConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["plug_login_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// BatchSetPlugOfLogin 批量设置封堵配置：  口令破解
func (crd *AgentExtend) BatchSetPlugOfLogin(orgId int64, agentIds []string, config enum_plug.PlugLoginConfig) error {
	conditions := agent_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}

	config.Format()
	updates := make(map[string]interface{})
	updates["plug_login_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByConditions(conditions, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
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

// BatchSetFlushAt 批量设置刷新时间
func (crd *AgentExtend) BatchSetFlushAt(orgId int64, agentIds []string, process bool) bool {
	updates := map[string]interface{}{}
	if process {
		updates["process_flush_at"] = date_helper.CurDatetime()
	} else {
		updates["dll_flush_at"] = date_helper.CurDatetime()
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates)
	return isOk
}

// GetPwJumpConfig 获取密码跳变
func (crd *AgentExtend) GetPwJumpConfig(orgId int64, agentId string) (config enum_pw_jump.PwJumpConfig) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "pw_jump_config")
	if model != nil && !string_helper.IsEmpty(model.PwJumpConfig) {
		_ = json_helper.ToStruct(model.PwJumpConfig, &config)
	}

	config.Format()
	return config
}

// SetPwJumpConfig 设置密码跳变
func (crd *AgentExtend) SetPwJumpConfig(orgId int64, agentId string, config enum_pw_jump.PwJumpConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["pw_jump_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// GetPwJumpNewest 获取密码跳变最新结果
func (crd *AgentExtend) GetPwJumpNewest(orgId int64, agentId string) (newest enum_pw_jump.PwJumpNewest) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "pw_jump_newest")
	if model != nil && !string_helper.IsEmpty(model.PwJumpNewest) {
		_ = json_helper.ToStruct(model.PwJumpNewest, &newest)
	}
	newest.Format()
	return newest
}

// SetPwJumpNewest 设置密码跳变最新结果
func (crd *AgentExtend) SetPwJumpNewest(orgId int64, agentId string, newest enum_pw_jump.PwJumpNewest) error {
	updates := make(map[string]interface{})
	if newest.Status == enum_pw_jump.NewestStatusFail {
		prefix := crd.GetPwJumpNewest(orgId, agentId)
		newest.User = prefix.User
		newest.Pw = prefix.Pw
	} else {
		newest.Reason = ""
	}

	updates["pw_jump_newest"] = json_helper.Encode(newest)
	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// GetIpJumpConfig 获取IP跳变
func (crd *AgentExtend) GetIpJumpConfig(orgId int64, agentId string) (config enum_ip_jump.IpJumpConfig) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "ip_jump_config")
	if model != nil && !string_helper.IsEmpty(model.IpJumpConfig) {
		_ = json_helper.ToStruct(model.IpJumpConfig, &config)
	}

	config.Format()
	return config
}

// SetIpJumpConfig 设置IP跳变
func (crd *AgentExtend) SetIpJumpConfig(orgId int64, agentId string, config enum_ip_jump.IpJumpConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["ip_jump_config"] = json_helper.Encode(config)

	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// GetIpJumpNewest 获取地址跳变最新结果
func (crd *AgentExtend) GetIpJumpNewest(orgId int64, agentId string) (newest enum_ip_jump.IpJumpNewest) {
	model := crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, "ip_jump_newest")
	if model != nil && !string_helper.IsEmpty(model.IpJumpNewest) {
		_ = json_helper.ToStruct(model.IpJumpNewest, &newest)
	}
	newest.Format()
	return newest
}

// SetIpJumpNewest 设置地址码跳变最新结果
func (crd *AgentExtend) SetIpJumpNewest(orgId int64, agentId string, newest enum_ip_jump.IpJumpNewest) error {
	updates := make(map[string]interface{})
	if newest.Status == enum_ip_jump.NewestStatusFail {
		prefix := crd.GetIpJumpNewest(orgId, agentId)
		newest.SourceIp = prefix.SourceIp
		newest.TargetIp = prefix.TargetIp
		newest.Gateway = prefix.Gateway
		newest.AgentIp = prefix.AgentIp
	} else {
		newest.Reason = ""
	}

	updates["ip_jump_newest"] = json_helper.Encode(newest)
	if isOK, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, agentId, updates); isOK {
		return nil
	} else {
		return errors.New("保存失败")
	}
}
