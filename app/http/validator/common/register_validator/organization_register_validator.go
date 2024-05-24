package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	agent3 "goskeleton/app/http/validator/organization/agent"
	auth3 "goskeleton/app/http/validator/organization/auth"
	department3 "goskeleton/app/http/validator/organization/department"
	home3 "goskeleton/app/http/validator/organization/home"
	procGlobal3 "goskeleton/app/http/validator/organization/proc_global"
	procSample3 "goskeleton/app/http/validator/organization/proc_sample"
	process3 "goskeleton/app/http/validator/organization/process"
	tamper3 "goskeleton/app/http/validator/organization/tamper"
)

// OrganizationRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func OrganizationRegisterValidator() {

	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//安全认证
	key = consts.ValidatorPrefixOrg + "AuthCheckUser"
	containers.Set(key, auth3.CheckUser{})
	key = consts.ValidatorPrefixOrg + "AuthGetSliderCode"
	containers.Set(key, auth3.GetSliderCode{})
	key = consts.ValidatorPrefixOrg + "AuthSendSms"
	containers.Set(key, auth3.SendSms{})
	key = consts.ValidatorPrefixOrg + "AuthNormalLogin"
	containers.Set(key, auth3.NormalLogin{})
	key = consts.ValidatorPrefixOrg + "AuthHsLogin"
	containers.Set(key, auth3.HsLogin{})

	//首页
	key = consts.ValidatorPrefixOrg + "HomeGetAgentStatistic"
	containers.Set(key, home3.GetAgentStatistic{})
	key = consts.ValidatorPrefixOrg + "HomeGetAlarmStatistic"
	containers.Set(key, home3.GetAlarmStatistic{})

	//组织架构
	key = consts.ValidatorPrefixOrg + "DepartmentTree"
	containers.Set(key, department3.Tree{})
	key = consts.ValidatorPrefixOrg + "DepartmentAdd"
	containers.Set(key, department3.Add{})
	key = consts.ValidatorPrefixOrg + "DepartmentEdit"
	containers.Set(key, department3.Edit{})
	key = consts.ValidatorPrefixOrg + "DepartmentDelete"
	containers.Set(key, department3.Delete{})

	//终端
	key = consts.ValidatorPrefixOrg + "AgentManagePage"
	containers.Set(key, agent3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "AgentManageDelete"
	containers.Set(key, agent3.ManageDelete{})
	key = consts.ValidatorPrefixOrg + "AgentManageUninstall"
	containers.Set(key, agent3.ManageUninstall{})
	key = consts.ValidatorPrefixOrg + "AgentManageUpgrade"
	containers.Set(key, agent3.ManageUpgrade{})
	key = consts.ValidatorPrefixOrg + "AgentManageBatchDelete"
	containers.Set(key, agent3.ManageBatchDelete{})
	key = consts.ValidatorPrefixOrg + "AgentManageBatchUninstall"
	containers.Set(key, agent3.ManageBatchUninstall{})
	key = consts.ValidatorPrefixOrg + "AgentManageBatchUpgrade"
	containers.Set(key, agent3.ManageBatchUpgrade{})
	key = consts.ValidatorPrefixOrg + "AgentManageBatchSetConfig"
	containers.Set(key, agent3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "AgentBasicGet"
	containers.Set(key, agent3.BasicGet{})
	key = consts.ValidatorPrefixOrg + "AgentBasicSave"
	containers.Set(key, agent3.BasicSave{})
	key = consts.ValidatorPrefixOrg + "AgentConfigGet"
	containers.Set(key, agent3.ConfigGet{})
	key = consts.ValidatorPrefixOrg + "AgentConfigSave"
	containers.Set(key, agent3.ConfigSave{})
	key = consts.ValidatorPrefixOrg + "AgentMonitorList"
	containers.Set(key, agent3.MonitorList{})
	key = consts.ValidatorPrefixOrg + "AgentMonitorInfo"
	containers.Set(key, agent3.MonitorInfo{})
	key = consts.ValidatorPrefixOrg + "AgentErrorPage"
	containers.Set(key, agent3.ErrorPage{})

	//进程
	key = consts.ValidatorPrefixOrg + "ProcessManagePage"
	containers.Set(key, process3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "ProcessManageSetSwitch"
	containers.Set(key, process3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "ProcessManageSetMode"
	containers.Set(key, process3.ManageSetMode{})
	key = consts.ValidatorPrefixOrg + "ProcessManageBatchSetConfig"
	containers.Set(key, process3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "ProcessManageBatchSetWhite"
	containers.Set(key, process3.ManageBatchSetWhite{})
	key = consts.ValidatorPrefixOrg + "ProcessManageBatchRefresh"
	containers.Set(key, process3.ManageBatchRefresh{})
	key = consts.ValidatorPrefixOrg + "ProcessManageBatchCollect"
	containers.Set(key, process3.ManageBatchCollect{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentPageAll"
	containers.Set(key, process3.AgentPageAll{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentPageRunning"
	containers.Set(key, process3.AgentPageRunning{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentSetStatus"
	containers.Set(key, process3.AgentSetStatus{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentCancel"
	containers.Set(key, process3.AgentCancel{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentClearAll"
	containers.Set(key, process3.AgentClearAll{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentClearRunning"
	containers.Set(key, process3.AgentClearRunning{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRefresh"
	containers.Set(key, process3.AgentRefresh{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentCollect"
	containers.Set(key, process3.AgentCollect{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentExtractRunningDll"
	containers.Set(key, process3.AgentExtractRunningDll{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentShowRunningDll"
	containers.Set(key, process3.AgentShowRunningDll{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentBatchWhite"
	containers.Set(key, process3.AgentShowRunningDll{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRepeatPage"
	containers.Set(key, process3.AgentRepeatHash{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentExtractSample"
	containers.Set(key, procSample3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentSetGlobal"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalPage"
	containers.Set(key, procGlobal3.PageFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalAdd"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalEdit"
	containers.Set(key, procGlobal3.Edit{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalDelete"
	containers.Set(key, procGlobal3.Delete{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalClear"
	containers.Set(key, procGlobal3.ClearFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessSearchPage"
	containers.Set(key, process3.SearchPage{})
	key = consts.ValidatorPrefixOrg + "ProcessSearchBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessSearchBatchGlobalBlack"
	containers.Set(key, procGlobal3.BatchBlackFromProcess{})

	//进程和模块样本
	key = consts.ValidatorPrefixOrg + "ProcSamplePage"
	containers.Set(key, procSample3.Page{})
	key = consts.ValidatorPrefixOrg + "ProcSampleDelete"
	containers.Set(key, procSample3.Delete{})
	key = consts.ValidatorPrefixOrg + "ProcSampleClear"
	containers.Set(key, procSample3.Clear{})
	key = consts.ValidatorPrefixOrg + "ProcSampleDownload"
	containers.Set(key, procSample3.Download{})

	//防篡改
	key = consts.ValidatorPrefixOrg + "TamperManagePage"
	containers.Set(key, tamper3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "TamperManageSetSwitch"
	containers.Set(key, tamper3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "TamperManageSetMode"
	containers.Set(key, tamper3.ManageSetMode{})
	key = consts.ValidatorPrefixOrg + "TamperManageBatchSetConfig"
	containers.Set(key, tamper3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyPage"
	containers.Set(key, tamper3.PolicyPage{})
	key = consts.ValidatorPrefixOrg + "TamperPolicySetSwitch"
	containers.Set(key, tamper3.PolicySetSwitch{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyAddFolder"
	containers.Set(key, tamper3.PolicyAddFolder{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyAddFile"
	containers.Set(key, tamper3.PolicyAddFile{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyEditFolder"
	containers.Set(key, tamper3.PolicyEditFolder{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyEditFile"
	containers.Set(key, tamper3.PolicyEditFile{})
	key = consts.ValidatorPrefixOrg + "TamperPolicyDelete"
	containers.Set(key, tamper3.PolicyDelete{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalPage"
	containers.Set(key, tamper3.GlobalPage{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalSetSwitch"
	containers.Set(key, tamper3.GlobalSetSwitch{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalAddFolder"
	containers.Set(key, tamper3.GlobalAddFolder{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalAddFile"
	containers.Set(key, tamper3.GlobalAddFile{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalEditFolder"
	containers.Set(key, tamper3.GlobalEditFolder{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalEditFile"
	containers.Set(key, tamper3.GlobalEditFile{})
	key = consts.ValidatorPrefixOrg + "TamperGlobalDelete"
	containers.Set(key, tamper3.GlobalDelete{})
}
