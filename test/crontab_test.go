package test

import (
	"fmt"
	"goskeleton/app/service/agent/plug_service/plug_calc_service"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console"
	"goskeleton/command/console/org_cmd"
	"testing"
)

func TestCrontabStart(t *testing.T) {
	console.TestStart("AgentSafe")
}

func TestCrontabStart1(t *testing.T) {
	//动态封堵
	if err := plug_calc_service.CreatePlugCalcFactory().CalcByLogin(1, "win2010_test_001", "192.168.4.66"); err != nil {
		fmt.Println(err.Error())
	}
}

func TestCrontabCollectOrgYesterday(t *testing.T) {
	org_cmd.CreateCmdFactory().YesterdayStatistics()
}
