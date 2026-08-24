package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type ProcSampleTrans struct {
	transformer.Common
	Lang string
}

// TransPage 转换成样本分页需要的数据
func (t *ProcSampleTrans) TransPage(total int64, models []*proc_models.ProcSampleModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("sample_status", model.Status, enum_process.EnumIntMap)

			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("sample_status_en", model.Status, enum_process.EnumIntMap)

			}

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
