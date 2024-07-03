package agent_cmd

import "goskeleton/app/global/variable"

func CreateTaskFactory() *TaskCmd {
	return &TaskCmd{}
}

type TaskCmd struct {
}

// Delete  删除
func (crd *TaskCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete")
}
