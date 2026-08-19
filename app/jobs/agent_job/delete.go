package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_error_models"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/agent_online_models"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/syslog/syslog_models"
)

func Delete(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))
	//终端
	agent_models.CreateAgentFactory().DeleteByOrgAndAgent(orgId, agentId)
	agent_models.CreateExtendFactory().DeleteByOrgAndAgent(orgId, agentId)
	//终端资源
	agent_monitor_models.CreateMonitorFactory().DeleteByOrgAndAgent(orgId, agentId)
	//任务
	task_models.CreateTaskFactory().DeleteByOrgAndAgent(orgId, agentId)
	//错误文件
	agent_error_models.CreateErrorFactory().DeleteByConditions(agent_error_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//进程或模块
	proc_models.CreateProcessFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	proc_models.CreateDllFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	proc_models.CreateProcSampleFactory().DeleteByConditions(proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//勒索
	blackmail_models.CreatePolicyFactory().DeleteByConditions(blackmail_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//防篡改
	tamper_models.CreatePolicyFactory().DeleteByConditions(tamper_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//端口
	port_models.CreateVirtualPolicyFactory().DeleteByConditions(port_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	port_models.CreateWorkingFactory().DeleteByConditions(port_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//封堵
	plug_models.CreateDynamicFactory().DeleteByConditions(plug_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	plug_models.CreateIpPolicyFactory().DeleteByConditions(plug_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//外设
	usb_models.CreateUsbFactory().DeleteByConditions(usb_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//告警
	alarm_stat_models.CreateStatDailyFactory().DeleteByConditions(alarm_stat_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateProcessFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateDllFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateBlackmailFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateTamperFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreatePortFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreatePlugFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateUsbFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateOutreachFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateLoginFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	alarm_models.CreateSelfFactory().DeleteByConditions(alarm_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//通讯日志
	syslog_models.CreateDnsFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	syslog_models.CreateOuterFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	syslog_models.CreateInnerFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	syslog_models.CreateProcessFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
	syslog_models.CreateLoginFactory().DeleteByConditions(syslog_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})

	//上线日志
	agent_online_models.CreateOnlineLogFactory().DeleteByConditions(agent_online_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	})
}
