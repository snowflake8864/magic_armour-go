package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers/map_helper"
)

type Dll struct {
	transformer.Common
	Lang string
}

// TransMasterPage 转换成全部进程分页需要的数据
func (t *Dll) TransMasterPage(total int64, models []*proc_models.DllModel, threatenMaps map[string]*threaten_models.CacheModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["dll_hash"] = model.DllHash
			data["dll_dir"] = model.Dir
			data["status"] = model.Status
			data["sync_status_at"] = model.SyncStatusAt
			data["global_status"] = model.GlobalStatus
			data["self_status"] = model.SelfStatus
			data["threaten_status"] = model.ThreatenStatus
			data["auto_white_switch"] = model.AutoWhiteSwitch
			data["grant_threaten_switch"] = model.GrantThreatenSwitch
			data["introduce"] = model.Introduce
			data["copyright"] = model.Copyright
			data["signature"] = model.Signature
			data["created_at"] = model.CreatedAt

			data["btn_change_enable"] = enum_process.BtnChangeEnable(model.GlobalStatus)
			data["btn_cancel_enable"] = enum_process.BtnCancelEnable(model.SelfStatus)

			if t.Lang == consts.LanguageZh {
				data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
				data["status_reason"] = enum_process.GetProcStatusReason(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			} else {
				data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status_en", model.ThreatenStatus, enum_process.EnumIntMap)
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status_en", model.Status, enum_process.EnumIntMap)
				data["status_reason"] = enum_process.GetProcStatusReasonEn(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			}
			//合并威胁情报数据
			data = map_helper.MergeStringAndInterface(data, t.AppendThreatenInfo(model.DllHash, threatenMaps))
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransRunningPage 转换成运行中进程分页需要的数据
func (t *Dll) TransRunningPage(total int64, models []*proc_models.DllAndRunning, threatenMaps map[string]*threaten_models.CacheModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["dll_hash"] = model.DllHash
			data["dll_dir"] = model.DllDir
			data["status"] = model.Status
			data["sync_status_at"] = model.SyncStatusAt
			data["global_status"] = model.GlobalStatus
			data["self_status"] = model.SelfStatus
			data["threaten_status"] = model.ThreatenStatus
			data["auto_white_switch"] = model.AutoWhiteSwitch
			data["grant_threaten_switch"] = model.GrantThreatenSwitch
			data["introduce"] = model.Introduce
			data["copyright"] = model.Copyright
			data["signature"] = model.Signature
			data["attribute"] = model.Attribute
			data["created_at"] = model.CreatedAt
			data["btn_change_enable"] = enum_process.BtnChangeEnable(model.GlobalStatus)
			data["btn_cancel_enable"] = enum_process.BtnCancelEnable(model.SelfStatus)

			if t.Lang == consts.LanguageZh {
				data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
				data["status_reason"] = enum_process.GetProcStatusReason(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			} else {
				data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status_en", model.ThreatenStatus, enum_process.EnumIntMap)
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status_en", model.Status, enum_process.EnumIntMap)
				data["status_reason"] = enum_process.GetProcStatusReasonEn(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			}

			//合并威胁情报数据
			data = map_helper.MergeStringAndInterface(data, t.AppendThreatenInfo(model.DllHash, threatenMaps))

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransRepeatPage 转换成运行中进程分页需要的数据
func (t *Dll) TransRepeatPage(total int64, models []*proc_models.DllModel, agentModels map[string]*agent_models.AgentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["dll_dir"] = model.Dir
			data["created_at"] = model.CreatedAt
			data["status"] = model.Status

			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status_en", model.Status, enum_process.EnumIntMap)
			}

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, nil)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
