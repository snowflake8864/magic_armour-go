package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	agent4 "goskeleton/app/http/validator/platform/agent"
	auth4 "goskeleton/app/http/validator/platform/auth"
	home4 "goskeleton/app/http/validator/platform/home"
	"goskeleton/app/http/validator/platform/org_menu"
	organize4 "goskeleton/app/http/validator/platform/organize"
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
	containers.Set(key, auth4.CheckUser{})
	key = consts.ValidatorPrefixPlatform + "AuthGetSliderCode"
	containers.Set(key, auth4.GetSliderCode{})
	key = consts.ValidatorPrefixPlatform + "AuthSendSms"
	containers.Set(key, auth4.SendSms{})
	key = consts.ValidatorPrefixPlatform + "AuthNormalLogin"
	containers.Set(key, auth4.NormalLogin{})

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeAgentDailyOnline"
	containers.Set(key, home4.AgentDailyOnline{})
	key = consts.ValidatorPrefixPlatform + "HomeUserDailyOnline"
	containers.Set(key, home4.UserDailyOnline{})

	//菜单
	key = consts.ValidatorPrefixPlatform + "OrgMenuTree"
	containers.Set(key, org_menu.Tree{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuAdd"
	containers.Set(key, org_menu.Add{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuEdit"
	containers.Set(key, org_menu.Edit{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuDelete"
	containers.Set(key, org_menu.Delete{})

	//企业
	key = consts.ValidatorPrefixPlatform + "OrganizePublicDetail"
	containers.Set(key, organize4.PublicDetail{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManagePage"
	containers.Set(key, organize4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManageLock"
	containers.Set(key, organize4.ManageLock{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManageUnlock"
	containers.Set(key, organize4.ManageUnlock{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManageClean"
	containers.Set(key, organize4.ManageClean{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManageSetHsSwitch"
	containers.Set(key, organize4.ManageSetHsSwitch{})
	key = consts.ValidatorPrefixPlatform + "OrganizeManageSetLoginSmsSwitch"
	containers.Set(key, organize4.ManageSetLoginSmsSwitch{})
	key = consts.ValidatorPrefixPlatform + "OrganizeVerifyPage"
	containers.Set(key, organize4.VerifyPage{})
	key = consts.ValidatorPrefixPlatform + "OrganizeVerifyAdd"
	containers.Set(key, organize4.VerifyAdd{})
	key = consts.ValidatorPrefixPlatform + "OrganizeVerifyApprove"
	containers.Set(key, organize4.VerifyApprove{})
	key = consts.ValidatorPrefixPlatform + "OrganizeVerifyRefuse"
	containers.Set(key, organize4.VerifyRefuse{})
	key = consts.ValidatorPrefixPlatform + "OrganizeUserPage"
	containers.Set(key, organize4.UserPage{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentPublicStatByOrg"
	containers.Set(key, agent4.PublicStatByOrg{})
	key = consts.ValidatorPrefixPlatform + "AgentOrgPage"
	containers.Set(key, organize4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "AgentOwnPage"
	containers.Set(key, agent4.OwnPage{})

	//威胁情报
	key = consts.ValidatorPrefixPlatform + "ThreatenPushPage"
	containers.Set(key, threaten_push.Page{})
}
