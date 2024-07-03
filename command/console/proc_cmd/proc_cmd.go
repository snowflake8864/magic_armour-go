package proc_cmd

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
)

func CreateCmdFactory() *ProcCmd {
	return &ProcCmd{}
}

type ProcCmd struct {
}

// AddFromProcessForTask  任务：从进程添加
func (crd *ProcCmd) AddFromProcessForTask() {
	variable.ZapLog.Sugar().Debug("AddFromProcessForTask")
}

// AddFromDllForTask  任务：从模块添加
func (crd *ProcCmd) AddFromDllForTask() {
	variable.ZapLog.Sugar().Debug("AddFromDllForCache")
}

// ClearRunning  清空运行中数据
func (crd *ProcCmd) ClearRunning() {
	variable.ZapLog.Sugar().Debug("ClearRunning start")
	proc_models.CreateProcessRunningFactory().TruncateTable()
	proc_models.CreateDllRunningFactory().TruncateTable()
	variable.ZapLog.Sugar().Debug("ClearRunning end")
}
