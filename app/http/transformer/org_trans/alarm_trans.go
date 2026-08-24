package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type AlarmTrans struct {
	transformer.Common
	Lang string
}

// TransProcessPage 转换成进程告警分页需要的数据
func (t *AlarmTrans) TransProcessPage(total int64, models []*alarm_models.ProcessModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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

// TransDllPage 转换成模块告警分页需要的数据
func (t *AlarmTrans) TransDllPage(total int64, models []*alarm_models.DllModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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

// TransBlackmailPage 转换成勒索告警分页需要的数据
func (t *AlarmTrans) TransBlackmailPage(total int64, models []*alarm_models.BlackmailModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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

// TransTamperPage 转换成防篡改告警分页需要的数据
func (t *AlarmTrans) TransTamperPage(total int64, models []*alarm_models.TamperModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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

// TransPortPage 转换成端口虚开告警分页需要的数据
func (t *AlarmTrans) TransPortPage(total int64, models []*alarm_models.PortModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
			}
			redirectPort := ""
			if model.RedirectPort > 0 {
				redirectPort = int_helper.IntToStr(model.RedirectPort)
			}
			data["redirect_port"] = redirectPort

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
				data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", model.PlugType, enum_plug.EnumMap)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
				data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type_en", model.PlugType, enum_plug.EnumMap)
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

// TransUsbPage 转换成外设告警分页需要的数据
func (t *AlarmTrans) TransUsbPage(total int64, models []*alarm_models.UsbModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel,
	agentUsbAlias map[string]map[string]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
			}

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
			}

			data["proc_hash"] = model.ProcHash
			if model.ProcId > 0 {
				data["proc_id"] = model.ProcId
			} else {
				data["proc_id"] = ""
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

// TransLoginPage 转换成口令破解告警分页需要的数据
func (t *AlarmTrans) TransLoginPage(total int64, models []*alarm_models.LoginModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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

// TransSelfPage 转换成自保告警分页需要的数据
func (t *AlarmTrans) TransSelfPage(total int64, models []*alarm_models.SelfModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status", model.HandleStatus, enum_alarm.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_alarm.EnumMapInt)
				data["level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level_en", model.Level, enum_alarm.EnumMapInt)
				data["handle_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("handle_status_en", model.HandleStatus, enum_alarm.EnumMapInt)
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
