package test

import (
	"fmt"
	"goskeleton/app/jobs/alarm_job"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/service/agent/plug_service/plug_calc_service"
	"goskeleton/app/utils/helpers/json_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console"
	"goskeleton/command/console/org_cmd"
	"testing"
)

func TestCrontabStart(t *testing.T) {
	console.HttpStart("ThreatenCacheUnknownSync")
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

func TestStatInStartAtAndEndAt(t *testing.T) {
	//org_cmd.CreateCmdFactory().YesterdayStatistics()
	param := alarm_stat_models.StatDailyModel{
		OrgId:     1,
		AgentId:   "win20108cabab2b8448c6d7ff8c94ggg",
		StatDay:   "2026-02-25",
		StatMonth: 202602,
		CreatedAt: "2026-02-25 23:59:59",
	}
	alarm_job.StatInStartAtAndEndAt(json_helper.Encode(param), "2026-02-25 00:00:00", "2026-02-25 23:59:59")
}
