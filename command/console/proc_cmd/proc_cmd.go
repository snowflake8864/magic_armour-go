package proc_cmd

import "goskeleton/app/global/variable"

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
