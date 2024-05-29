package agent_service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: agent_models.CreateAgentFactory()}
}

type AgentCurd struct {
	AgentBase
	masterModel *agent_models.AgentModel
}

// SetLivingByAgent 设置终端活动中
func (crd *AgentCurd) SetLivingByAgent(orgId int64, id string) bool {
	nowDate := date_helper.CurDatetime()

	updates := map[string]interface{}{
		"live_at":   nowDate,
		"online_at": gorm.Expr("IF(is_online=?, ?, online_at)", enums.IsNot, nowDate),
		"is_online": enums.IsYes,
	}
	isOk, _ := crd.masterModel.UpdateByOrgAndAgent(orgId, id, updates)
	return isOk
}

// BatchDelete 批量删除
func (crd *AgentCurd) BatchDelete(orgId int64, ids []string) (errCode int, err error) {
	var validIds []string
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnSimplify()); len(agentModels) == 0 {
		for _, model := range agentModels {
			validIds = append(validIds, model.AgentId)
		}
	}
	if len(validIds) == 0 {
		return -1, errors.New("终端不存在")
	}

	if isOK := crd.masterModel.DeleteByOrgAndAgents(orgId, validIds); !isOK {
		return -1, errors.New("删除失败")
	}

	//丢入队列
	ctx := context.Background()
	for _, val := range validIds {
		msg := &kafka_go.MessageStruct{
			Type:    kafka_go.TypeAgentDelete,
			OrgId:   orgId,
			AgentId: val,
		}
		producer_kafka.CreateProducerFactory().ProduceCommon(ctx, msg)
	}

	return errCode, err
}

// BatchUpgrade 批量升级
func (crd *AgentCurd) BatchUpgrade(orgId int64, ids []string) (errCode int, err error) {
	softwareMaps := software_service.CreateCurdFactory().GetMapsByPublished()
	if softwareMaps == nil || len(softwareMaps) == 0 {
		return errCode, err
	}
	var validIds []string
	if agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnProc(true)); len(agentModels) == 0 {
		for _, model := range agentModels {
			if isOk, _ := enum_software.CanUpgradeIndirect(softwareMaps, model.OsType, model.Ver); isOk {
				validIds = append(validIds, model.AgentId)
			}
		}
	}

	if len(validIds) > 0 {
		//添加任务
		task_service.CreateCurdFactory().AddTasksForAgents(orgId, validIds, []int{enum_task.TypeAgentUpgrade}, "")
	}
	return errCode, err
}

// BatchUpdate 批量修改
func (crd *AgentCurd) BatchUpdate(orgId int64, ids []string, updateColumn UpdateColumn, remark string) (errCode int, err error) {
	agentModels := crd.masterModel.GetItemsByOrgAndAgents(orgId, ids, crd.masterModel.ColumnConfig(true))
	if len(agentModels) == 0 {
		return -1, errors.New("终端不存在")
	}

	return errCode, err
}
