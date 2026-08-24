package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_syslog"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type SyslogTrans struct {
	transformer.Common
	Lang string
}

// TransDnsPage 转换成DNS分页需要的数据
func (t *SyslogTrans) TransDnsPage(total int64, models []*syslog_models.DnsModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["proto_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.Proto, enum_syslog.EnumMapInt)

			if t.Lang == consts.LanguageZh {
				data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.TargetPortStatus, enum_syslog.EnumMapInt)
			} else {
				data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status_en", model.TargetPortStatus, enum_syslog.EnumMapInt)

			}

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status", model.TargetPortStatus, enum_syslog.EnumMapInt)
			} else {
				data["target_port_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("port_status_en", model.TargetPortStatus, enum_syslog.EnumMapInt)

			}

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("login_status", model.Status, enum_syslog.EnumMapInt)
			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("login_status_en", model.Status, enum_syslog.EnumMapInt)
			}

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
