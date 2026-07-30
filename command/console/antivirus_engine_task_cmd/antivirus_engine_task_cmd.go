package antivirus_engine_task_cmd

import "goskeleton/app/global/variable"

func CreateAntivirusEngineTaskFactory() *AntivirusEngineTaskCmd {
	return &AntivirusEngineTaskCmd{}
}

type AntivirusEngineTaskCmd struct {
}

// Start  探测任务启动
func (crd *AntivirusEngineTaskCmd) Start() {
	variable.ZapLog.Info("Start start")

	variable.ZapLog.Info("Start end")
}

// Overdue  探测任务逾期
func (crd *AntivirusEngineTaskCmd) Overdue() {
	variable.ZapLog.Info("overdue start")

	variable.ZapLog.Info("overdue end")
}
