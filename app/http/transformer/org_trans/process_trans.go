package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
)

type Process struct {
	transformer.Common
}

// TransMasterPage 转换成全部进程分页需要的数据
func (t *Process) TransMasterPage(total int64, models []*proc_models.ProcessModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["proc_hash"] = model.ProcHash
			data["proc_dir"] = model.Dir
			data["status"] = model.Status
			data["sync_status_at"] = model.SyncStatusAt
			data["global_status"] = model.GlobalStatus
			data["self_status"] = model.SelfStatus
			data["threaten_status"] = model.ThreatenStatus
			data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
			data["threaten_level"] = model.ThreatenLevel
			data["threaten_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.ThreatenLevel, enum_threaten.EnumMap)
			data["threaten_is_sync"] = 0
			if enum_threaten.HasLevel(model.ThreatenLevel) {
				data["threaten_is_sync"] = 1
			}
			data["auto_white_switch"] = model.AutoWhiteSwitch
			data["grant_threaten_switch"] = model.GrantThreatenSwitch
			data["introduce"] = model.Introduce
			data["copyright"] = model.Copyright
			data["signature"] = model.Signature
			data["created_at"] = model.CreatedAt
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
			data["status_reason"] = enum_process.GetProcStatusReason(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			data["btn_change_enable"] = enum_process.BtnChangeEnable(model.GlobalStatus)
			data["btn_cancel_enable"] = enum_process.BtnCancelEnable(model.SelfStatus)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransRunningPage 转换成运行中进程分页需要的数据
func (t *Process) TransRunningPage(total int64, models []*proc_models.ProcessAndRunning) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["proc_hash"] = model.ProcHash
			data["proc_dir"] = model.ProcDir
			data["proc_id"] = model.ProcId
			data["procUser"] = model.ProcUser
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
			data["sync_status_at"] = model.SyncStatusAt
			data["global_status"] = model.GlobalStatus
			data["self_status"] = model.SelfStatus
			data["threaten_status"] = model.ThreatenStatus
			data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
			data["threaten_level"] = model.ThreatenLevel
			data["threaten_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("level", model.ThreatenLevel, enum_threaten.EnumMap)
			data["threaten_is_sync"] = 0
			if enum_threaten.HasLevel(model.ThreatenLevel) {
				data["threaten_is_sync"] = 1
			}
			data["auto_white_switch"] = model.AutoWhiteSwitch
			data["grant_threaten_switch"] = model.GrantThreatenSwitch
			data["introduce"] = model.Introduce
			data["copyright"] = model.Copyright
			data["signature"] = model.Signature
			data["created_at"] = model.CreatedAt

			data["status_reason"] = enum_process.GetProcStatusReason(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
			data["btn_change_enable"] = enum_process.BtnChangeEnable(model.GlobalStatus)
			data["btn_cancel_enable"] = enum_process.BtnCancelEnable(model.SelfStatus)

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransRepeatPage 转换成运行中进程分页需要的数据
func (t *Process) TransRepeatPage(total int64, models []*proc_models.ProcessModel, agentModels map[string]*agent_models.AgentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["proc_dir"] = model.Dir
			data["created_at"] = model.CreatedAt
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, nil)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
