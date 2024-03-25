package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/platform/agent"
	"goskeleton/app/http/validator/platform/auth"
	"goskeleton/app/http/validator/platform/department"
	"goskeleton/app/http/validator/platform/home"
	"goskeleton/app/http/validator/platform/menu"
	"goskeleton/app/http/validator/platform/proc_global"
	"goskeleton/app/http/validator/platform/proc_sample"
	"goskeleton/app/http/validator/platform/process"
	"goskeleton/app/http/validator/platform/role"
	"goskeleton/app/http/validator/platform/threaten_push"
)

// PlatformRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func PlatformRegisterValidator() {

	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//安全认证
	key = consts.ValidatorPrefixPlatform + "AuthCheckUser"
	containers.Set(key, auth.CheckUser{})
	key = consts.ValidatorPrefixPlatform + "AuthGetSliderCode"
	containers.Set(key, auth.GetSliderCode{})
	key = consts.ValidatorPrefixPlatform + "AuthSendSms"
	containers.Set(key, auth.SendSms{})
	key = consts.ValidatorPrefixPlatform + "AuthNormalLogin"
	containers.Set(key, auth.NormalLogin{})
	key = consts.ValidatorPrefixPlatform + "AuthLogout"
	containers.Set(key, auth.Logout{})
	key = consts.ValidatorPrefixPlatform + "AuthHsLogin"
	containers.Set(key, auth.HsLogin{})

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeGetAgentStatistic"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefixPlatform + "HomeGetAlarmStatistic"
	containers.Set(key, home.GetAlarmStatistic{})

	//菜单
	key = consts.ValidatorPrefixPlatform + "MenuTree"
	containers.Set(key, menu.Tree{})
	key = consts.ValidatorPrefixPlatform + "MenuAdd"
	containers.Set(key, menu.Add{})
	key = consts.ValidatorPrefixPlatform + "MenuEdit"
	containers.Set(key, menu.Edit{})
	key = consts.ValidatorPrefixPlatform + "MenuSort"
	containers.Set(key, menu.Sort{})
	key = consts.ValidatorPrefixPlatform + "MenuDelete"
	containers.Set(key, menu.Delete{})

	//角色
	key = consts.ValidatorPrefixPlatform + "RoleManageAll"
	containers.Set(key, role.ManageAll{})
	key = consts.ValidatorPrefixPlatform + "RoleManageAdd"
	containers.Set(key, role.ManageAdd{})
	key = consts.ValidatorPrefixPlatform + "RoleManageEdit"
	containers.Set(key, role.ManageEdit{})
	key = consts.ValidatorPrefixPlatform + "RoleManageDelete"
	containers.Set(key, role.ManageDelete{})
	key = consts.ValidatorPrefixPlatform + "RolePowerLoginMenuTree"
	containers.Set(key, role.PowerLoginMenuTree{})
	key = consts.ValidatorPrefixPlatform + "RolePowerEditableMenuTree"
	containers.Set(key, role.PowerEditableMenuTree{})
	key = consts.ValidatorPrefixPlatform + "RolePowerGetSelectedMenus"
	containers.Set(key, role.GetSelectedMenus{})
	key = consts.ValidatorPrefixPlatform + "RolePowerSetSelectedMenus"
	containers.Set(key, role.SetSelectedMenus{})

	//组织架构
	key = consts.ValidatorPrefixPlatform + "DepartmentTree"
	containers.Set(key, department.Tree{})
	key = consts.ValidatorPrefixPlatform + "DepartmentAdd"
	containers.Set(key, department.Add{})
	key = consts.ValidatorPrefixPlatform + "DepartmentEdit"
	containers.Set(key, department.Edit{})
	key = consts.ValidatorPrefixPlatform + "DepartmentDelete"
	containers.Set(key, department.Delete{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentManagePage"
	containers.Set(key, agent.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "AgentManageDelete"
	containers.Set(key, agent.ManageDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUninstall"
	containers.Set(key, agent.ManageUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageUpgrade"
	containers.Set(key, agent.ManageUpgrade{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchDelete"
	containers.Set(key, agent.ManageBatchDelete{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUninstall"
	containers.Set(key, agent.ManageBatchUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchUpgrade"
	containers.Set(key, agent.ManageBatchUninstall{})
	key = consts.ValidatorPrefixPlatform + "AgentManageBatchSetConfig"
	containers.Set(key, agent.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicGet"
	containers.Set(key, agent.BasicGet{})
	key = consts.ValidatorPrefixPlatform + "AgentBasicSave"
	containers.Set(key, agent.BasicSave{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigGet"
	containers.Set(key, agent.ConfigGet{})
	key = consts.ValidatorPrefixPlatform + "AgentConfigSave"
	containers.Set(key, agent.ConfigSave{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorList"
	containers.Set(key, agent.MonitorList{})
	key = consts.ValidatorPrefixPlatform + "AgentMonitorInfo"
	containers.Set(key, agent.MonitorInfo{})
	key = consts.ValidatorPrefixPlatform + "AgentErrorPage"
	containers.Set(key, agent.ErrorPage{})

	//进程
	key = consts.ValidatorPrefixPlatform + "ProcessManagePage"
	containers.Set(key, process.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetSwitch"
	containers.Set(key, process.ManageSetSwitch{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageSetMode"
	containers.Set(key, process.ManageSetMode{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetConfig"
	containers.Set(key, process.ManageBatchSetConfig{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchSetWhite"
	containers.Set(key, process.ManageBatchSetWhite{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchRefresh"
	containers.Set(key, process.ManageBatchRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchCollect"
	containers.Set(key, process.ManageBatchCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessManageBatchThreatenWhite"
	containers.Set(key, process.ManageBatchThreatenWhite{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageAll"
	containers.Set(key, process.AgentPageAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentPageRunning"
	containers.Set(key, process.AgentPageRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetStatus"
	containers.Set(key, process.AgentSetStatus{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCancel"
	containers.Set(key, process.AgentCancel{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearAll"
	containers.Set(key, process.AgentClearAll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentClearRunning"
	containers.Set(key, process.AgentClearRunning{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRefresh"
	containers.Set(key, process.AgentRefresh{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentCollect"
	containers.Set(key, process.AgentCollect{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentShowRunningDll"
	containers.Set(key, process.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchWhite"
	containers.Set(key, process.AgentShowRunningDll{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentRepeatPage"
	containers.Set(key, process.AgentRepeatHash{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentExtractSample"
	containers.Set(key, proc_sample.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentSetGlobal"
	containers.Set(key, proc_global.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchGlobalWhite"
	containers.Set(key, proc_global.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchThreatenWhite"
	containers.Set(key, threaten_push.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessAgentBatchThreatenBlack"
	containers.Set(key, threaten_push.BatchBlackFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalPage"
	containers.Set(key, proc_global.PageFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalAdd"
	containers.Set(key, proc_global.AddFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalEdit"
	containers.Set(key, proc_global.Edit{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalDelete"
	containers.Set(key, proc_global.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcessGlobalClear"
	containers.Set(key, proc_global.ClearFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchPage"
	containers.Set(key, process.SearchPage{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalWhite"
	containers.Set(key, proc_global.BatchWhiteFromProcess{})
	key = consts.ValidatorPrefixPlatform + "ProcessSearchBatchGlobalBlack"
	containers.Set(key, proc_global.BatchBlackFromProcess{})

	//进程和模块样本
	key = consts.ValidatorPrefixPlatform + "ProcSamplePage"
	containers.Set(key, proc_sample.Page{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDelete"
	containers.Set(key, proc_sample.Delete{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleClear"
	containers.Set(key, proc_sample.Clear{})
	key = consts.ValidatorPrefixPlatform + "ProcSampleDownload"
	containers.Set(key, proc_sample.Download{})

}
