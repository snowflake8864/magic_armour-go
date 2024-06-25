package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_syslog"
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
			data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.TargetPortStatus, enum_syslog.EnumMapInt)
			data["source_ip"] = model.SourceIp
			data["source_port"] = model.SourcePort
			data["proto"] = model.Proto
			data["proto_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.Proto, enum_syslog.EnumMapInt)

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransOuterPage 转换成对外通讯分页需要的数据
func (t *SyslogTrans) TransOuterPage(total int64, models []*syslog_models.OuterModel,
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
			data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.TargetPortStatus, enum_syslog.EnumMapInt)
			data["source_ip"] = model.SourceIp
			data["source_port"] = model.SourcePort
			data["proto"] = model.Proto
			data["proto_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.Proto, enum_syslog.EnumMapInt)

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransProcessPage 转换成进程分页需要的数据
func (t *SyslogTrans) TransProcessPage(total int64, models []*syslog_models.ProcessModel,
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
			data["proc_param"] = model.ProcParam
			data["parent_hash"] = model.ParentHash
			data["parent_id"] = model.ParentId
			data["parent_dir"] = model.ParentDir
			data["parent_param"] = model.ParentParam

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransLoginPage 转换成登录分页需要的数据
func (t *SyslogTrans) TransLoginPage(total int64, models []*syslog_models.LoginModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["agent_at"] = model.AgentAt
			data["created_at"] = model.CreatedAt

			data["source_ip"] = model.SourceIp
			data["username"] = model.Username
			data["type"] = model.Type
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("login_status", model.Status, enum_syslog.EnumMapInt)

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
