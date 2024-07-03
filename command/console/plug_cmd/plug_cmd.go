package plug_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *PlugCmd {
	return &PlugCmd{}
}

type PlugCmd struct {
}

// DynamicExpired  动态封堵已过期的
func (crd *PlugCmd) DynamicExpired() {
	variable.ZapLog.Sugar().Debug("DynamicExpired")
}
