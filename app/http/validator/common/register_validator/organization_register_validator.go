package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	agent3 "goskeleton/app/http/validator/organization/agent"
	auth3 "goskeleton/app/http/validator/organization/auth"
	blackmail3 "goskeleton/app/http/validator/organization/blackmail"
	department3 "goskeleton/app/http/validator/organization/department"
	dll3 "goskeleton/app/http/validator/organization/dll"
	findPw3 "goskeleton/app/http/validator/organization/find_pw"
	home3 "goskeleton/app/http/validator/organization/home"
	plug3 "goskeleton/app/http/validator/organization/plug"
	port3 "goskeleton/app/http/validator/organization/port"
	procGlobal3 "goskeleton/app/http/validator/organization/proc_global"
	procSample3 "goskeleton/app/http/validator/organization/proc_sample"
	process3 "goskeleton/app/http/validator/organization/process"
	registry3 "goskeleton/app/http/validator/organization/registry"
	tamper3 "goskeleton/app/http/validator/organization/tamper"
	threaten3 "goskeleton/app/http/validator/organization/threaten"
	trustDir3 "goskeleton/app/http/validator/organization/trust_dir"
	usb3 "goskeleton/app/http/validator/organization/usb"
	user3 "goskeleton/app/http/validator/organization/user"
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

	//注册
	key = consts.ValidatorPrefixOrg + "RegistryCheckOrgName"
	containers.Set(key, registry3.CheckOrgName{})
	key = consts.ValidatorPrefixOrg + "RegistryCheckUserName"
	containers.Set(key, registry3.CheckUserName{})
	key = consts.ValidatorPrefixOrg + "RegistryCheckUserPhone"
	containers.Set(key, registry3.CheckUserPhone{})
	key = consts.ValidatorPrefixOrg + "RegistryGetSliderCode"
	containers.Set(key, registry3.GetSliderCode{})
	key = consts.ValidatorPrefixOrg + "RegistrySendSms"
	containers.Set(key, registry3.SendSms{})
	key = consts.ValidatorPrefixOrg + "RegistrySignup"
	containers.Set(key, registry3.Signup{})

	//找回密码
	key = consts.ValidatorPrefixOrg + "FindPwGetSliderCode"
	containers.Set(key, findPw3.GetSliderCode{})
	key = consts.ValidatorPrefixOrg + "FindPwSendSms"
	containers.Set(key, findPw3.SendSms{})
	key = consts.ValidatorPrefixOrg + "FindPwCheckSms"
	containers.Set(key, findPw3.CheckSms{})
	key = consts.ValidatorPrefixOrg + "FindPwModifyPw"
	containers.Set(key, findPw3.ModifyPw{})

	//首页
	key = consts.ValidatorPrefixOrg + "HomeAlarmDailyStat"
	containers.Set(key, home3.AlarmDailyStat{})

	//组织架构
	key = consts.ValidatorPrefixOrg + "DepartmentAdd"
	containers.Set(key, department3.Add{})
	key = consts.ValidatorPrefixOrg + "DepartmentEdit"
	containers.Set(key, department3.Edit{})
	key = consts.ValidatorPrefixOrg + "DepartmentDelete"
	containers.Set(key, department3.Delete{})

	//用户管理
	key = consts.ValidatorPrefixOrg + "UserPage"
	containers.Set(key, user3.Page{})
	key = consts.ValidatorPrefixOrg + "UserAdd"
	containers.Set(key, user3.Add{})
	key = consts.ValidatorPrefixOrg + "UserEdit"
	containers.Set(key, user3.Edit{})
	key = consts.ValidatorPrefixOrg + "UserDelete"
	containers.Set(key, user3.Delete{})
	key = consts.ValidatorPrefixOrg + "UserModifyPw"
	containers.Set(key, user3.ModifyPw{})

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
	key = consts.ValidatorPrefixOrg + "AgentManageBatchModify"
	containers.Set(key, agent3.ManageBatchModify{})
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

	//信任目录
	key = consts.ValidatorPrefixOrg + "TrustDirGlobalPage"
	containers.Set(key, trustDir3.GlobalPage{})
	key = consts.ValidatorPrefixOrg + "TrustDirGlobalAdd"
	containers.Set(key, trustDir3.GlobalAdd{})
	key = consts.ValidatorPrefixOrg + "TrustDirGlobalEdit"
	containers.Set(key, trustDir3.GlobalEdit{})
	key = consts.ValidatorPrefixOrg + "TrustDirGlobalSetSwitch"
	containers.Set(key, trustDir3.GlobalSetSwitch{})
	key = consts.ValidatorPrefixOrg + "TrustDirGlobalDelete"
	containers.Set(key, trustDir3.GlobalDelete{})

	//进程
	key = consts.ValidatorPrefixOrg + "ProcessManagePage"
	containers.Set(key, process3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "ProcessManageAgentProcStat"
	containers.Set(key, process3.ManageAgentProcStat{})
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
	key = consts.ValidatorPrefixOrg + "ProcessAgentSearchWhile"
	containers.Set(key, process3.AgentSearchWhile{})
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
	key = consts.ValidatorPrefixOrg + "ProcessAgentBatchWhite"
	containers.Set(key, process3.AgentBatchWhite{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRunningDllExtract"
	containers.Set(key, process3.AgentRunningDllExtract{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRunningDllShow"
	containers.Set(key, process3.AgentRunningDllShow{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRepeatStat"
	containers.Set(key, process3.AgentRepeatStat{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentRepeatPage"
	containers.Set(key, process3.AgentRepeatPage{})

	//进程全局名单
	key = consts.ValidatorPrefixOrg + "ProcessAgentSetGlobal"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessAgentBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalPage"
	containers.Set(key, procGlobal3.PageFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalAdd"
	containers.Set(key, procGlobal3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalSetSwitch"
	containers.Set(key, procGlobal3.SetSwitchFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalDelete"
	containers.Set(key, procGlobal3.DeleteFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessGlobalClear"
	containers.Set(key, procGlobal3.ClearFromProcess{})
	/*key = consts.ValidatorPrefixOrg + "ProcessSearchPage"
	containers.Set(key, process3.SearchPage{})
	key = consts.ValidatorPrefixOrg + "ProcessSearchBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixOrg + "ProcessSearchBatchGlobalBlack"
	containers.Set(key, procGlobal3.BatchBlackFromProcess{})*/

	//模块
	key = consts.ValidatorPrefixOrg + "DllManagePage"
	containers.Set(key, dll3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "DllManageAgentProcStat"
	containers.Set(key, dll3.ManageAgentProcStat{})
	key = consts.ValidatorPrefixOrg + "DllManageSetSwitch"
	containers.Set(key, dll3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "DllManageSetMode"
	containers.Set(key, dll3.ManageSetMode{})
	key = consts.ValidatorPrefixOrg + "DllManageBatchSetConfig"
	containers.Set(key, dll3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "DllManageBatchSetWhite"
	containers.Set(key, dll3.ManageBatchSetWhite{})
	key = consts.ValidatorPrefixOrg + "DllManageBatchRefresh"
	containers.Set(key, dll3.ManageBatchRefresh{})
	key = consts.ValidatorPrefixOrg + "DllManageBatchCollect"
	containers.Set(key, dll3.ManageBatchCollect{})
	key = consts.ValidatorPrefixOrg + "DllAgentPageAll"
	containers.Set(key, dll3.AgentPageAll{})
	key = consts.ValidatorPrefixOrg + "DllAgentPageRunning"
	containers.Set(key, dll3.AgentPageRunning{})
	key = consts.ValidatorPrefixOrg + "DllAgentSetStatus"
	containers.Set(key, dll3.AgentSetStatus{})
	key = consts.ValidatorPrefixOrg + "DllAgentCancel"
	containers.Set(key, dll3.AgentCancel{})
	key = consts.ValidatorPrefixOrg + "DllAgentClearAll"
	containers.Set(key, dll3.AgentClearAll{})
	key = consts.ValidatorPrefixOrg + "DllAgentClearRunning"
	containers.Set(key, dll3.AgentClearRunning{})
	key = consts.ValidatorPrefixOrg + "DllAgentRefresh"
	containers.Set(key, dll3.AgentRefresh{})
	key = consts.ValidatorPrefixOrg + "DllAgentCollect"
	containers.Set(key, dll3.AgentCollect{})
	key = consts.ValidatorPrefixOrg + "DllAgentBatchWhite"
	containers.Set(key, dll3.AgentBatchWhite{})
	key = consts.ValidatorPrefixOrg + "DllAgentRunningDllExtract"
	key = consts.ValidatorPrefixOrg + "DllAgentRepeatStat"
	containers.Set(key, dll3.AgentRepeatStat{})
	key = consts.ValidatorPrefixOrg + "DllAgentRepeatPage"
	containers.Set(key, dll3.AgentRepeatPage{})

	//模块全局名单
	key = consts.ValidatorPrefixOrg + "DllAgentSetGlobal"
	containers.Set(key, procGlobal3.AddFromDll{})
	key = consts.ValidatorPrefixOrg + "DllAgentBatchGlobalWhite"
	containers.Set(key, procGlobal3.BatchWhiteFromDll{})
	key = consts.ValidatorPrefixOrg + "DllGlobalPage"
	containers.Set(key, procGlobal3.PageFromDll{})
	key = consts.ValidatorPrefixOrg + "DllGlobalAdd"
	containers.Set(key, procGlobal3.AddFromDll{})
	key = consts.ValidatorPrefixOrg + "DllGlobalSetSwitch"
	containers.Set(key, procGlobal3.SetSwitchFromDll{})
	key = consts.ValidatorPrefixOrg + "DllGlobalDelete"
	containers.Set(key, procGlobal3.DeleteFromDll{})
	key = consts.ValidatorPrefixOrg + "DllGlobalClear"
	containers.Set(key, procGlobal3.ClearFromDll{})

	//进程和模块样本
	key = consts.ValidatorPrefixOrg + "ProcessAgentSampleExtract"
	containers.Set(key, procSample3.AddFromProcess{})
	key = consts.ValidatorPrefixOrg + "DllAgentSampleExtract"
	containers.Set(key, procSample3.AddFromDll{})
	key = consts.ValidatorPrefixOrg + "ProcSamplePage"
	containers.Set(key, procSample3.Page{})
	key = consts.ValidatorPrefixOrg + "ProcSampleDelete"
	containers.Set(key, procSample3.Delete{})
	key = consts.ValidatorPrefixOrg + "ProcSampleClear"
	containers.Set(key, procSample3.Clear{})
	key = consts.ValidatorPrefixOrg + "ProcSampleDownload"
	containers.Set(key, procSample3.Download{})

	//威胁情报
	key = consts.ValidatorPrefixOrg + "ThreatenGetLevels"
	containers.Set(key, threaten3.GetLevels{})
	key = consts.ValidatorPrefixOrg + "ThreatenGetAttr"
	containers.Set(key, threaten3.GetAttr{})

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

	//勒索防护
	key = consts.ValidatorPrefixOrg + "BlackmailGroupPage"
	containers.Set(key, blackmail3.GroupPage{})
	key = consts.ValidatorPrefixOrg + "BlackmailGroupGetAll"
	containers.Set(key, blackmail3.GroupGetAll{})
	key = consts.ValidatorPrefixOrg + "BlackmailGroupAdd"
	containers.Set(key, blackmail3.GroupAdd{})
	key = consts.ValidatorPrefixOrg + "BlackmailGroupEdit"
	containers.Set(key, blackmail3.GroupEdit{})
	key = consts.ValidatorPrefixOrg + "BlackmailGroupDelete"
	containers.Set(key, blackmail3.GroupDelete{})
	key = consts.ValidatorPrefixOrg + "BlackmailManagePage"
	containers.Set(key, blackmail3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "BlackmailManageSetSwitch"
	containers.Set(key, blackmail3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "BlackmailManageSetMode"
	containers.Set(key, blackmail3.ManageSetMode{})
	key = consts.ValidatorPrefixOrg + "BlackmailManageBatchSetConfig"
	containers.Set(key, blackmail3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "BlackmailPolicyPage"
	containers.Set(key, blackmail3.PolicyPage{})
	key = consts.ValidatorPrefixOrg + "BlackmailPolicySetSwitch"
	containers.Set(key, blackmail3.PolicySetSwitch{})
	key = consts.ValidatorPrefixOrg + "BlackmailPolicyAdd"
	containers.Set(key, blackmail3.PolicyAdd{})
	key = consts.ValidatorPrefixOrg + "BlackmailPolicyEdit"
	containers.Set(key, blackmail3.PolicyEdit{})
	key = consts.ValidatorPrefixOrg + "BlackmailPolicyDelete"
	containers.Set(key, blackmail3.PolicyDelete{})
	key = consts.ValidatorPrefixOrg + "BlackmailGlobalPage"
	containers.Set(key, blackmail3.GlobalPage{})
	key = consts.ValidatorPrefixOrg + "BlackmailGlobalSetSwitch"
	containers.Set(key, blackmail3.GlobalSetSwitch{})
	key = consts.ValidatorPrefixOrg + "BlackmailGlobalAdd"
	containers.Set(key, blackmail3.GlobalAdd{})
	key = consts.ValidatorPrefixOrg + "BlackmailGlobalEdit"
	containers.Set(key, blackmail3.GlobalEdit{})
	key = consts.ValidatorPrefixOrg + "BlackmailGlobalDelete"
	containers.Set(key, blackmail3.GlobalDelete{})

	//端口虚开
	key = consts.ValidatorPrefixOrg + "PortWorkingPage"
	containers.Set(key, port3.WorkingPage{})
	key = consts.ValidatorPrefixOrg + "PortWorkingRefresh"
	containers.Set(key, port3.WorkingRefresh{})
	key = consts.ValidatorPrefixOrg + "PortManagePage"
	containers.Set(key, port3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "PortManageSetSwitch"
	containers.Set(key, port3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PortManageBatchSetConfig"
	containers.Set(key, port3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "PortPolicyPage"
	containers.Set(key, port3.PolicyPage{})
	key = consts.ValidatorPrefixOrg + "PortPolicySetSwitch"
	containers.Set(key, port3.PolicySetSwitch{})
	key = consts.ValidatorPrefixOrg + "PortPolicyAdd"
	containers.Set(key, port3.PolicyAdd{})
	key = consts.ValidatorPrefixOrg + "PortPolicyEdit"
	containers.Set(key, port3.PolicyEdit{})
	key = consts.ValidatorPrefixOrg + "PortPolicyDelete"
	containers.Set(key, port3.PolicyDelete{})
	key = consts.ValidatorPrefixOrg + "PortGlobalPage"
	containers.Set(key, port3.GlobalPage{})
	key = consts.ValidatorPrefixOrg + "PortGlobalSetSwitch"
	containers.Set(key, port3.GlobalSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PortGlobalAdd"
	containers.Set(key, port3.GlobalAdd{})
	key = consts.ValidatorPrefixOrg + "PortGlobalEdit"
	containers.Set(key, port3.GlobalEdit{})
	key = consts.ValidatorPrefixOrg + "PortGlobalDelete"
	containers.Set(key, port3.GlobalDelete{})

	//动态封堵
	key = consts.ValidatorPrefixOrg + "PlugManagePage"
	containers.Set(key, plug3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "PlugManageSetSwitch"
	containers.Set(key, plug3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PlugManageBatchSetConfig"
	containers.Set(key, plug3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyConfigPortGet"
	containers.Set(key, plug3.PolicyConfigPortGet{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyConfigPortSet"
	containers.Set(key, plug3.PolicyConfigPortSet{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyConfigLoginGet"
	containers.Set(key, plug3.PolicyConfigLoginGet{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyConfigLoginSet"
	containers.Set(key, plug3.PolicyConfigLoginSet{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicEnabledPage"
	containers.Set(key, plug3.PolicyDynamicEnabledPage{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicDisabledPage"
	containers.Set(key, plug3.PolicyDynamicDisabledPage{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicDisable"
	containers.Set(key, plug3.PolicyDynamicDisable{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicDisableBatch"
	containers.Set(key, plug3.PolicyDynamicDisableBatch{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicDelete"
	containers.Set(key, plug3.PolicyDynamicDelete{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyDynamicDeleteBatch"
	containers.Set(key, plug3.PolicyDynamicDeleteBatch{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhitePage"
	containers.Set(key, plug3.PolicyIpWhitePage{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhiteSetSwitch"
	containers.Set(key, plug3.PolicyIpWhiteSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhiteAdd"
	containers.Set(key, plug3.PolicyIpWhiteAdd{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhiteEdit"
	containers.Set(key, plug3.PolicyIpWhiteEdit{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhiteDelete"
	containers.Set(key, plug3.PolicyIpWhiteDelete{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpWhiteClear"
	containers.Set(key, plug3.PolicyIpWhiteClear{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackPage"
	containers.Set(key, plug3.PolicyIpBlackPage{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackSetSwitch"
	containers.Set(key, plug3.PolicyIpBlackSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackAdd"
	containers.Set(key, plug3.PolicyIpBlackAdd{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackEdit"
	containers.Set(key, plug3.PolicyIpBlackEdit{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackDelete"
	containers.Set(key, plug3.PolicyIpBlackDelete{})
	key = consts.ValidatorPrefixOrg + "PlugPolicyIpBlackClear"
	containers.Set(key, plug3.PolicyIpBlackClear{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalConfigPortGet"
	containers.Set(key, plug3.GlobalConfigPortGet{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalConfigPortSet"
	containers.Set(key, plug3.GlobalConfigPortSet{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalConfigLoginGet"
	containers.Set(key, plug3.GlobalConfigLoginGet{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalConfigLoginSet"
	containers.Set(key, plug3.GlobalConfigLoginSet{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicEnabledPage"
	containers.Set(key, plug3.GlobalDynamicEnabledPage{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicDisabledPage"
	containers.Set(key, plug3.GlobalDynamicDisabledPage{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicDisable"
	containers.Set(key, plug3.GlobalDynamicDisable{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicDisableBatch"
	containers.Set(key, plug3.GlobalDynamicDisableBatch{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicDelete"
	containers.Set(key, plug3.GlobalDynamicDelete{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalDynamicDeleteBatch"
	containers.Set(key, plug3.GlobalDynamicDeleteBatch{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhitePage"
	containers.Set(key, plug3.GlobalIpWhitePage{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhiteSetSwitch"
	containers.Set(key, plug3.GlobalIpWhiteSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhiteAdd"
	containers.Set(key, plug3.GlobalIpWhiteAdd{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhiteEdit"
	containers.Set(key, plug3.GlobalIpWhiteEdit{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhiteDelete"
	containers.Set(key, plug3.GlobalIpWhiteDelete{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpWhiteClear"
	containers.Set(key, plug3.GlobalIpWhiteClear{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackPage"
	containers.Set(key, plug3.GlobalIpBlackPage{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackSetSwitch"
	containers.Set(key, plug3.GlobalIpBlackSetSwitch{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackAdd"
	containers.Set(key, plug3.GlobalIpBlackAdd{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackEdit"
	containers.Set(key, plug3.GlobalIpBlackEdit{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackDelete"
	containers.Set(key, plug3.GlobalIpBlackDelete{})
	key = consts.ValidatorPrefixOrg + "PlugGlobalIpBlackClear"
	containers.Set(key, plug3.GlobalIpBlackClear{})

	//外接设备
	key = consts.ValidatorPrefixOrg + "UsbManagePage"
	containers.Set(key, usb3.ManagePage{})
	key = consts.ValidatorPrefixOrg + "UsbManageSetSwitch"
	containers.Set(key, usb3.ManageSetSwitch{})
	key = consts.ValidatorPrefixOrg + "UsbManageSetMode"
	containers.Set(key, usb3.ManageSetMode{})
	key = consts.ValidatorPrefixOrg + "UsbManageBatchSetConfig"
	containers.Set(key, usb3.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyPageRunning"
	containers.Set(key, usb3.PolicyPageRunning{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyPageHistory"
	containers.Set(key, usb3.PolicyPageHistory{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyPageLog"
	containers.Set(key, usb3.PolicyPageLog{})
	key = consts.ValidatorPrefixOrg + "UsbPolicySetWhite"
	containers.Set(key, usb3.PolicySetWhite{})
	key = consts.ValidatorPrefixOrg + "UsbPolicySetBlack"
	containers.Set(key, usb3.PolicySetBlack{})
	key = consts.ValidatorPrefixOrg + "UsbPolicySetAlias"
	containers.Set(key, usb3.PolicySetAlias{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyRefresh"
	containers.Set(key, usb3.PolicyRefresh{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyHasUpload"
	containers.Set(key, usb3.PolicyHasUpload{})
	key = consts.ValidatorPrefixOrg + "UsbPolicyAddGlobal"
	containers.Set(key, usb3.GlobalAdd{})
	key = consts.ValidatorPrefixOrg + "UsbGlobalPage"
	containers.Set(key, usb3.GlobalPage{})
	key = consts.ValidatorPrefixOrg + "UsbGlobalSetWhite"
	containers.Set(key, usb3.GlobalSetWhite{})
	key = consts.ValidatorPrefixOrg + "UsbGlobalSetBlack"
	containers.Set(key, usb3.GlobalSetBlack{})
	key = consts.ValidatorPrefixOrg + "UsbGlobalSetAlias"
	containers.Set(key, usb3.GlobalSetAlias{})
	key = consts.ValidatorPrefixOrg + "UsbGlobalDelete"
	containers.Set(key, usb3.GlobalDelete{})
}
