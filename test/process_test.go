package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/jobs/proc_job"
	"goskeleton/app/utils/helpers/date_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
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
