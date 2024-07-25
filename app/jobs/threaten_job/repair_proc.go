package threaten_job

import (
	"fmt"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func RepairProc(dataFrom, level int, jsonStr string) {
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("RepairProc dataFrom=%d, level=%d, jsonSrt=%s", dataFrom, level, jsonStr))
	ids := json_helper.ToSliceString(jsonStr)
	if len(ids) == 0 {
		return
	}

	if dataFrom == enum_threaten.DataFromProcess {
		proc_service.CreateProcessCurdFactory().ThreatenCallbackByHashes(ids, level)
	} else {
		proc_service.CreateDllCurdFactory().ThreatenCallbackByHashes(ids, level)
	}
}
