package agent_callback_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/helpers/date_helper"
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

// DeleteByOrgAndAgent 通过企业和终端删除
func (crd *AgentCurd) DeleteByOrgAndAgent(orgId int64, id string) bool {
	return crd.masterModel.DeleteByOrgAndAgent(orgId, id)
}

// BatchGrant 批量授权
func (crd *AgentCurd) BatchGrant(orgId int64, agentIds []string, orderNo string, grantDuration int, grantModule string, productTitle string) bool {

	updates := map[string]interface{}{
		"grant_product_title": productTitle,
		"grant_config":        grantModule,
		"grant_order_no":      orderNo,
		"grant_status":        enums.IsYes,
		"grant_expire_at":     date_helper.TimestampToDatetime(date_helper.AddTimestampByDays(0, grantDuration, 0)),
	}
	if isOk, _ := crd.masterModel.UpdateByOrgAndAgents(orgId, agentIds, updates); isOk {
		return true
	} else {
		return false
	}
}
