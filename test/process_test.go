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

	timestamp := date_helper.DatetimeToTimestamp(curDatetime, variable.DateFormatNow)
	fmt.Println("时间戳：", timestamp)
	date := date_helper.TimestampToDatetime(timestamp)
	fmt.Println("时间戳还原时间：", date)

	fmt.Println("--------------------------")
	dateDay := "2025-02-14"
	fmt.Println("现在时间：", dateDay)

	timestampDay := date_helper.DatetimeToTimestamp(dateDay, variable.DateFormatToday)
	fmt.Println("时间戳：", timestampDay)
	dateDay1 := date_helper.TimestampToDatetime(timestampDay)
	fmt.Println("时间戳还原时间：", dateDay1)
}

func TestOpenAutoWhite(t *testing.T) {
	agent_job.OpenAutoWhite(1, "881d508cabab2b8448c6d7ff8c941bf8")
}

func TestThreatenRepairStatus(t *testing.T) {
	proc_cmd.CreateCmdFactory().RepairStatusByThreaten()
}
