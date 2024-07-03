package agent_cmd

import "goskeleton/app/global/variable"

func CreateAccessFactory() *AccessCmd {
	return &AccessCmd{}
}

type AccessCmd struct {
}

// Clear  清空
func (crd *AccessCmd) Clear() {
	variable.ZapLog.Sugar().Debug("Clear")
}
