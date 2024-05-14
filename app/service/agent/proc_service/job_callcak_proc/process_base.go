package job_callcak_proc

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/utils/helpers"
	"strings"
)

type ProcessBase struct {
}

// ColumnCalcOfUsage 字段： 统计使用数
func (crd *ProcessBase) ColumnCalcOfUsage() string {
	var columns []string
	columns = append(columns, "SUM(buy_number) AS buy_number")
	columns = append(columns, "SUM(usage_number) AS usage_number")
	return strings.Join(columns, ", ")
}

// getThreatenStatus 获取威胁情报状态
func (crd *ProcessBase) getThreatenStatus(hash string) (isExist bool, status int) {
	isExist, level := threaten_service.CreateCacheFactory().GetLevelById(hash)
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
	model, err := proc_models.CreateProcGlobalFactory().GetItemByOrgAndId(orgId, hash)
	if err != nil || model == nil || model.Type == enum_process.RelyStatusUnset {
		return enum_process.RelyStatusUnset
	}
	return model.Type
}

// getProcStatus 获取最终状态
func (crd *ProcessBase) getProcStatus(globalStatus, selfStatus, threatenStatus, autoWhiteSwitch int) (status int, at string) {
	at = helpers.DatetimeCurr()
	if globalStatus == enum_process.RelyStatusWhite {
		status = enum_process.StatusWhite
	} else if globalStatus == enum_process.RelyStatusBlack {
		status = enum_process.StatusBlack
	} else if selfStatus == enum_process.RelyStatusWhite {
		status = enum_process.StatusWhite
	} else if selfStatus == enum_process.RelyStatusBlack {
		status = enum_process.StatusBlack
	} else if selfStatus == enum_process.RelyStatusUnset {
		status = enum_process.StatusUnset
	} else if autoWhiteSwitch == enums.SwitchOn {
		if threatenStatus == enum_process.RelyStatusWhite {
			status = enum_process.StatusWhite
		} else if threatenStatus == enum_process.RelyStatusBlack {
			status = enum_process.StatusBlack
		} else {
			status = enum_process.StatusUnset
			at = helpers.TimestampToDatetime(helpers.AddTimestampByHour(-3))
		}
	} else {
		status = enum_process.StatusUnset
		at = helpers.TimestampToDatetime(helpers.AddTimestampByHour(-3))
	}

	return status, at
}
