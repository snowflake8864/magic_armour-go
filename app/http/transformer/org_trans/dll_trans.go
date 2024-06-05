package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
)

type Dll struct {
	transformer.Common
}

// TransMasterPage 转换成全部进程分页需要的数据
func (t *Dll) TransMasterPage(total int64, models []*proc_models.DllModel) interface{} {
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
			data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
			data["auto_white_switch"] = model.AutoWhiteSwitch
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
func (t *Dll) TransRunningPage(total int64, models []*proc_models.DllAndRunning) interface{} {
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
			data["threaten_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rely_status", model.ThreatenStatus, enum_process.EnumIntMap)
			data["auto_white_switch"] = model.AutoWhiteSwitch
			data["introduce"] = model.Introduce
			data["copyright"] = model.Copyright
			data["signature"] = model.Signature
			data["attribute"] = model.Attribute
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
			//终端信息
			data["host_name"] = ""
			data["mac_id"] = ""
			data["agent_ip"] = ""
			data["net_ip"] = ""
			data["os_name"] = ""
			data["contact_name"] = ""
			data["contact_tel"] = ""
			data["alias"] = ""
			if agentModel, isOK := agentModels[model.AgentId]; isOK {
				data["host_name"] = agentModel.HostName
				data["mac_id"] = agentModel.MacId
				data["agent_ip"] = agentModel.AgentIp
				data["net_ip"] = agentModel.NetIp
				data["os_name"] = agentModel.OsName
				data["contact_name"] = agentModel.ContactName
				data["contact_tel"] = agentModel.ContactTel
				data["alias"] = agentModel.Alias
			}

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
