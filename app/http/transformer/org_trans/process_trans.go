package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/proc_models"
)

type Process struct {
	transformer.Common
}

func (t *Process) pageGeneral(model *proc_models.ProcessModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["agent_id"] = model.AgentId
	data["proc_hash"] = model.ProcHash
	data["dir"] = model.Dir
	data["status"] = model.Status
	data["sync_status_at"] = model.SyncStatusAt
	data["global_status"] = model.GlobalStatus
	data["self_status"] = model.SelfStatus
	data["threaten_status"] = model.ThreatenStatus
	data["auto_white_switch"] = model.AutoWhiteSwitch
	data["introduce"] = model.Introduce
	data["copyright"] = model.Copyright
	data["signature"] = model.Signature
	data["created_at"] = model.CreatedAt
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("all_status", model.Status, enum_process.EnumIntMap)
	data["status_reason"] = enum_process.GetProcStatusReason(model.GlobalStatus, model.SelfStatus, model.ThreatenStatus, model.AutoWhiteSwitch)
	data["btn_change_enable"] = enum_process.BtnChangeEnable(model.GlobalStatus)
	data["btn_cancel_enable"] = enum_process.BtnCancelEnable(model.SelfStatus)

	return data
}

// TransMasterPage 转换成管理分页需要的数据
func (t *Process) TransMasterPage(total int64, models []*proc_models.ProcessModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			tmp := t.pageGeneral(model)

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
