package agent_cmd

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/dir_models"
)

func CreateDirFactory() *DirCmd {
	return &DirCmd{}
}

type DirCmd struct {
}

// Clear  清空
func (crd *DirCmd) Clear() {
	variable.ZapLog.Sugar().Debug("Clear start")
	dir_models.CreateDirFactory().TruncateTable()

	variable.ZapLog.Sugar().Debug("Clear end")
}
