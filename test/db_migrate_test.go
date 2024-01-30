package test

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/database"
	"testing"
	"time"
)

func TestDbMigrateType(t *testing.T) {
	val := enums.GETFirstIsStringReturnSliceOfString(enums.DictionaryModuleAdmin, enums.EnumDictionarySlice)
	fmt.Println(val)
}

func TestDbMigrate(t *testing.T) {

	//database.TestStart("init")
	database.TestStart("menu")
	database.TestStart("dictionary")
	//_ = modelAgent.CreateOutreachGlobalFactory().InitData()
	//modelPlatform.CreateRoleFactory().InitData()
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
	result := modelFact.InsertBatch(st)
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

	items, err := modelFact.GetAllMaps()
	if items == nil {
		fmt.Println("items is nil")
		fmt.Println(err)
		return
	}
	if err != nil {
		fmt.Println(items)
		fmt.Println(err)
	}

	for key, value := range items {
		fmt.Println(key)
		fmt.Println(value)
	}
}

func TestToSql(t *testing.T) {
	modelFact := modelPlatform.CreateMigrateFactory()
	_ = modelFact.TestToSQL()

}
