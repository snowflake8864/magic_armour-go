package agent_cmd

import (
	"context"
	"fmt"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/agent_service/agent_callback_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
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

func (crd *AgentCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// Offline 离线
func (crd *AgentCmd) Offline() {
	variable.ZapLog.Debug("Offline start")

	floorAt := date_helper.TimestampToDatetime(date_helper.AddTimestampByMinute(-6))
	variable.ZapLog.Debug(fmt.Sprintf("Offline floorAt=%s", floorAt))

	agent_service.CreateCurdFactory().SetOffline(floorAt)

	variable.ZapLog.Debug("Offline end")
}

// HandleSafe 处理安全状态
func (crd *AgentCmd) HandleSafe() {
	variable.ZapLog.Debug("HandleSafe start")

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Debug("HandleSafe orgModels empty")
		return
	}

	ctx := context.Background()
	producerFact := producer_kafka.CreateProducerFactory()

	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			//进入队列
			producerFact.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabAgentSafe,
				OrgId:   orgModel.Id,
				AgentId: agentModel.AgentId,
			})
		}

	} //

	variable.ZapLog.Debug("HandleSafe end")
}

// GrantExpired  授权已过期的
func (crd *AgentCmd) GrantExpired() {
	variable.ZapLog.Debug("GrantExpired")
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	orgModels := crd.getOrgModels()
	if len(orgModels) > 0 {
		for _, model := range orgModels {
			_ = agent_callback_service.CreateCurdFactory().GrantAutoExpire(model.Id, yesterday)
		}
	}

	variable.ZapLog.Debug("OrderExpired end dayEnd:" + yesterday)
}
