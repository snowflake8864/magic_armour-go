package proc_job

import (
	"fmt"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func TaskFromProc(orgId int64, agentId string, at string, dataFrom int) {
	//variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s, at=%s, dataFrom=%d", orgId, agentId, at, dataFrom))
	//任务
	var taskTypes []int

	//条件
	condition := proc_models.SearchConditions{
		OrgId:           orgId,
		AgentId:         agentId,
		SyncStatusAtGte: date_helper.DatetimeToTimestamp(at, variable.DateFormatNow),
	}

	//进程已开启
	if (dataFrom & enum_process.DataFromProcess) > 0 {
		if total := proc_models.CreateProcessFactory().Count(condition); total > 0 {
			taskTypes = append(taskTypes, enum_task.TypeProcessWhiteDownload, enum_task.TypeProcessBlackDownload)
		}
	}

	//模块已开启
	if (dataFrom & enum_process.DataFromDll) > 0 {
		if total := proc_models.CreateDllFactory().Count(condition); total > 0 {
			taskTypes = append(taskTypes, enum_task.TypeDllWhiteDownload, enum_task.TypeDllBlackDownload)
		}
	}

	if len(taskTypes) > 0 {
		variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s, at=%s, dataFrom=%d", orgId, agentId, at, dataFrom))

		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, taskTypes, "同步状态")
	}
}
