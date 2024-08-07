package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/jobs/agent_job"
	"goskeleton/app/jobs/proc_job"
	"goskeleton/app/utils/helpers/date_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/proc_cmd"
	"testing"
)

func TestTaskFromProcess(t *testing.T) {
	proc_job.TaskFromProcess(1, "f7f8a4cbcad0b9885c9186eb0fbfcf78", "2024-07-17 10:45:01")
}

func TestDate(t *testing.T) {
	fmt.Println("现在时间戳：", date_helper.CurTimestamp())

	curDatetime := date_helper.CurDatetime()
	fmt.Println("现在时间：", curDatetime)

	timestamp := date_helper.DatetimeToTimestamp(curDatetime, variable.DateFormat)
	fmt.Println("时间戳：", timestamp)
	date := date_helper.TimestampToDatetime(timestamp)
	fmt.Println("时间戳还原时间：", date)
}

func TestOpenAutoWhite(t *testing.T) {
	agent_job.OpenAutoWhite(1, "881d508cabab2b8448c6d7ff8c941bf8")
}

func TestThreatenRepairStatus(t *testing.T) {
	proc_cmd.CreateCmdFactory().RepairStatusByThreaten()
}
func TestUploadProcess(t *testing.T) {
	jsonStr := "{\"id\":\"8cd274dc6542f0f5caa66044f39d7265\",\"org_id\":1,\"agent_id\":\"win2010_test_001\"" +
		",\"proc_hash\":\"2b67b9e85ca951de7d360bc2aceea291\",\"dir\":\"C:\\\\Program Files\\\\SplunkUniversalForwarder\\\\bin\\\\splunk-MonitorNoHandle.exe\"," +
		"\"status\":0,\"sync_status_at\":0,\"global_status\":0,\"self_status\":0,\"threaten_status\":0,\"threaten_level\":0," +
		"\"auto_white_switch\":0,\"grant_threaten_switch\":0,\"introduce\":\"SplunkMonNoHandle Control Program \"," +
		"\"copyright\":\"© Microsoft Corporation. All rights reserved. \",\"signature\":\"\",\"created_at\":\"2024-08-07 18:58:44\"}"
	proc_job.UploadProcess(jsonStr)
}
