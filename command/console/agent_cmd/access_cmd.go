package agent_cmd

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_access_models"
)

func CreateAccessFactory() *AccessCmd {
	return &AccessCmd{}
}

type AccessCmd struct {
}

// Clear  清空
func (crd *AccessCmd) Clear() {
	variable.ZapLog.Sugar().Debug("Clear start")
	agent_access_models.CreateAccessLogFactory().TruncateTable()

	variable.ZapLog.Sugar().Debug("Clear end")
}
