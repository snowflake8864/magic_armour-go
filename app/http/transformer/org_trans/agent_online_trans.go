package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

type AgentOnline struct {
	transformer.Common
	Lang string
}

func (t *AgentOnline) pageGeneral(model *agent_online_models.OnlineLogModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["agent_id"] = model.AgentId
	data["host_name"] = model.HostName
	data["mac_id"] = model.MacId
	data["agent_ip"] = model.AgentIp
	data["net_ip"] = model.NetIp
	data["os_name"] = model.OsName
	data["ver"] = model.Ver
	data["is_online"] = model.IsOnline

	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []interface{}{}
	if departments != nil && len(departments) > 0 {
		if deptModel, ok := departments[model.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
			data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
		}
	}
	data["contact_name"] = model.ContactName
	data["contact_tel"] = model.ContactTel
	data["alias"] = model.Alias
	data["remark"] = model.Remark
	data["created_at"] = model.CreatedAt

	return data
}

// TransPage 转换成需要的数据
func (t *AgentOnline) TransPage(total int64, models []*agent_online_models.OnlineLogModel, departments map[int64]*department_models.DepartmentModel) interface{} {
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
