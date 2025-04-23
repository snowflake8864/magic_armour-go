package service_cmd

import (
	"bytes"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/service/platform/service_monitor_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"os/exec"
	"strings"
)

func CreateCmdFactory() *ServiceCmd {
	return &ServiceCmd{}
}

type ServiceCmd struct {
}

// Delete 删除监控数据
func (crd *ServiceCmd) Delete() {
	variable.ZapLog.Info("Delete start")
	monitorDay, _, _, _ := date_helper.CalcByDay(-90)
	server_monitor_models.CreateServerMonitorFactory().DeleteByConditions(server_monitor_models.SearchConditions{
		MonitorDayLte: monitorDay,
	})

	variable.ZapLog.Info("Delete end")
}

// Monitor  服务器状态监控
func (crd *ServiceCmd) Monitor() {

	serviceCode := cur_app.GetServiceCode()
	cmdInit := cur_app.ShellInit()
	cmdServiceMonitor := cur_app.ShellServiceMonitor()
	variable.ZapLog.Info(fmt.Sprintf("Monitor code:%s, cmd:%s %s", serviceCode, cmdInit, cmdServiceMonitor))

	cmd := exec.Command(cmdInit, cmdServiceMonitor)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	if err != nil {
		variable.ZapLog.Warn(fmt.Sprintf("exec command Run failed: %s", err.Error()))
		return
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if string_helper.IsEmpty(outStr) || !string_helper.IsEmpty(errStr) {
		variable.ZapLog.Warn(fmt.Sprintf("exec command output empty out:\n%s\nerr:\n%s\n", outStr, errStr))
		return
	}

	//监控日期
	monitorDay, _, _, _ := date_helper.CalcByDay(0)

	model := server_monitor_models.ServerMonitorModel{
		ServerCode: serviceCode,
		MonitorDay: monitorDay,
		CreatedAt:  date_helper.CurDatetime(),
		CpuUsage:   "0.00",
		MemUsage:   "0.00",
		HdUsage:    "0.00",
	}
	outList := strings.Split(outStr, "\n")
	for _, line := range outList {
		arr := strings.Split(line, " ")
		if n := len(arr); n >= 2 {
			if strings.ToLower(arr[0]) == "cpu" {
				model.CpuUsage = arr[1]
			} else if strings.ToLower(arr[0]) == "mem" {
				model.MemUsage = arr[1]
			} else if strings.ToLower(arr[0]) == "disk" {
				model.HdUsage = arr[1]
			}
		}
	}
	if !service_monitor_service.CreateCurdFactory().Add(&model) {
		variable.ZapLog.Warn("Add failed")
	}

	variable.ZapLog.Info(fmt.Sprintf("Monitor code:%s, cmd:%s %s", serviceCode, cmdInit, cmdServiceMonitor))
}
