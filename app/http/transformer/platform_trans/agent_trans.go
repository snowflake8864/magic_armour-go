package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/utils/helpers"
)

type Agent struct {
	transformer.Common
}

func (t *Agent) pageGeneral(model *model_agent.AgentModel, departments map[int64]*model_organization.DepartmentModel) map[string]interface{} {
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
	data["grant_product_title"] = model.GrantProductTitle
	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []string{}
	if deptModel, ok := departments[model.DeptId]; ok && !helpers.StringIsEmpty(deptModel.FullName) {
		data["dept_fullname"] = helpers.Json2Map(deptModel.FullName)
	}

	data["created_at"] = model.CreatedAt
	return data
}

// TransAgentManagePage 转换成管理页需要的数据
func (t *Agent) TransAgentManagePage(total int64, models []*model_agent.AgentModel, departments map[int64]*model_organization.DepartmentModel) interface{} {
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
