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
	key = consts.ValidatorPrefixPlatform + "AuthCheckUser"
	containers.Set(key, auth3.CheckUser{})
	key = consts.ValidatorPrefixPlatform + "AuthGetSliderCode"
	containers.Set(key, auth3.GetSliderCode{})
	key = consts.ValidatorPrefixPlatform + "AuthSendSms"
	containers.Set(key, auth3.SendSms{})
	key = consts.ValidatorPrefixPlatform + "AuthNormalLogin"
	containers.Set(key, auth3.NormalLogin{})
	key = consts.ValidatorPrefixPlatform + "AuthLogout"
	containers.Set(key, auth3.Logout{})
	key = consts.ValidatorPrefixPlatform + "AuthHsLogin"
	containers.Set(key, auth3.HsLogin{})

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeGetAgentStatistic"
	containers.Set(key, home3.GetAgentStatistic{})
	key = consts.ValidatorPrefixPlatform + "HomeGetAlarmStatistic"
	containers.Set(key, home3.GetAlarmStatistic{})

	//组织架构
	key = consts.ValidatorPrefixPlatform + "DepartmentTree"
	containers.Set(key, department3.Tree{})
	key = consts.ValidatorPrefixPlatform + "DepartmentAdd"
	containers.Set(key, department3.Add{})
	key = consts.ValidatorPrefixPlatform + "DepartmentEdit"
	containers.Set(key, department3.Edit{})
	key = consts.ValidatorPrefixPlatform + "DepartmentDelete"
	containers.Set(key, department3.Delete{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentManagePage"
	containers.Set(key, agent3.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "AgentManageDelete"
	containers.Set(key, agent3.ManageDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUninstall"
	containers.Set(key, agent3.ManageUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUpgrade"
	containers.Set(key, agent3.ManageUpgrade{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchDelete"
	containers.Set(key, agent3.ManageBatchDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUninstall"
	containers.Set(key, agent3.ManageBatchUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUpgrade"
	containers.Set(key, agent3.ManageBatchUpgrade{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchSetConfig"
	containers.Set(key, agent3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicGet"
	containers.Set(key, agent3.BasicGet{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicSave"
	containers.Set(key, agent3.BasicSave{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigGet"
	containers.Set(key, agent3.ConfigGet{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigSave"
	containers.Set(key, agent3.ConfigSave{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorList"
	containers.Set(key, agent3.MonitorList{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorInfo"
	containers.Set(key, agent3.MonitorInfo{})
	key = consts.ValidatorPrefixPlatform + "AgentErrorPage"
	containers.Set(key, agent3.ErrorPage{})

	//进程
	key = consts.ValidatorPrefixPlatform + "ProcessManagePage"
	containers.Set(key, process3.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetSwitch"
	containers.Set(key, process3.ManageSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetMode"
	containers.Set(key, process3.ManageSetMode{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetConfig"
	containers.Set(key, process3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetWhite"
	containers.Set(key, process3.ManageBatchSetWhite{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchRefresh"
	containers.Set(key, process3.ManageBatchRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchCollect"
	containers.Set(key, process3.ManageBatchCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageAll"
	containers.Set(key, process3.AgentPageAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageRunning"
	containers.Set(key, process3.AgentPageRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetStatus"
	containers.Set(key, process3.AgentSetStatus{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCancel"
	containers.Set(key, process3.AgentCancel{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearAll"
	containers.Set(key, process3.AgentClearAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearRunning"
	containers.Set(key, process3.AgentClearRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRefresh"
	containers.Set(key, process3.AgentRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCollect"
	containers.Set(key, process3.AgentCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentExtractRunningDll"
	containers.Set(key, process3.AgentExtractRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentShowRunningDll"
	containers.Set(key, process3.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchWhite"
	containers.Set(key, process3.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRepeatPage"
	containers.Set(key, process3.AgentRepeatHash{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentExtractSample"
	containers.Set(key, procSample3.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetGlobal"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalPage"
	containers.Set(key, procGlobal3.PageFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalAdd"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalEdit"
	containers.Set(key, procGlobal3.Edit{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalDelete"
	containers.Set(key, procGlobal3.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalClear"
	containers.Set(key, procGlobal3.ClearFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchPage"
	containers.Set(key, process3.SearchPage{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalBlack"
	containers.Set(key, procGlobal3.BatchBlackFromProcess{})

	//进程和模块样本
	key = consts.ValidatorPrefixPlatform + "ProcSamplePage"
	containers.Set(key, procSample3.Page{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDelete"
	containers.Set(key, procSample3.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleClear"
	containers.Set(key, procSample3.Clear{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDownload"
	containers.Set(key, procSample3.Download{})

	//防篡改
	key = consts.ValidatorPrefixPlatform + "TamperManagePage"
	containers.Set(key, tamper3.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "TamperManageSetSwitch"
	containers.Set(key, tamper3.ManageSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperManageSetMode"
	containers.Set(key, tamper3.ManageSetMode{})
	key = consts.ValidatorPrefixPlatform + "TamperManageBatchSetConfig"
	containers.Set(key, tamper3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyPage"
	containers.Set(key, tamper3.PolicyPage{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicySetSwitch"
	containers.Set(key, tamper3.PolicySetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyAddFolder"
	containers.Set(key, tamper3.PolicyAddFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyAddFile"
	containers.Set(key, tamper3.PolicyAddFile{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyEditFolder"
	containers.Set(key, tamper3.PolicyEditFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyEditFile"
	containers.Set(key, tamper3.PolicyEditFile{})
	key = consts.ValidatorPrefixPlatform + "TamperPolicyDelete"
	containers.Set(key, tamper3.PolicyDelete{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalPage"
	containers.Set(key, tamper3.GlobalPage{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalSetSwitch"
	containers.Set(key, tamper3.GlobalSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalAddFolder"
	containers.Set(key, tamper3.GlobalAddFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalAddFile"
	containers.Set(key, tamper3.GlobalAddFile{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalEditFolder"
	containers.Set(key, tamper3.GlobalEditFolder{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalEditFile"
	containers.Set(key, tamper3.GlobalEditFile{})
	key = consts.ValidatorPrefixPlatform + "TamperGlobalDelete"
	containers.Set(key, tamper3.GlobalDelete{})
}
