package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	agent2 "goskeleton/app/http/validator/organization/agent"
	auth2 "goskeleton/app/http/validator/organization/auth"
	department2 "goskeleton/app/http/validator/organization/department"
	home2 "goskeleton/app/http/validator/organization/home"
	proc_global2 "goskeleton/app/http/validator/organization/proc_global"
	proc_sample2 "goskeleton/app/http/validator/organization/proc_sample"
	process2 "goskeleton/app/http/validator/organization/process"
	tamper2 "goskeleton/app/http/validator/organization/tamper"
	"goskeleton/app/http/validator/platform/threaten_push"
)

// OrganizationRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func OrganizationRegisterValidator() {

	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//安全认证
	key = consts.ValidatorPrefixPlatform + "AuthCheckUser"
	containers.Set(key, auth2.CheckUser{})
	key = consts.ValidatorPrefixPlatform + "AuthGetSliderCode"
	containers.Set(key, auth2.GetSliderCode{})
	key = consts.ValidatorPrefixPlatform + "AuthSendSms"
	containers.Set(key, auth2.SendSms{})
	key = consts.ValidatorPrefixPlatform + "AuthNormalLogin"
	containers.Set(key, auth2.NormalLogin{})
	key = consts.ValidatorPrefixPlatform + "AuthLogout"
	containers.Set(key, auth2.Logout{})
	key = consts.ValidatorPrefixPlatform + "AuthHsLogin"
	containers.Set(key, auth2.HsLogin{})

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeGetAgentStatistic"
	containers.Set(key, home2.GetAgentStatistic{})
	key = consts.ValidatorPrefixPlatform + "HomeGetAlarmStatistic"
	containers.Set(key, home2.GetAlarmStatistic{})

	//组织架构
	key = consts.ValidatorPrefixPlatform + "DepartmentTree"
	containers.Set(key, department2.Tree{})
	key = consts.ValidatorPrefixPlatform + "DepartmentAdd"
	containers.Set(key, department2.Add{})
	key = consts.ValidatorPrefixPlatform + "DepartmentEdit"
	containers.Set(key, department2.Edit{})
	key = consts.ValidatorPrefixPlatform + "DepartmentDelete"
	containers.Set(key, department2.Delete{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentManagePage"
	containers.Set(key, agent2.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "AgentManageDelete"
	containers.Set(key, agent2.ManageDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUninstall"
	containers.Set(key, agent2.ManageUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUpgrade"
	containers.Set(key, agent2.ManageUpgrade{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchDelete"
	containers.Set(key, agent2.ManageBatchDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUninstall"
	containers.Set(key, agent2.ManageBatchUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUpgrade"
	containers.Set(key, agent2.ManageBatchUpgrade{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchSetConfig"
	containers.Set(key, agent2.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicGet"
	containers.Set(key, agent2.BasicGet{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicSave"
	containers.Set(key, agent2.BasicSave{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigGet"
	containers.Set(key, agent2.ConfigGet{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigSave"
	containers.Set(key, agent2.ConfigSave{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorList"
	containers.Set(key, agent2.MonitorList{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorInfo"
	containers.Set(key, agent2.MonitorInfo{})
	key = consts.ValidatorPrefixPlatform + "AgentErrorPage"
	containers.Set(key, agent2.ErrorPage{})

	//进程
	key = consts.ValidatorPrefixPlatform + "ProcessManagePage"
	containers.Set(key, process2.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetSwitch"
	containers.Set(key, process2.ManageSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetMode"
	containers.Set(key, process2.ManageSetMode{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetConfig"
	containers.Set(key, process2.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetWhite"
	containers.Set(key, process2.ManageBatchSetWhite{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchRefresh"
	containers.Set(key, process2.ManageBatchRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchCollect"
	containers.Set(key, process2.ManageBatchCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchThreatenWhite"
	containers.Set(key, process2.ManageBatchThreatenWhite{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageAll"
	containers.Set(key, process2.AgentPageAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageRunning"
	containers.Set(key, process2.AgentPageRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetStatus"
	containers.Set(key, process2.AgentSetStatus{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCancel"
	containers.Set(key, process2.AgentCancel{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearAll"
	containers.Set(key, process2.AgentClearAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearRunning"
	containers.Set(key, process2.AgentClearRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRefresh"
	containers.Set(key, process2.AgentRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCollect"
	containers.Set(key, process2.AgentCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentExtractRunningDll"
	containers.Set(key, process2.AgentExtractRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentShowRunningDll"
	containers.Set(key, process2.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchWhite"
	containers.Set(key, process2.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRepeatPage"
	containers.Set(key, process2.AgentRepeatHash{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentExtractSample"
	containers.Set(key, proc_sample2.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetGlobal"
	containers.Set(key, proc_global2.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchGlobalWhite"
	containers.Set(key, proc_global2.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchThreatenWhite"
	containers.Set(key, threaten_push.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchThreatenBlack"
	containers.Set(key, threaten_push.BatchBlackFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalPage"
	containers.Set(key, proc_global2.PageFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalAdd"
	containers.Set(key, proc_global2.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalEdit"
	containers.Set(key, proc_global2.Edit{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalDelete"
	containers.Set(key, proc_global2.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalClear"
	containers.Set(key, proc_global2.ClearFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchPage"
	containers.Set(key, process2.SearchPage{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalWhite"
	containers.Set(key, proc_global2.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalBlack"
	containers.Set(key, proc_global2.BatchBlackFromProcess{})

	//进程和模块样本
	key = consts.ValidatorPrefixPlatform + "ProcSamplePage"
	containers.Set(key, proc_sample2.Page{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDelete"
	containers.Set(key, proc_sample2.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleClear"
	containers.Set(key, proc_sample2.Clear{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDownload"
	containers.Set(key, proc_sample2.Download{})

	//防篡改
	key = consts.ValidatorPrefixPlatform + "TamperManagePage"
	containers.Set(key, tamper2.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "TamperManageSetSwitch"
	containers.Set(key, tamper2.ManageSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperManageSetMode"
	containers.Set(key, tamper2.ManageSetMode{})
	key = consts.ValidatorPrefixPlatform + "TamperManageBatchSetConfig"
	containers.Set(key, tamper2.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyPage"
	containers.Set(key, tamper2.PolicyPage{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicySetSwitch"
	containers.Set(key, tamper2.PolicySetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyAddFolder"
	containers.Set(key, tamper2.PolicyAddFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyAddFile"
	containers.Set(key, tamper2.PolicyAddFile{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyEditFolder"
	containers.Set(key, tamper2.PolicyEditFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyEditFile"
	containers.Set(key, tamper2.PolicyEditFile{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyDelete"
	containers.Set(key, tamper2.PolicyDelete{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalPage"
	containers.Set(key, tamper2.GlobalPage{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalSetSwitch"
	containers.Set(key, tamper2.GlobalSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalAddFolder"
	containers.Set(key, tamper2.GlobalAddFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalAddFile"
	containers.Set(key, tamper2.GlobalAddFile{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalEditFolder"
	containers.Set(key, tamper2.GlobalEditFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalEditFile"
	containers.Set(key, tamper2.GlobalEditFile{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalDelete"
	containers.Set(key, tamper2.GlobalDelete{})
}
