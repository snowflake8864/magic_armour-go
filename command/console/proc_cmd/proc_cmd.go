package proc_cmd

import (
	"context"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/redis_factory"
)

func CreateCmdFactory() *ProcCmd {
	return &ProcCmd{
		queryClient: producer_kafka.CreateProducerFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "task_form_proc",
		cacheExpire: 86400,
	}
}

type ProcCmd struct {
	queryClient *producer_kafka.ProducerClient
	redisClient *redis_factory.RedisClient
	cacheKey    string
	cacheExpire int
}

func (crd *ProcCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *ProcCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnProc(true))
}

func (crd *ProcCmd) getAndSetTaskAt() string {
	preAt := crd.redisClient.Get(crd.cacheKey)
	if string_helper.IsEmpty(preAt) {
		preAt = dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformInstallAt)
	}
	crd.redisClient.SetEx(crd.cacheKey, date_helper.CurDatetime(), crd.cacheExpire)

	return preAt
}

// AddFromProcForTask  任务：从进程或模块添加
func (crd *ProcCmd) AddFromProcForTask() {
	variable.ZapLog.Sugar().Debug("AddFromProcForTask start")

	//获取上一次执行时间， 并保存当前时间
	preAt := crd.getAndSetTaskAt()

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Sugar().Debug("AddFromProcForTask orgModels empty")
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

			//进程
			if agentModel.ProcessSwitch == enums.SwitchOn {
				//进入队列
				producerFact.ProduceCrontab(ctx, &kafka_go.MessageStruct{
					Type:    kafka_go.TypeCrontabTaskFormProcess,
					OrgId:   orgModel.Id,
					AgentId: agentModel.AgentId,
					Data:    preAt,
				})
			}

			//模块
			if agentModel.DllSwitch == enums.SwitchOn {
				//进入队列
				producerFact.ProduceCrontab(ctx, &kafka_go.MessageStruct{
					Type:    kafka_go.TypeCrontabTaskFormDll,
					OrgId:   orgModel.Id,
					AgentId: agentModel.AgentId,
					Data:    preAt,
				})
			}
		}

	} //

	variable.ZapLog.Sugar().Debug("AddFromProcForTask end")
}

// ClearRunning  清空运行中数据
func (crd *ProcCmd) ClearRunning() {
	variable.ZapLog.Sugar().Debug("ClearRunning start")
	proc_models.CreateProcessRunningFactory().TruncateTable()
	proc_models.CreateDllRunningFactory().TruncateTable()
	variable.ZapLog.Sugar().Debug("ClearRunning end")
}
