package service_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *ServiceCmd {
	return &ServiceCmd{}
}

type ServiceCmd struct {
}

// Monitor  服务器状态监控
func (crd *ServiceCmd) Monitor() {
	variable.ZapLog.Sugar().Debug("Monitor")
}
