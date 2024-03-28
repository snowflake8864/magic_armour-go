package agent_service

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
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

// search2Where 查询条件转化
func (crd *AgentCurd) search2Where(search map[string]interface{}) (wheres map[string]interface{}) {
	if val, ok := search["agentId"]; ok {
		wheres["agent_id"] = val
	}
	if val, ok := search["hostName"]; ok {
		wheres["agent_id"] = val
	}

	return wheres
}

func (crd *AgentCurd) Page(search map[string]interface{}, offset, pageSize int, orderField, orderSort string) (total int, items []model_agent.AgentModel) {

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
