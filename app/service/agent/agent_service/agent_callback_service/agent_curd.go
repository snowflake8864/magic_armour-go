package agent_callback_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/task_service"
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
func (crd *AgentCurd) BatchGrant(orgId int64, agentIds []string, orderNo string, grantType int, grantModule, productTitle, expireAt string, grantDuration int) bool {
	if enum_order.GetPrefixByOrderNo(orderNo) == enum_order.OrderNoPrefixCloud {
		//导入授权终端过期时间与证书过期时间相同
		//购买授权终端过期时间为当前时间加上授权时长
		//todo bug[10110] 改成已订单过期为准
		//expireAt = date_helper.TimestampToDatetimeByFormat(date_helper.AddTimestampByDays(0, grantDuration, 0), variable.DateFormatToday)

	}
	updates := map[string]interface{}{
		"grant_product_title": productTitle,
		"grant_config":        grantModule,
		"grant_order_no":      orderNo,
		"grant_status":        enum_agent.GrantStatusAuthorized,
		"grant_type":          grantType,
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
	types := enums.Get1stIsStrReturnSliceOfInt("type", enum_product.WebMapInt)
	for _, grantType := range types {
		_ = crd.doGrantAutoExpire(orgId, expireAt, grantType)
	}
	return nil
}

func (crd *AgentCurd) doGrantAutoExpire(orgId int64, expireAt string, grantType int) error {
	conditions := agent_models.SearchConditions{
		OrgId:            orgId,
		GrantStatus:      enum_agent.GrantStatusAuthorized,
		GrantType:        grantType,
		GrantExpireAtLte: expireAt,
	}
	agentIds := crd.masterModel.GetIdsByByConditions(conditions)
	if agentIds == nil || len(agentIds) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"grant_product_title": "",
		"grant_order_no":      "",
	}

	//测试授权到期初始化终端配置
	//正式授权到期保留授权配置
	if grantType == enum_product.TypeGold {
		updates["grant_status"] = enum_agent.GrantStatusExpired
	} else {
		//授权配置
		grantConfig := enum_product.ProductGrantConfig{}
		grantConfig.Default()
		updates["grant_config"] = json_helper.Encode(grantConfig)
		updates["grant_status"] = enum_agent.GrantStatusUnauthorized
		updates["common_time"] = enum_agent.HeartbeatCommonTime
		updates["fast_time"] = enum_agent.HeartbeatFastTime
		updates["process_switch"] = enums.SwitchOff
		updates["process_mode"] = enum_agent.ModeMonitor
		updates["dll_switch"] = enums.SwitchOff
		updates["dll_mode"] = enum_agent.ModeMonitor
		updates["tamper_switch"] = enums.SwitchOff
		updates["tamper_mode"] = enum_agent.ModeMonitor
		updates["blackmail_switch"] = enums.SwitchOff
		updates["blackmail_mode"] = enum_agent.ModeMonitor
		updates["port_switch"] = enums.SwitchOff
		updates["plug_switch"] = enums.SwitchOff
		updates["usb_switch"] = enums.SwitchOff
		updates["usb_mode"] = enum_agent.ModeMonitor
		updates["syslog_dns_switch"] = enums.SwitchOff
		updates["syslog_outer_switch"] = enums.SwitchOff
		updates["syslog_inner_switch"] = enums.SwitchOff
		updates["syslog_process_switch"] = enums.SwitchOff
		updates["syslog_login_switch"] = enums.SwitchOff
		updates["outreach_alarm_switch"] = enums.SwitchOff
		updates["login_solve_switch"] = enums.SwitchOff
	}

	if isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates); isOk {
		if grantType != enum_product.TypeGold {
			task_service.CreateCurdFactory().AddTaskForAgents(orgId, agentIds, enum_task.TypeAgentConfigDownload, "Grant Expiration")
		}
	}

	return nil
}
