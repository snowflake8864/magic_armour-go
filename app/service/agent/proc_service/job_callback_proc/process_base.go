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
func (crd *ProcessBase) getThreatenStatus(hash string) (isExist bool, threatenStatus int, threatenLevel int) {
	isExist, threatenLevel = threaten_models.CreateCacheFactory().GetLevelById(hash)
	if !isExist {
		return isExist, enum_process.RelyStatusUnset, threatenLevel
	}

	if enum_threaten.IsWhite(threatenLevel) {
		return isExist, enum_process.RelyStatusWhite, threatenLevel
	} else if enum_threaten.IsBlack(threatenLevel) {
		return isExist, enum_process.RelyStatusBlack, threatenLevel
	}
	return isExist, enum_process.RelyStatusUnset, threatenLevel
}

// getGlobalStatus 获取全局名单状态
func (crd *ProcessBase) getGlobalStatus(orgId int64, hash string) (status int) {
	model := proc_models.CreateProcGlobalFactory().GetItemByOrgAndId(orgId, hash)
	if model == nil {
		return enum_process.RelyStatusUnset
	}
	if model.OpenSwitch == enums.SwitchOff {
		return enum_process.RelyStatusUnset
	}

	if model.Type == enum_process.GlobalTypeWhite {
		return enum_process.RelyStatusWhite
	} else if model.Type == enum_process.GlobalTypeBlack {
		return enum_process.RelyStatusBlack
	} else {
		return enum_process.RelyStatusUnset
	}

}

// getThreatenStatusByHashes 获取威胁情报状态
func (crd *ProcessBase) getThreatenStatusByHashes(hashes []string) map[string]int {
	result := make(map[string]int)

	threatenLevels := threaten_models.CreateCacheFactory().GetLevelsByIds(hashes)
	if len(threatenLevels) == 0 {
		return result
	}

	for hash, level := range threatenLevels {
		if enum_threaten.IsWhite(level) {
			result[hash] = enum_process.RelyStatusWhite
		} else if enum_threaten.IsBlack(level) {
			result[hash] = enum_process.RelyStatusBlack
		} else {
			result[hash] = enum_process.RelyStatusUnset
		}
	}
	return result

}

// getGlobalStatus 获取全局名单状态
func (crd *ProcessBase) getGlobalStatusByHashes(orgId int64, hashes []string) map[string]int {
	result := make(map[string]int)

	models := proc_models.CreateProcGlobalFactory().GetItemsByOrgAndIds(orgId, hashes)
	if models == nil {
		return result
	}
	for _, model := range models {
		if model.OpenSwitch == enums.SwitchOn {
			if model.Type == enum_process.GlobalTypeWhite {
				result[model.Id] = enum_process.RelyStatusWhite
			} else if model.Type == enum_process.GlobalTypeBlack {
				result[model.Id] = enum_process.RelyStatusBlack
			} else {
				result[model.Id] = enum_process.RelyStatusUnset
			}
		}
	}
	return result
}

// getProcStatus 获取最终状态
func (crd *ProcessBase) getProcStatus(globalStatus, selfStatus, threatenStatus, autoWhiteSwitch int) (status int, at int64) {
	at = date_helper.CurTimestamp()
	if globalStatus == enum_process.RelyStatusWhite {
		status = enum_process.ProcStatusWhite
	} else if globalStatus == enum_process.RelyStatusBlack {
		status = enum_process.ProcStatusBlack
	} else if selfStatus == enum_process.RelyStatusWhite {
		status = enum_process.ProcStatusWhite
	} else if selfStatus == enum_process.RelyStatusBlack {
		status = enum_process.ProcStatusBlack
	} else if selfStatus == enum_process.RelyStatusUnknown {
		status = enum_process.ProcStatusUnknown
	} else if autoWhiteSwitch == enums.SwitchOn {
		if threatenStatus == enum_process.RelyStatusWhite {
			status = enum_process.ProcStatusWhite
		} else if threatenStatus == enum_process.RelyStatusBlack {
			status = enum_process.ProcStatusBlack
		} else {
			status = enum_process.ProcStatusUnset
			at = 0
		}
	} else {
		status = enum_process.ProcStatusUnset
		at = 0
	}

	return status, at
}
