package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
			//终端信息
			data["host_name"] = ""
			data["mac_id"] = ""
			data["agent_ip"] = ""
			data["net_ip"] = ""
			data["os_name"] = ""
			data["contact_name"] = ""
			data["contact_tel"] = ""
			data["alias"] = ""
			data["dept_id"] = 0
			data["dept_fullname"] = []interface{}{}
			if agentModel, isOK := agentModels[model.AgentId]; isOK {
				data["host_name"] = agentModel.HostName
				data["mac_id"] = agentModel.MacId
				data["agent_ip"] = agentModel.AgentIp
				data["net_ip"] = agentModel.NetIp
				data["os_name"] = agentModel.OsName
				data["contact_name"] = agentModel.ContactName
				data["contact_tel"] = agentModel.ContactTel
				data["alias"] = agentModel.Alias
				if departments != nil && len(departments) > 0 {
					if deptModel, ok := departments[agentModel.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
						data["dept_id"] = deptModel.Id
						data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
					}
				}
			}

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
