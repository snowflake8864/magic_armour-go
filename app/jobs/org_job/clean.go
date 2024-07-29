package org_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_error_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/syslog/syslog_models"
)

func Clean(orgId int64) {
	variable.ZapLog.Sugar().Info(fmt.Sprintf("orgId=%d", orgId))
	//终端
	agent_models.CreateAgentFactory().DeleteByOrg(orgId)
	agent_models.CreateExtendFactory().DeleteByOrg(orgId)
	//终端资源
	agent_monitor_models.CreateMonitorFactory().DeleteByOrg(orgId)
	//任务
	task_models.CreateTaskFactory().DeleteByOrg(orgId)
	//错误文件
	agent_error_models.CreateErrorFactory().DeleteByConditions(agent_error_models.SearchConditions{
		OrgId: orgId,
	})

	//进程或模块
	proc_models.CreateProcessFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId: orgId,
	})
	proc_models.CreateDllFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId: orgId,
	})
	proc_models.CreateProcSampleFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId: orgId,
	})
	proc_models.CreateProcGlobalFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId: orgId,
	})

	//勒索
	blackmail_models.CreatePolicyFactory().DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: orgId,
	})
	blackmail_models.CreateGlobalFactory().DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: orgId,
	})
	blackmail_models.CreateGroupFactory().DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: orgId,
	})

	//防篡改
	tamper_models.CreatePolicyFactory().DeleteByConditions(tamper_models.SearchConditions{
		OrgId: orgId,
	})
	tamper_models.CreateGlobalFactory().DeleteByConditions(tamper_models.SearchConditions{
		OrgId: orgId,
	})

	//端口
	port_models.CreateVirtualPolicyFactory().DeleteByConditions(port_models.SearchConditions{
		OrgId: orgId,
	})
	port_models.CreateVirtualGlobalFactory().DeleteByConditions(port_models.SearchConditions{
		OrgId: orgId,
	})
	port_models.CreateWorkingFactory().DeleteByConditions(port_models.SearchConditions{
		OrgId: orgId,
	})

	//封堵
	plug_models.CreateDynamicFactory().DeleteByConditions(plug_models.SearchConditions{
		OrgId: orgId,
	})
	plug_models.CreateIpPolicyFactory().DeleteByConditions(plug_models.SearchConditions{
		OrgId: orgId,
	})
	plug_models.CreateIpGlobalFactory().DeleteByConditions(plug_models.SearchConditions{
		OrgId: orgId,
	})

	//外设
	usb_models.CreateUsbFactory().DeleteByConditions(usb_models.SearchConditions{
		OrgId: orgId,
	})
	usb_models.CreateGlobalFactory().DeleteByConditions(usb_models.SearchConditions{
		OrgId: orgId,
	})

	//外联
	outreach_models.CreateGlobalFactory().DeleteByConditions(outreach_models.SearchConditions{
		OrgId: orgId,
	})

	//告警
	alarm_models.CreateProcessFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateDllFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateBlackmailFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateTamperFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreatePortFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreatePlugFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateUsbFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateOutreachFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateLoginFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})
	alarm_models.CreateSelfFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId: orgId,
	})

	//通讯日志
	syslog_models.CreateDnsFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId: orgId,
	})
	syslog_models.CreateOuterFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId: orgId,
	})
	syslog_models.CreateInnerFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId: orgId,
	})
	syslog_models.CreateProcessFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId: orgId,
	})
	syslog_models.CreateLoginFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId: orgId,
	})
}
