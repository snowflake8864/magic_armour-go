package agent_service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/kafka_go"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: model_agent.CreateAgentFactory()}
}

type AgentCurd struct {
	AgentBase
	masterModel *model_agent.AgentModel
}

// GetItemByOrgAndAgent 通过企业id和终端id获取一条记录
func (crd *AgentCurd) GetItemByOrgAndAgent(orgId int64, agentId, column string) (model *model_agent.AgentModel) {
	model, _ = crd.masterModel.GetItemByOrgAndAgent(orgId, agentId, column)
	return model
}

// GetItemByAgent 通过终端id获取一条记录
func (crd *AgentCurd) GetItemByAgent(id, column string) (model *model_agent.AgentModel) {
	model, _ = crd.masterModel.GetItemById(id, column)
	return model
}

// GetMapsByAgents 通过终端ids获取多条记录, 有映射关系
func (crd *AgentCurd) GetMapsByAgents(ids []string, column string) map[string]*model_agent.AgentModel {
	models, err := crd.masterModel.GetItemsByIds(ids, column)
	if err != nil {
		return nil
	}
	if models == nil {
		return nil
	}
	result := make(map[string]*model_agent.AgentModel)
	for _, model := range models {
		result[model.AgentId] = model
	}
	return result
}

// Page 分页
func (crd *AgentCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_agent.AgentModel) {
	where := crd.search2Where(search)
	order := ""
	if orderField == enum_agent.OrderFieldDefault {
		order = "is_online DESC, is_safe ASC, created_at DESC"
	} else {
		order = "created_at DESC"
	}
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// SetLivingByAgent 设置终端活动中
func (crd *AgentCurd) SetLivingByAgent(id string) bool {
	nowDate := helpers.GetDatetime()

	updates := map[string]interface{}{
		"live_at":   nowDate,
		"online_at": gorm.Expr("IF(is_online=?, ?, online_at)", enums.IsNot, nowDate),
		"is_online": enums.IsYes,
	}
	return crd.masterModel.UpdateById(id, updates)
}

// BatchDelete 批量删除
func (crd *AgentCurd) BatchDelete(ids []string) (errCode int, err error) {
	var validIds []string
	if agentModels, err := crd.masterModel.GetItemsByIds(ids, "agent_id"); err != nil || len(agentModels) == 0 {
		for _, model := range agentModels {
			validIds = append(validIds, model.AgentId)
		}
	}
	if len(validIds) == 0 {
		return -1, errors.New("终端不存在")
	}

	if isOK := crd.masterModel.DeleteByIds(validIds); !isOK {
		return -1, errors.New("删除失败")
	}

	//丢入队列
	ctx := context.Background()
	for _, val := range validIds {
		msg := &kafka_go.MessageStruct{
			Sign: kafka_go.SignAgent,
			Type: kafka_go.TypeAgentDelete,
			Data: val,
		}
		kafka_go.ProduceCommon(ctx, msg)
	}

	return errCode, err
}

// BatchUpgrade 批量升级
func (crd *AgentCurd) BatchUpgrade(ids []string) (errCode int, err error) {
	softwareMaps := software_service.CreateCurdFactory().GetMapsByPublished()
	if softwareMaps == nil || len(softwareMaps) == 0 {
		return errCode, err
	}
	var validIds []string
	if agentModels, err := crd.masterModel.GetItemsByIds(ids, "agent_id, os_type, ver"); err != nil || len(agentModels) == 0 {
		for _, model := range agentModels {
			if isOk, _ := enum_software.CanUpgradeIndirect(softwareMaps, model.OsType, model.Ver); isOk {
				validIds = append(validIds, model.AgentId)
			}
		}
	}

	if len(validIds) > 0 {
		//添加任务
		task_service.CreateCurdFactory().AddTasksForAgents(validIds, []int{enum_task.TypeAgentUpgrade}, "")
	}
	return errCode, err
}

// BatchUpdate 批量修改
func (crd *AgentCurd) BatchUpdate(ids []string, updateColumn map[string]struct{}, param *model_agent.AgentModel, remark string) (errCode int, err error) {
	agentModels, err := crd.masterModel.GetItemsByIds(ids, crd.ColumnConfig(true))
	if err != nil || len(agentModels) == 0 {
		return -1, errors.New("终端不存在")
	}

	return errCode, err
}
