package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/int_helper"
)

type AlarmTrans struct {
	transformer.Common
}

// TransProcessPage 转换成进程告警分页需要的数据
func (t *AlarmTrans) TransProcessPage(total int64, models []*alarm_models.ProcessModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_id"] = model.ProcId
			data["proc_dir"] = model.ProcDir
			data["proc_param"] = model.ProcParam

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransDllPage 转换成模块告警分页需要的数据
func (t *AlarmTrans) TransDllPage(total int64, models []*alarm_models.DllModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["dll_hash"] = model.DllHash
			data["dll_dir"] = model.DllDir
			data["proc_dir"] = model.ProcDir

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransBlackmailPage 转换成勒索告警分页需要的数据
func (t *AlarmTrans) TransBlackmailPage(total int64, models []*alarm_models.BlackmailModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_dir"] = model.ProcDir
			data["file_dir"] = model.FileDir

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransTamperPage 转换成防篡改告警分页需要的数据
func (t *AlarmTrans) TransTamperPage(total int64, models []*alarm_models.TamperModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_dir"] = model.ProcDir
			data["file_dir"] = model.FileDir

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransPortPage 转换成端口虚开告警分页需要的数据
func (t *AlarmTrans) TransPortPage(total int64, models []*alarm_models.PortModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["weight"] = model.Weight
			data["attack_ip"] = model.AttackIp
			data["target_ip"] = model.TargetIp
			data["target_port"] = model.TargetPort
			data["redirect_ip"] = model.RedirectIp
			redirectPort := ""
			if model.RedirectPort > 0 {
				redirectPort = int_helper.IntToStr(model.RedirectPort)
			}
			data["redirect_port"] = redirectPort
			data["attack_total"] = model.AttackTotal

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransPlugPage 转换成封堵告警分页需要的数据
func (t *AlarmTrans) TransPlugPage(total int64, models []*alarm_models.PlugModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["ip"] = model.Ip
			data["plug_type"] = model.PlugType
			data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", model.PlugType, enum_plug.EnumMap)
			data["direction"] = model.Direction
			data["duration"] = model.Duration
			data["start_time"] = model.StartTime
			data["end_time"] = model.EndTime

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransUsbPage 转换成外设告警分页需要的数据
func (t *AlarmTrans) TransUsbPage(total int64, models []*alarm_models.UsbModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel,
	agentUsbAlias map[string]map[string]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["usb_hash"] = model.UsbHash
			data["usb_name"] = model.UsbName
			data["usb_intro"] = model.UsbIntro
			data["usb_alias"] = ""
			if hashAlias, isHash := agentUsbAlias[model.AgentId][model.UsbHash]; isHash {
				data["usb_alias"] = hashAlias
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

// TransOutreachPage 转换成外联告警分页需要的数据
func (t *AlarmTrans) TransOutreachPage(total int64, models []*alarm_models.OutreachModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			if model.ProcId > 0 {
				data["proc_id"] = model.ProcId
			} else {
				data["proc_id"] = ""
			}
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

// TransLoginPage 转换成口令破解告警分页需要的数据
func (t *AlarmTrans) TransLoginPage(total int64, models []*alarm_models.LoginModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["source_ip"] = model.SourceIp
			data["solve_total"] = model.SolveTotal

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransSelfPage 转换成自保告警分页需要的数据
func (t *AlarmTrans) TransSelfPage(total int64, models []*alarm_models.SelfModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["proc_hash"] = model.ProcHash
			data["proc_dir"] = model.ProcDir
			data["proc_param"] = model.ProcParam
			data["file_dir"] = model.FileDir
			data["target_dir"] = model.TargetDir

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransVirusPage 转换成病毒告警分页需要的数据
func (t *AlarmTrans) TransVirusPage(total int64, models []*alarm_models.VirusModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})

			data["id"] = fmt.Sprintf("%d", model.Id)
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
			data["level"] = model.Level
			data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
			data["agent_at"] = model.AgentAt
			data["handle_status"] = model.HandleStatus
			data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			data["handle_remark"] = model.HandleRemark
			data["handle_user"] = model.HandleUser
			data["handle_at"] = model.HandleAt
			data["created_at"] = model.CreatedAt

			data["file_hash"] = model.FileHash
			data["file_dir"] = model.FileDir
			data["virus_type"] = model.VirusType
			data["virus_desc"] = model.VirusDesc

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
