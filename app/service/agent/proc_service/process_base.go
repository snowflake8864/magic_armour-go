package proc_service

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/utils/helpers/date_helper"
)

type ProcessBase struct {
}

// getUpdatesForSetStatus 根据设置状态动作得到修改值
func (crd *ProcessBase) getUpdatesForSetStatus(procStatus int) map[string]interface{} {
	updates := map[string]interface{}{
		"status":         procStatus,
		"sync_status_at": 0,
	}

	switch procStatus {
	case enum_process.ProcStatusUnknown:
		updates["self_status"] = enum_process.RelyStatusUnknown
		break
	case enum_process.ProcStatusWhite:
		updates["self_status"] = enum_process.RelyStatusWhite
		break
	case enum_process.ProcStatusBlack:
		updates["self_status"] = enum_process.RelyStatusBlack
		break
	}

	return updates
}

// getUpdatesForCancel 根据取消动作得到修改值
func (crd *ProcessBase) getUpdatesForCancel() map[string]interface{} {
	updates := map[string]interface{}{
		"self_status":    enum_process.RelyStatusUnset,
		"sync_status_at": 0,
	}

	str := "(CASE " +
		"WHEN auto_white_switch=%d THEN %d " +
		"WHEN threaten_status=%d THEN %d " +
		"WHEN threaten_status=%d THEN %d " +
		"ELSE %d " +
		"END)"
	str = fmt.Sprintf(str,
		enums.IsNot, enum_process.ProcStatusUnset,
		enum_process.RelyStatusWhite, enum_process.ProcStatusWhite,
		enum_process.RelyStatusBlack, enum_process.ProcStatusBlack,
		enum_process.ProcStatusUnset)
	updates["status"] = gorm.Expr(str)
	return updates
}

// getUpdatesForBatchWhite 根据批量白名单动作得到修改值
func (crd *ProcessBase) getUpdatesForBatchWhite() map[string]interface{} {
	updates := map[string]interface{}{
		"self_status":    enum_process.RelyStatusWhite,
		"sync_status_at": 0,
		"status":         enum_process.ProcStatusWhite,
	}
	return updates
}

// getUpdatesForBatchWhite 根据加入全局名单动作得到修改值
func (crd *ProcessBase) getUpdatesForAddGlobal(globalType int) map[string]interface{} {
	updates := map[string]interface{}{
		"sync_status_at": date_helper.CurTimestamp(),
	}

	switch globalType {
	case enum_process.GlobalTypeWhite:
		updates["global_status"] = enum_process.RelyStatusWhite
		updates["status"] = enum_process.ProcStatusWhite
		break
	case enum_process.GlobalTypeBlack:
		updates["global_status"] = enum_process.RelyStatusBlack
		updates["status"] = enum_process.ProcStatusBlack
		break
	}

	return updates
}

// getUpdatesForDeleteGlobal 根据删除全局名单动作得到修改值
func (crd *ProcessBase) getUpdatesForDeleteGlobal() map[string]interface{} {
	updates := map[string]interface{}{
		"sync_status_at": date_helper.CurTimestamp(),
		"global_status":  enum_process.RelyStatusUnset,
	}
	str := "(CASE " +
		"WHEN self_status=%d THEN %d " +
		"WHEN self_status=%d THEN %d " +
		"WHEN self_status=%d THEN %d " +
		"WHEN auto_white_switch=%d THEN %d " +
		"WHEN threaten_status=%d THEN %d " +
		"WHEN threaten_status=%d THEN %d " +
		"ELSE %d " +
		"END)"
	str = fmt.Sprintf(str,
		enum_process.RelyStatusWhite, enum_process.ProcStatusWhite,
		enum_process.RelyStatusBlack, enum_process.ProcStatusBlack,
		enum_process.RelyStatusUnknown, enum_process.ProcStatusUnknown,
		enums.IsNot, enum_process.ProcStatusUnset,
		enum_process.RelyStatusWhite, enum_process.ProcStatusWhite,
		enum_process.RelyStatusBlack, enum_process.ProcStatusBlack,
		enum_process.ProcStatusUnset)
	updates["status"] = gorm.Expr(str)

	return updates
}

// getUpdatesForAutoWhiteSwitch 根据自动匹配白名单开关动作得到修改值
func (crd *ProcessBase) getUpdatesForAutoWhiteSwitch(openSwitch int) map[string]interface{} {
	updates := map[string]interface{}{
		"auto_white_switch": openSwitch,
	}

	if openSwitch == enums.SwitchOn {
		statusStr := "(CASE " +
			"WHEN status=%d AND threaten_status=%d THEN %d " +
			"WHEN status=%d AND threaten_status=%d THEN %d " +
			"ELSE status " +
			"END)"
		statusStr = fmt.Sprintf(statusStr,
			enum_process.ProcStatusUnset, enum_process.RelyStatusWhite, enum_process.ProcStatusWhite,
			enum_process.ProcStatusUnset, enum_process.RelyStatusBlack, enum_process.ProcStatusBlack,
		)

		updates["status"] = gorm.Expr(statusStr)
		updates["sync_status_at"] = date_helper.CurTimestamp()
	}
	return updates
}

// getUpdatesForThreatenCallback 根据威胁情报回调动作得到修改值
func (crd *ProcessBase) getUpdatesForThreatenCallback(threatenLevel int) map[string]interface{} {
	updates := map[string]interface{}{}

	statusStr := "IF(auto_white_switch=%d AND status=%d, %d, status)"
	syncStatusAt := "IF(auto_white_switch=%d AND status=%d, %d, sync_status_at)"
	switch threatenLevel {
	case enum_threaten.LevelWhite:
		updates["threaten_status"] = enum_process.RelyStatusWhite
		statusStr = fmt.Sprintf(statusStr, enums.IsYes, enum_process.ProcStatusUnset, enum_process.ProcStatusWhite)
		syncStatusAt = fmt.Sprintf(syncStatusAt, enums.IsYes, enum_process.ProcStatusUnset, date_helper.CurTimestamp())
		break
	case enum_threaten.LevelBlack:
		updates["threaten_status"] = enum_process.RelyStatusBlack
		statusStr = fmt.Sprintf(statusStr, enums.IsYes, enum_process.ProcStatusUnset, enum_process.ProcStatusBlack)
		syncStatusAt = fmt.Sprintf(syncStatusAt, enums.IsYes, enum_process.ProcStatusUnset, date_helper.CurTimestamp())
		break
	default:
		return nil
	}
	updates["status"] = gorm.Expr(statusStr)
	updates["sync_status_at"] = gorm.Expr(syncStatusAt)
	return updates
}

// getUpdatesForThreatenRepair 根据威胁情报修复动作得到修改值
func (crd *ProcessBase) getUpdatesForThreatenRepair(threatenStatus int) map[string]interface{} {
	updates := map[string]interface{}{}

	switch threatenStatus {
	case enum_process.RelyStatusWhite:
		updates["status"] = enum_process.ProcStatusWhite
		updates["sync_status_at"] = date_helper.CurTimestamp()
		break
	case enum_process.ProcStatusBlack:
		updates["status"] = enum_process.ProcStatusBlack
		updates["sync_status_at"] = date_helper.CurTimestamp()
		break
	}

	return updates
}
