package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent_access"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type AgentAccessTrans struct {
	transformer.Common
	Lang string
}

func (t *AgentAccessTrans) pageGeneral(model *agent_access_models.AccessLogModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	if t.Lang == consts.LanguageZh {
		data["path_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("path", model.Path, enum_agent_access.EnumStrMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_agent_access.EnumMapInt)
	} else {
		data["path_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("path_en", model.Path, enum_agent_access.EnumStrMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_agent_access.EnumMapInt)
	}

	return data
}

// TransPage 转换成管理页需要的数据
func (t *AgentAccessTrans) TransPage(total int64, models []*agent_access_models.AccessLogModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)

			tmp["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
