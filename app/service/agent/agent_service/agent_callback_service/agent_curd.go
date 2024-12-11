package agent_callback_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: agent_models.CreateAgentFactory()}
}

type AgentCurd struct {
	masterModel *agent_models.AgentModel
}

// DeleteByOrg 通过企业删除
func (crd *AgentCurd) DeleteByOrg(orgId int64) bool {
	return crd.masterModel.DeleteByOrg(orgId)
}

// BatchGrant 批量授权
func (crd *AgentCurd) BatchGrant(orgId int64, agentIds []string, orderNo string, grantDuration int, grantModule string, productTitle string, expireAt string) bool {
	if enum_order.GetTypeByNum(orderNo) == enum_order.OrderNoPrefixCloud {
		//导入授权终端过期时间与证书过期时间相同
		//购买授权终端过期时间为当前时间加上授权时长
		expireAt, _, _, _ = date_helper.CalcByDay(grantDuration)
	}
	updates := map[string]interface{}{
		"grant_product_title": productTitle,
		"grant_config":        grantModule,
		"grant_order_no":      orderNo,
		"grant_status":        enum_agent.GrantStatusAuthorized,
		"grant_expire_at":     expireAt,
	}
	if isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates); isOk {
		return true
	} else {
		return false
	}
}

// GrantAutoExpire 授权自动到期
func (crd *AgentCurd) GrantAutoExpire(orgId int64, expireAt string) error {
	conditions := agent_models.SearchConditions{
		OrgId:            orgId,
		GrantStatus:      enum_agent.GrantStatusAuthorized,
		GrantExpireAtLte: expireAt,
	}
	agentIds := crd.masterModel.GetIdsByByConditions(conditions)
	if agentIds == nil || len(agentIds) == 0 {
		return nil
	}
	//授权配置
	grantConfig := enum_product.ProductGrantConfig{}
	grantConfig.Default()
	updates := map[string]interface{}{
		"grant_product_title":   "",
		"grant_config":          json_helper.Encode(grantConfig),
		"grant_order_no":        "",
		"grant_status":          enum_agent.GrantStatusUnauthorized,
		"grant_expire_at":       "",
		"common_time":           enum_agent.HeartbeatCommonTime,
		"fast_time":             enum_agent.HeartbeatFastTime,
		"process_switch":        enums.SwitchOff,
		"process_mode":          enum_agent.ModeMonitor,
		"dll_switch":            enums.SwitchOff,
		"dll_mode":              enum_agent.ModeMonitor,
		"tamper_switch":         enums.SwitchOff,
		"tamper_mode":           enum_agent.ModeMonitor,
		"blackmail_switch":      enums.SwitchOff,
		"blackmail_mode":        enum_agent.ModeMonitor,
		"plug_switch":           enums.SwitchOff,
		"usb_switch":            enums.SwitchOff,
		"usb_mode":              enum_agent.ModeMonitor,
		"syslog_dns_switch":     enums.SwitchOff,
		"syslog_outer_switch":   enums.SwitchOff,
		"syslog_inner_switch":   enums.SwitchOff,
		"syslog_process_switch": enums.SwitchOff,
		"syslog_login_switch":   enums.SwitchOff,
		"outreach_alarm_switch": enums.SwitchOff,
		"login_solve_switch":    enums.SwitchOff,
	}
	if isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates); isOk {
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, agentIds, enum_task.TypeAgentConfigDownload, "授权到期")
	}

	return nil
}
