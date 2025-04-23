package plug_cmd

import (
	"context"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateCmdFactory() *PlugCmd {
	return &PlugCmd{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type PlugCmd struct {
	queryClient *producer_kafka.ProducerClient
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
	variable.ZapLog.Debug("DynamicExpired start")

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Error("HandleSafe orgModels empty")
		return
	}

	ctx := context.Background()
	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			//进入队列
			crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabPlugDisable,
				OrgId:   orgModel.Id,
				AgentId: agentModel.AgentId,
			})
		}

	} //

	variable.ZapLog.Debug("DynamicExpired end")
}
