package test

import (
	"encoding/json"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/migrate_models"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/database"
	"testing"
	"time"
)

func TestDbMigrateType(t *testing.T) {

	fmt.Println(enums.Get1stIsStrReturnSliceOfStr(enum_dictionary.ModulePlatform, enum_dictionary.EnumSlice))

	fmt.Println(enums.Get1stIsStrAnd2ndIsIntReturnStr("ip_type", 1, enum_plug.EnumMap))
	fmt.Println(enums.Get1stIsStrAnd2ndIsIntReturnStr("ip_type", 3, enum_plug.EnumMap))
	fmt.Println(enums.Get1stIsStrAnd2ndIsIntReturnStr("ip_type_test", 3, enum_plug.EnumMap))

	c := enum_plug.PlugPortConfig{
		Switch:   1,
		PlugTime: 40,
	}
	c.Format()
	jsonPlugPort, _ := json.Marshal(c)
	fmt.Println(string(jsonPlugPort))

	d := enum_plug.PlugLoginConfig{
		Switch:   1,
		Type:     2,
		PlugTime: 20,
	}
	d.Format()
	fmt.Println(d)
}

func TestDbMigrate(t *testing.T) {

	database.TestStart("init")
	database.TestStart("menu")
	database.TestStart("dictionary")
}

func TestInsertMany(t *testing.T) {
	modelFact := migrate_models.CreateMigrateFactory()

	var st = &[]migrate_models.MigrateModel{
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

func TestShowAllMAP(t *testing.T) {
	modelFact := migrate_models.CreateMigrateFactory()

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
