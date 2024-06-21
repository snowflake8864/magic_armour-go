package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/map_helper"
)

type AlarmTrans struct {
	transformer.Common
}

// TransProcessPage 转换成进程告警分页需要的数据
func (t *AlarmTrans) TransProcessPage(total int64, models []*alarm_models.ProcessModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_id"] = model.ProcId
			data["proc_dir"] = model.ProcDir
			data["proc_param"] = model.ProcParam

			agentInfo := t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, map_helper.MergeStringAndInterface(data, agentInfo))
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
