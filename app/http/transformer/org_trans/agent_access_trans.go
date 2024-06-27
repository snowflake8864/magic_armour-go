package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent_access"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
)

type AgentAccess struct {
	transformer.Common
}

func (t *AgentAccess) pageGeneral(model *agent_access_models.AccessLogModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["agent_id"] = model.AgentId
	data["path"] = model.Path
	data["path_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("path", model.Path, enum_agent_access.EnumStrMap)
	data["status"] = model.Status
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_agent_access.EnumMapInt)

	data["net_ip"] = model.NetIp
	data["input"] = model.Input
	data["method"] = model.Method
	data["http_code"] = model.HttpCode
	data["resp"] = model.Resp
	data["delay_time"] = model.DelayTime

	data["created_at"] = model.CreatedAt

	return data
}

// TransPage 转换成管理页需要的数据
func (t *AgentAccess) TransPage(total int64, models []*agent_access_models.AccessLogModel,
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
