package agent_cmd

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateTaskFactory() *TaskCmd {
	return &TaskCmd{}
}

type TaskCmd struct {
}

// Delete  删除
func (crd *TaskCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete start")
	_, _, dayEnd, _ := date_helper.CalcByDay(-60)
	task_models.CreateTaskFactory().DeleteByConditions(task_models.SearchConditions{
		CreateAtLte: dayEnd,
	})
	variable.ZapLog.Sugar().Debug("Delete end")
}
