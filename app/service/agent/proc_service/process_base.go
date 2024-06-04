package proc_service

import "goskeleton/app/global/enums/enum_process"

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
