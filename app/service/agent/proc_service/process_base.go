package proc_service

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/utils/helpers/date_helper"
)

type ProcessBase struct {
}

func (crd *ProcessBase) getUpdatesForSetStatus(status int) map[string]interface{} {
	updates := map[string]interface{}{
		"status":         status,
		"sync_status_at": 0,
	}

	switch status {
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
	str = fmt.Sprintf(str, enums.IsNot, enum_process.ProcStatusUnset,
		enum_process.RelyStatusWhite, enum_process.ProcStatusWhite,
		enum_process.RelyStatusBlack, enum_process.ProcStatusBlack,
		enum_process.ProcStatusUnset)
	updates["status"] = gorm.Expr(str)
	return updates
}

func (crd *ProcessBase) getUpdatesForBatchWhite() map[string]interface{} {
	updates := map[string]interface{}{
		"self_status":    enum_process.RelyStatusWhite,
		"sync_status_at": 0,
		"status":         enum_process.ProcStatusWhite,
	}
	return updates
}

func (crd *ProcessBase) getUpdatesForAddGlobal(globalType int) map[string]interface{} {
	updates := map[string]interface{}{
		"sync_status_at": date_helper.CurDatetime(),
	}

	switch globalType {
	case enum_process.GlobalTypeWhite:
		updates["global_status"] = enum_process.RelyStatusWhite
		updates["status"] = enum_process.ProcStatusWhite
		break
	case enum_process.ProcStatusBlack:
		updates["global_status"] = enum_process.RelyStatusWhite
		updates["status"] = enum_process.ProcStatusWhite
		break
	}

	return updates
}
