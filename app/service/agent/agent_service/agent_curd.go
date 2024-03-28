package agent_service

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: model_agent.CreateAgentFactory()}
}

type AgentCurd struct {
	AgentBase
	masterModel *model_agent.AgentModel
}

// GetItemByAgent 通过终端id获取一条记录
func (crd *AgentCurd) GetItemByAgent(agentId, column string) (model *model_agent.AgentModel) {
	model, _ = crd.masterModel.GetItemById(agentId, column)
	return model
}

// GetMapsByAgents 通过终端ids获取多条记录, 有映射关系
func (crd *AgentCurd) GetMapsByAgents(agentIds []string, column string) map[string]model_agent.AgentModel {
	models, err := crd.masterModel.GetItemsByIds(agentIds, column)
	if err != nil {
		return nil
	}
	if models == nil {
		return nil
	}
	result := make(map[string]model_agent.AgentModel)
	for _, model := range models {
		result[model.AgentId] = model
	}
	return result
}

func (crd *AgentCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []model_agent.AgentModel) {
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
func (crd *AgentCurd) SetLivingByAgent(agentId string) bool {
	nowDate := helpers.GetDatetime()

	updates := map[string]interface{}{
		"live_at":   nowDate,
		"online_at": gorm.Expr("IF(is_online=?, ?, online_at)", enums.IsNot, nowDate),
		"is_online": enums.IsYes,
	}
	return crd.masterModel.UpdateById(agentId, updates)
}
