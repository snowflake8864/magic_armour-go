package test

import (
	model "goskeleton/app/model/platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestDbMigrate(t *testing.T) {
	modelFact := model.CreateAgentFactory()
	modelFact.CreateTable()
}
