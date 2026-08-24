package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type AgentOnlineTrans struct {
	transformer.Common
	Lang string
}

func (t *AgentOnlineTrans) pageGeneral(model *agent_online_models.OnlineLogModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

	return data
}

// TransPage 转换成需要的数据
func (t *AgentOnlineTrans) TransPage(total int64, models []*agent_online_models.OnlineLogModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
