package kafka_go

type MessageStruct struct {
	Type    string `json:"type"`
	Sign    string `json:"sign"`
	OrgId   int64  `json:"oId"`
	AgentId string `json:"aId"`
	AgentOS string `json:"aOS"`
	Data    string `json:"data"`
}

const (
	TopicCommon   string = "TopicCommon"
	TopicProcess  string = "TopicProcess"
	TopicDll      string = "TopicDll"
	TopicThreaten string = "TopicThreaten"
	TopicAlarm    string = "TopicAlarm"
	TopicSyslog   string = "TopicSyslog"
	TopicCrontab  string = "TopicCrontab"

	TypeCommonAgentDelete        string = "common_agent_delete"          //类型：终端删除
	TypeCommonAgentPutMonitor    string = "common_agent_put_monitor"     //类型：终端上传监控数据
	TypeCommonAgentOpenAutoWhite string = "common_agent_open_auto_white" //类型：终端开启自动白名单
	TypeCommonOrgClean           string = "common_org_clean"             //类型：企业清空业务数据
	TypeCommonEmailAlarmCycle    string = "common_email_alarm_cycle"     //类型：邮件安全告警通知(周期)

	TypeProcessBatchPutMaster  string = "process_batch_put_master"  //类型：上传主进程
	TypeProcessBatchPutRunning string = "process_batch_put_running" //类型：上传运行中进程
	TypeDllBatchPutMaster      string = "dll_batch_put_master"      //类型：上传主模块
	TypeDllBatchPutRunning     string = "dll_batch_put_running"     //类型：上传运行中模块

	TypeAlarmPutProcess   string = "alarm_put_process"   //类型：进程告警上传
	TypeAlarmPutDll       string = "alarm_put_dll"       //类型：模块告警上传
	TypeAlarmPutTamper    string = "alarm_put_tamper"    //类型：防篡改告警上传
	TypeAlarmPutBlackmail string = "alarm_put_blackmail" //类型：勒索告警上传
	TypeAlarmPutUsb       string = "alarm_put_usb"       //类型：外设告警上传
	TypeAlarmPutSelf      string = "alarm_put_self"      //类型：自保告警上传
	TypeAlarmPutOutreach  string = "alarm_put_outreach"  //类型：外联告警上传
	TypeAlarmPutPort      string = "alarm_put_port"      //类型：端口虚开告警上传

	TypeSyslogPutDns          string = "syslog_put_dns"           //类型：通讯日志-DNS上传
	TypeSyslogPutInner        string = "syslog_put_inner"         //类型：通讯日志-内联上传
	TypeSyslogPutOuter        string = "syslog_put_outer"         //类型：通讯日志-外联上传
	TypeSyslogPutProcess      string = "syslog_put_process"       //类型：通讯日志-进程启动上传
	TypeSyslogPutLoginSuccess string = "syslog_put_login_success" //类型：通讯日志-登录日志上传(成功)
	TypeSyslogPutLoginFail    string = "syslog_put_login_fail"    //类型：通讯日志-登录日志上传(失败)

	TypeThreatenSyncRewrite string = "threaten_sync_rewrite" //类型：威胁情报同步后重写
	TypeThreatenImport      string = "threaten_import"       //类型：威胁情报导入

	TypeCrontabAlarmStat          string = "crontab_alarm_stat"            //类型：定时任务-告警统计
	TypeCrontabAgentSafe          string = "crontab_agent_safe"            //类型：定时任务-终端安全状态
	TypeCrontabPlugDisable        string = "crontab_plug_disable"          //类型：定时任务-终端封堵解除
	TypeCrontabTaskFormProcess    string = "crontab_task_form_process"     //类型：定时任务-从进程新增任务
	TypeCrontabTaskFormDll        string = "crontab_task_form_dll"         //类型：定时任务-从模块新增任务
	TypeCrontabThreatenRepairProc string = "crontab_threaten_repair_proc " //类型：定时任务-威胁情报修复进程或模块
)
