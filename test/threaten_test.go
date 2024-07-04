package test

import (
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/threaten_cmd"
	"testing"
)

func TestThreatenSync(t *testing.T) {
	threaten_cmd.CreateCmdFactory().SyncUnknownForCache()
}
