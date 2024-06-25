package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/syslog/syslog_models"
)

type SyslogTrans struct {
	transformer.Common
}

// TransDnsPage 转换成DNS分页需要的数据
func (t *SyslogTrans) TransDnsPage(total int64, models []*syslog_models.DnsModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["agent_at"] = model.AgentAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_id"] = model.ProcId
			data["proc_dir"] = model.ProcDir
			data["domain_name"] = model.DomainName
			data["target_ip"] = model.TargetIp

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransInnerPage 转换成对内通讯分页需要的数据
func (t *SyslogTrans) TransInnerPage(total int64, models []*syslog_models.InnerModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["agent_at"] = model.AgentAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_id"] = model.ProcId
			data["proc_dir"] = model.ProcDir
			data["target_ip"] = model.TargetIp
			data["target_port"] = model.TargetPort
			data["target_port_status"] = model.TargetPortStatus
			data["target_port_status"] = model.TargetPortStatus

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
