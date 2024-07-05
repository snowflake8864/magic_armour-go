package huanshi_cmd

import (
	"goskeleton/app/global/variable"
)

func CreateCmdFactory() *HuanshiCmd {
	return &HuanshiCmd{}
}

type HuanshiCmd struct {
}

// Push  推送
func (crd *HuanshiCmd) Push() {
	variable.ZapLog.Sugar().Debug("push start")

	variable.ZapLog.Sugar().Debug("push end")
}
