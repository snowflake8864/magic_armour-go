package test

import (
	"fmt"
	"goskeleton/app/service/agent/task_service"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestCloseTask(t *testing.T) {
	unfinishedTasks := []int{1, 2, 3}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(unfinishedTasks), cap(unfinishedTasks), unfinishedTasks)

	unfinishedTasks = []int{10, 11}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(unfinishedTasks), cap(unfinishedTasks), unfinishedTasks)

	task_service.CreateCurdFactory().CloseTasks("20240108000001_init_paltforma' or '2'='2", unfinishedTasks, "168.10.1.1")
}
