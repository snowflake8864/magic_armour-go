package proc_cmd

import (
	"context"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
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
	isOnline := enums.IsYes
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId:    orgId,
		IsOnline: &isOnline,
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
	variable.ZapLog.Info("AddFromProcForTask start")

	//获取上一次执行时间， 并保存当前时间
	preAt := crd.getAndSetTaskAt()

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Info("AddFromProcForTask orgModels empty")
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
			//数据来源
			dataFrom := 0

			//进程
			if agentModel.ProcessSwitch == enums.SwitchOn {
				dataFrom += enum_process.DataFromProcess
			}

			//模块
			if agentModel.DllSwitch == enums.SwitchOn {
				dataFrom += enum_process.DataFromDll
			}

			if dataFrom > 0 {
				//进入队列
				producerFact.ProduceCrontab(ctx, &kafka_go.MessageStruct{
					Type:    kafka_go.TypeCrontabTaskFormProc,
					OrgId:   orgModel.Id,
					AgentId: agentModel.AgentId,
					AgentOS: int_helper.IntToStr(dataFrom),
					Data:    preAt,
				})
			}
		}

	} //

	variable.ZapLog.Info("AddFromProcForTask end")
}

// ClearRunning  清空运行中数据
func (crd *ProcCmd) ClearRunning() {
	variable.ZapLog.Info("ClearRunning start")
	proc_models.CreateProcessRunningFactory().TruncateTable()
	proc_models.CreateDllRunningFactory().TruncateTable()
	variable.ZapLog.Info("ClearRunning end")
}

// RepairStatusByThreaten  通过威胁情报修复状态
func (crd *ProcCmd) RepairStatusByThreaten() {
	variable.ZapLog.Info("RepairStatusByThreaten start")

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Info("RepairStatusByThreaten orgModels empty")
		return
	}

	processCurdFact := proc_service.CreateProcessCurdFactory()
	dllCurdFact := proc_service.CreateDllCurdFactory()
	for _, orgModel := range orgModels {
		variable.ZapLog.Info(fmt.Sprintf("RepairStatusByThreaten orgId=%d start", orgModel.Id))
		processCurdFact.ThreatenRepairStatus(orgModel.Id, enum_process.RelyStatusWhite)
		processCurdFact.ThreatenRepairStatus(orgModel.Id, enum_process.RelyStatusBlack)
		dllCurdFact.ThreatenRepairStatus(orgModel.Id, enum_process.RelyStatusWhite)
		dllCurdFact.ThreatenRepairStatus(orgModel.Id, enum_process.RelyStatusBlack)
		variable.ZapLog.Info(fmt.Sprintf("RepairStatusByThreaten orgId=%d end", orgModel.Id))
	}
	variable.ZapLog.Info("RepairStatusByThreaten end")
}

// SyncStatusByThreaten  通过威胁情报修复状态
func (crd *ProcCmd) SyncStatusByThreaten() {
	variable.ZapLog.Info("SyncStatusByThreaten start")

	threatenCacheFact := threaten_models.CreateCacheFactory()
	page := 0
	limit := 2000
	conditions := threaten_models.SearchConditions{
		LevelList: []int{enum_threaten.LevelWhite, enum_threaten.LevelBlack},
	}

	var processWhites, processBlacks, dllWhites, dllBlacks []string
	for {
		page++
		offset := (page - 1) * limit
		variable.ZapLog.Info(fmt.Sprintf("SyncStatusByThreaten query page=%d, offset=%d, limit=%d ", page, offset, limit))
		models := threatenCacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			variable.ZapLog.Info(fmt.Sprintf("SyncStatusByThreaten break page=%d, offset=%d, limit=%d ", page, offset, limit))
			break
		}

		for _, model := range models {
			if model.DataFrom == enum_threaten.DataFromProcess {
				if model.Level == enum_threaten.LevelWhite {
					processWhites = append(processWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					processBlacks = append(processBlacks, model.Id)
				}
			} else if model.DataFrom == enum_threaten.DataFromDll {
				if model.Level == enum_threaten.LevelWhite {
					dllWhites = append(dllWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					dllBlacks = append(dllBlacks, model.Id)
				}
			} else {
				if model.Level == enum_threaten.LevelWhite {
					processWhites = append(processWhites, model.Id)
					dllWhites = append(dllWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					processBlacks = append(processBlacks, model.Id)
					dllBlacks = append(dllBlacks, model.Id)
				}
			}
		}

	}
	variable.ZapLog.Info(fmt.Sprintf("SyncStatusByThreaten processWhites=%d, processBlacks=%d,dllWhites=%d, dllBlacks=%d",
		len(processWhites), len(processBlacks), len(dllWhites), len(dllBlacks)))

	//每次修改最大数
	supportNumber := 50
	//进程白名单
	if tmpLen := len(processWhites); tmpLen > 0 {
		variable.ZapLog.Info("SyncStatusByThreaten processWhites start")
		for i := 0; i < tmpLen; i += supportNumber {
			end := i + supportNumber
			if end > tmpLen {
				end = tmpLen
			}

			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHashes(processWhites[i:end], enum_threaten.LevelWhite)
		}
		variable.ZapLog.Info("SyncStatusByThreaten processWhites end")
	}

	//进程黑名单
	if tmpLen := len(processBlacks); tmpLen > 0 {
		variable.ZapLog.Info("SyncStatusByThreaten processBlacks start")
		for i := 0; i < tmpLen; i += supportNumber {
			end := i + supportNumber
			if end > tmpLen {
				end = tmpLen
			}

			proc_service.CreateProcessCurdFactory().ThreatenCallbackByHashes(processBlacks[i:end], enum_threaten.LevelBlack)
		}
		variable.ZapLog.Info("SyncStatusByThreaten processBlacks end")
	}

	//模块白名单
	if tmpLen := len(dllWhites); tmpLen > 0 {
		variable.ZapLog.Info("SyncStatusByThreaten dllWhites start")
		for i := 0; i < tmpLen; i += supportNumber {
			end := i + supportNumber
			if end > tmpLen {
				end = tmpLen
			}

			proc_service.CreateDllCurdFactory().ThreatenCallbackByHashes(dllWhites[i:end], enum_threaten.LevelWhite)
		}
		variable.ZapLog.Info("SyncStatusByThreaten dllWhites end")
	}

	//模块黑名单
	if tmpLen := len(dllBlacks); tmpLen > 0 {
		variable.ZapLog.Info("SyncStatusByThreaten dllBlacks start")
		for i := 0; i < tmpLen; i += supportNumber {
			end := i + supportNumber
			if end > tmpLen {
				end = tmpLen
			}

			proc_service.CreateDllCurdFactory().ThreatenCallbackByHashes(dllBlacks[i:end], enum_threaten.LevelBlack)
		}
		variable.ZapLog.Info("SyncStatusByThreaten dllBlacks end")
	}

	variable.ZapLog.Info("SyncStatusByThreaten end")
}
