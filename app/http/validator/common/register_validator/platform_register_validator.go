package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	auth2 "goskeleton/app/http/validator/platform/auth"
	home2 "goskeleton/app/http/validator/platform/home"
	"goskeleton/app/http/validator/platform/org_menu"
)

// PlatformRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func PlatformRegisterValidator() {

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

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeAgentDailyOnline"
	containers.Set(key, home2.AgentDailyOnline{})
	key = consts.ValidatorPrefixPlatform + "HomeUserDailyOnline"
	containers.Set(key, home2.UserDailyOnline{})

	//菜单
	key = consts.ValidatorPrefixPlatform + "OrgMenuTree"
	containers.Set(key, org_menu.Tree{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuAdd"
	containers.Set(key, org_menu.Add{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuEdit"
	containers.Set(key, org_menu.Edit{})
	key = consts.ValidatorPrefixPlatform + "MenuDelete"
	containers.Set(key, org_menu.Delete{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentManagePage"

}
