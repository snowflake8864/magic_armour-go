package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/map_helper"
)

type ProcSample struct {
	transformer.Common
}

// TransPage 转换成样本分页需要的数据
func (t *ProcSample) TransPage(total int64, models []*proc_models.ProcSampleModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["dir"] = model.Dir
			data["original_name"] = model.OriginalName
			data["data_from"] = model.DataFrom
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("sample_status", model.Status, enum_process.EnumIntMap)
			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt

			agentInfo := t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, map_helper.MergeStringAndInterface(data, agentInfo))
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
