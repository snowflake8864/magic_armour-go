package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Agent struct {
	transformer.Common
	Lang string
}

func (t *Agent) pageGeneral(model *agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	data["grant_config"] = enum_product.Json2Config(model.GrantConfig)
	data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

	return data
}

// TransAgentManagePage 转换成管理页需要的数据
func (t *Agent) TransAgentManagePage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
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
