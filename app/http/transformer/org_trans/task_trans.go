package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/task_models"
)

type Task struct {
	transformer.Common
}

func (t *Task) pageGeneral(model *task_models.TaskModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = fmt.Sprintf("%d", model.Id)
	data["agent_id"] = model.AgentId
	data["type"] = model.Type
	data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_task.EnumMap)
	data["status"] = model.Status
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_task.EnumMap)
	data["remark"] = model.Remark
	data["net_ip"] = model.NetIp
	data["created_at"] = model.CreatedAt
	data["completed_at"] = model.CompletedAt

	return data
}

// TransPage 转换成管理页需要的数据
func (t *Task) TransPage(total int64, models []*task_models.TaskModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
