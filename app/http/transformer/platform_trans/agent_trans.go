package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

type Agent struct {
	transformer.Common
}

func (t *Agent) pageGeneral(model *agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["agent_id"] = model.AgentId
	data["host_name"] = model.HostName
	data["mac_id"] = model.MacId
	data["agent_ip"] = model.AgentIp
	data["net_ip"] = model.NetIp
	data["os_name"] = model.OsName
	data["platform_type"] = model.PlatformType
	data["os_type"] = model.OsType
	data["is_online"] = model.IsOnline
	data["online_at"] = model.IsOnline
	data["is_safe"] = model.IsSafe
	data["contact_name"] = model.ContactName
	data["alias"] = model.Alias
	data["grant_order_no"] = model.GrantOrderNo
	data["grant_product_title"] = model.GrantProductTitle
	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []string{}
	if departments != nil && len(departments) > 0 {
		if deptModel, ok := departments[model.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
			data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
		}
	}

	data["created_at"] = model.CreatedAt
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
