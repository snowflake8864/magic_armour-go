package proc_job

import (
	"fmt"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func TaskFromProcess(orgId int64, agentId string, at string) {
	condition := proc_models.SearchConditions{
		OrgId:           orgId,
		AgentId:         agentId,
		SyncStatusAtGte: date_helper.DatetimeToTimestamp(at, variable.DateFormat),
	}
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("TaskFromProcess orgId=%d, agentId=%s, at=%s, SyncStatusAtGte=%d", orgId, agentId, at, condition.SyncStatusAtGte))

	if total := proc_models.CreateProcessFactory().Count(condition); total > 0 {
		types := []int{
			enum_task.TypeProcessWhiteDownload,
			enum_task.TypeProcessBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "同步状态")
	}
}

func TaskFromDll(orgId int64, agentId string, at string) {
	condition := proc_models.SearchConditions{
		OrgId:           orgId,
		AgentId:         agentId,
		SyncStatusAtGte: date_helper.DatetimeToTimestamp(at, variable.DateFormat),
	}
	variable.ZapLog.Sugar().Debug(fmt.Sprintf("TaskFromDll orgId=%d, agentId=%s, at=%s, SyncStatusAtGte=%d", orgId, agentId, at, condition.SyncStatusAtGte))

	if total := proc_models.CreateDllFactory().Count(condition); total > 0 {
		types := []int{
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "同步状态")
	}
}
