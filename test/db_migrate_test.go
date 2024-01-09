package test

import (
	"goskeleton/app/global/variable"
	model "goskeleton/app/model/platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
	"time"
)

func TestDbMigrate(t *testing.T) {
	modelFact := model.CreateMigrateFactory()
	modelFact.CreateTable()
	modelFact.AddOne()
}

func TestInsert(t *testing.T) {
	var st = &model.MigrateModel{
		Migration: "20240108000001_init_paltform111",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	result := variable.GormDbPlatform.Create(st)
	if result.RowsAffected < 0 {
		t.Error("新增失败，错误详情：", result.Error.Error())
	}
}
