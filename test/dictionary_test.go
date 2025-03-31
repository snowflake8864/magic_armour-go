package test

import (
	"fmt"
	"goskeleton/app/model/platform/dictionary_models"
	"goskeleton/app/service/platform/dictionary_service"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestGetItems(t *testing.T) {
	fact := dictionary_service.CreateCurdFactory()
	fmt.Println("------------------------------")
	fmt.Println(fact.GetMapByModule("admin"))
	fmt.Println("------------------------------")
	fmt.Println(fact.GetMapByModule("admin_login_error_max_test"))

	fmt.Println("------------------------------")
	l1 := fact.GetAll()
	for _, obj := range l1 {
		fmt.Println("IdSingle:" + obj.Id + " VAL:" + obj.Val)
	}
}

func TestUpdate(t *testing.T) {
	fact := dictionary_service.CreateCurdFactory()
	fmt.Println("------------------------------")
	fmt.Println(fact.UpdateValById("agent_save_access_log", "1"))
	fmt.Println(fact.UpdateValById("agent_debug_ids", "win2010_test_001"))
}

func TestRebuildCache(t *testing.T) {
	dictionary_models.CreateDictionaryFactory().InitData()

	fact := dictionary_service.CreateCacheFactory()
	fmt.Println("------------------------------")
	fmt.Println(fact.Rebuild())
}

func TestGetMapByIds(t *testing.T) {
	fact := dictionary_service.CreateCacheFactory()
	fmt.Println("------------------------------")
	fmt.Println(fact.GetMapByIds([]string{"admin_login_lock_duration", "ffffff", "admin_pw_length_max", "agent_test"}))
}
