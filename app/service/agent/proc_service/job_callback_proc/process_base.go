package job_callback_proc

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers/date_helper"
)

type ProcessBase struct {
}

// getThreatenStatus 获取威胁情报状态
func (crd *ProcessBase) getThreatenStatus(hash string) (isExist bool, status int) {
	isExist, level := threaten_models.CreateThreatenCacheFactory().GetLevelById(hash)
	if !isExist {
		return isExist, enum_process.RelyStatusUnset
	}
	if enum_threaten.IsWhite(level) {
		return isExist, enum_process.RelyStatusWhite
	} else if enum_threaten.IsBlack(level) {
		return isExist, enum_process.RelyStatusBlack
	} else {
		return isExist, enum_process.RelyStatusUnset
	}
}

// getGlobalStatus 获取全局名单状态z
func (crd *ProcessBase) getGlobalStatus(orgId int64, hash string) (status int) {
	model := proc_models.CreateProcGlobalFactory().GetItemByOrgAndId(orgId, hash)
	if model == nil || model.Type == enum_process.RelyStatusUnset {
		return enum_process.RelyStatusUnset
	}
	return model.Type
}

// getProcStatus 获取最终状态
func (crd *ProcessBase) getProcStatus(globalStatus, selfStatus, threatenStatus, autoWhiteSwitch int) (status int, at string) {
	at = date_helper.CurDatetime()
	if globalStatus == enum_process.RelyStatusWhite {
		status = enum_process.ProcStatusWhite
	} else if globalStatus == enum_process.RelyStatusBlack {
		status = enum_process.ProcStatusBlack
	} else if selfStatus == enum_process.RelyStatusWhite {
		status = enum_process.ProcStatusWhite
	} else if selfStatus == enum_process.RelyStatusBlack {
		status = enum_process.ProcStatusBlack
	} else if selfStatus == enum_process.RelyStatusUnset {
		status = enum_process.ProcStatusUnset
	} else if autoWhiteSwitch == enums.SwitchOn {
		if threatenStatus == enum_process.RelyStatusWhite {
			status = enum_process.ProcStatusWhite
		} else if threatenStatus == enum_process.RelyStatusBlack {
			status = enum_process.ProcStatusBlack
		} else {
			status = enum_process.ProcStatusUnset
			at = date_helper.TimestampToDatetime(date_helper.AddTimestampByHour(-3))
		}
	} else {
		status = enum_process.ProcStatusUnset
		at = date_helper.TimestampToDatetime(date_helper.AddTimestampByHour(-3))
	}

	return status, at
}
