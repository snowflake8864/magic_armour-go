package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/platform/auth"
	"goskeleton/app/http/validator/platform/department"
	"goskeleton/app/http/validator/platform/home"
	"goskeleton/app/http/validator/platform/menu"
)

// PlatformRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func PlatformRegisterValidator() {

	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//安全认证
	key = consts.ValidatorPrefix + "AuthCheckUser"
	containers.Set(key, auth.CheckUser{})
	key = consts.ValidatorPrefix + "AuthGetSliderCode"
	containers.Set(key, auth.GetSliderCode{})
	key = consts.ValidatorPrefix + "AuthSendSms"
	containers.Set(key, auth.SendSms{})
	key = consts.ValidatorPrefix + "AuthNormalLogin"
	containers.Set(key, auth.NormalLogin{})
	key = consts.ValidatorPrefix + "AuthLogout"
	containers.Set(key, auth.Logout{})
	key = consts.ValidatorPrefix + "AuthHsLogin"
	containers.Set(key, auth.HsLogin{})

	//首页
	key = consts.ValidatorPrefix + "HomeGetAgentStatistic"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "HomeGetAlarmStatistic"
	containers.Set(key, home.GetAlarmStatistic{})

	//菜单
	key = consts.ValidatorPrefix + "MenuTree"
	containers.Set(key, menu.Tree{})
	key = consts.ValidatorPrefix + "MenuAdd"
	containers.Set(key, menu.Add{})
	key = consts.ValidatorPrefix + "MenuEdit"
	containers.Set(key, menu.Edit{})
	key = consts.ValidatorPrefix + "MenuSort"
	containers.Set(key, menu.Sort{})
	key = consts.ValidatorPrefix + "MenuDelete"
	containers.Set(key, menu.Delete{})

	//角色
	key = consts.ValidatorPrefix + "RoleManageAll"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RoleManageAdd"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RoleManageEdit"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RoleManageDelete"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RolePowerLoginMenuTree"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RolePowerEditableMenuTree"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RolePowerGetSelectedMenus"
	containers.Set(key, home.GetAgentStatistic{})
	key = consts.ValidatorPrefix + "RolePowerSetSelectedMenus"
	containers.Set(key, home.GetAgentStatistic{})

	//组织架构
	key = consts.ValidatorPrefix + "DepartmentTree"
	containers.Set(key, department.Tree{})
	key = consts.ValidatorPrefix + "DepartmentAdd"
	containers.Set(key, department.Add{})
	key = consts.ValidatorPrefix + "DepartmentEdit"
	containers.Set(key, department.Edit{})
	key = consts.ValidatorPrefix + "DepartmentDelete"
	containers.Set(key, department.Delete{})

}
