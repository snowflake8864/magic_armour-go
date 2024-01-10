package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
	"time"
)

func TestDbMigrate(t *testing.T) {
	migrateFact := modelPlatform.CreateMigrateFactory()
	migrateFact.CreateTable()

	agentFact := modelPlatform.CreateAgentFactory()
	agentFact.CreateTable()
}

func TestInsertOne(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()
	var st = &modelPlatform.MigrateModel{
		Migration: "20240108000001_init_paltform111",
		CreatedAt: time.Now().Format(variable.DateFormat),
	}

	result := modelFact.InsertItem(st)
	fmt.Println(result)
}

func TestInsertMany(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()

	var st = &[]modelPlatform.MigrateModel{
		{
			Migration: "20240108000001_init_paltforma' or '2'='2",
			CreatedAt: time.Now().Format(variable.DateFormat),
		},
		{
			Migration: "20240108000001_init_paltforma' or '3'='3",
			CreatedAt: time.Now().Format(variable.DateFormat),
		},
	}
	result := modelFact.InsertItems(st)

	var list []modelPlatform.MigrateModel
	var s1 = modelPlatform.MigrateModel{
		Migration: "20240108000001_init_paltforma333333333",
		CreatedAt: time.Now().Format(variable.DateFormat),
	}
	var s2 = modelPlatform.MigrateModel{
		Migration: "20240108000001_init_paltform44444444",
		CreatedAt: time.Now().Format(variable.DateFormat),
	}
	list = append(list, s1, s2)
	result = modelFact.InsertItems(&list)
	fmt.Println(result)
}

func TestShowOne(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()

	item, err := modelFact.ShowOneItem()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(item)

}

func TestShowAll(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()

	items, err := modelFact.GetALLItems()
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range items {
		fmt.Println(key)
		fmt.Println(value.Migration)
	}
}

func TestShowAllMAP(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()

	items, err := modelFact.GetAllMap()
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range items {
		fmt.Println(key)
		fmt.Println(value.Migration)
	}
}

func TestToSql(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()
	_ = modelFact.TestToSQL()

}
