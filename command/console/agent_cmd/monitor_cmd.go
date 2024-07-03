package agent_cmd

import "goskeleton/app/global/variable"

func CreateMonitorFactory() *MonitorCmd {
	return &MonitorCmd{}
}

type MonitorCmd struct {
}

// Delete  删除
func (crd *MonitorCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete")
}
