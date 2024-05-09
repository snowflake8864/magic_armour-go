package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	uploadFiles4 "goskeleton/app/http/validator/common/upload_files"
	account4 "goskeleton/app/http/validator/platform/account"
	agent4 "goskeleton/app/http/validator/platform/agent"
	auth4 "goskeleton/app/http/validator/platform/auth"
	dictionary4 "goskeleton/app/http/validator/platform/dictionary"
	home4 "goskeleton/app/http/validator/platform/home"
	operation4 "goskeleton/app/http/validator/platform/operation"
	order4 "goskeleton/app/http/validator/platform/order"
	orgMenu4 "goskeleton/app/http/validator/platform/org_menu"
	organize4 "goskeleton/app/http/validator/platform/organize"
	product4 "goskeleton/app/http/validator/platform/product"
	software4 "goskeleton/app/http/validator/platform/software"
	system4 "goskeleton/app/http/validator/platform/system"
	threaten4 "goskeleton/app/http/validator/platform/threaten"
)

// PlatformRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func PlatformRegisterValidator() {

	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string
	// 文件上传
	key = consts.ValidatorPrefixPlatform + "UploadFiles"
	containers.Set(key, uploadFiles4.PlatformUpFiles{})

	//安全认证
	key = consts.ValidatorPrefixPlatform + "AuthCheckUser"
	containers.Set(key, auth4.CheckUser{})
	key = consts.ValidatorPrefixPlatform + "AuthGetSliderCode"
	containers.Set(key, auth4.GetSliderCode{})
	key = consts.ValidatorPrefixPlatform + "AuthSendSms"
	containers.Set(key, auth4.SendSms{})
	key = consts.ValidatorPrefixPlatform + "AuthNormalLogin"
	containers.Set(key, auth4.NormalLogin{})

	//

	//首页
	key = consts.ValidatorPrefixPlatform + "HomeAgentDailyOnline"
	containers.Set(key, home4.AgentDailyOnline{})
	key = consts.ValidatorPrefixPlatform + "HomeUserDailyOnline"
	containers.Set(key, home4.UserDailyOnline{})

	//账户管理
	key = consts.ValidatorPrefixPlatform + "AccountModifyPw"
	containers.Set(key, account4.ModifyPw{})

	//数据字典
	key = consts.ValidatorPrefixPlatform + "DictionarySave"
	containers.Set(key, dictionary4.Save{})

	//菜单
	key = consts.ValidatorPrefixPlatform + "OrgMenuTree"
	containers.Set(key, orgMenu4.Tree{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuAdd"
	containers.Set(key, orgMenu4.Add{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuEdit"
	containers.Set(key, orgMenu4.Edit{})
	key = consts.ValidatorPrefixPlatform + "OrgMenuDelete"
	containers.Set(key, orgMenu4.Delete{})

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
	key = consts.ValidatorPrefixPlatform + "OrganizeLogPage"
	containers.Set(key, organize4.LogPage{})

	//终端
	key = consts.ValidatorPrefixPlatform + "AgentPublicStatByOrg"
	containers.Set(key, agent4.PublicStatByOrg{})
	key = consts.ValidatorPrefixPlatform + "AgentOrgPage"
	containers.Set(key, organize4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "AgentOrgAddGrant"
	containers.Set(key, order4.OrgAddGrant{})
	key = consts.ValidatorPrefixPlatform + "AgentOwnPage"
	containers.Set(key, agent4.OwnPage{})

	//订单
	key = consts.ValidatorPrefixPlatform + "OrderPublicValidByOrg"
	containers.Set(key, order4.ValidByOrg{})
	key = consts.ValidatorPrefixPlatform + "OrderGrantPage"
	containers.Set(key, order4.GrantPage{})
	key = consts.ValidatorPrefixPlatform + "OrderManagePage"
	containers.Set(key, order4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "OrderManageDetail"
	containers.Set(key, order4.ManageDetail{})
	key = consts.ValidatorPrefixPlatform + "OrderManageLogs"
	containers.Set(key, order4.ManageLogs{})
	key = consts.ValidatorPrefixPlatform + "OrderManageDecline"
	containers.Set(key, order4.ManageDecline{})
	key = consts.ValidatorPrefixPlatform + "OrderManageCancel"
	containers.Set(key, order4.ManageCancel{})
	key = consts.ValidatorPrefixPlatform + "OrderManagePaid"
	containers.Set(key, order4.ManagePaid{})
	key = consts.ValidatorPrefixPlatform + "OrderManageUnpaid"
	containers.Set(key, order4.ManageUnpaid{})

	//产品
	key = consts.ValidatorPrefixPlatform + "ProductManagePage"
	containers.Set(key, product4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "ProductManageAdd"
	containers.Set(key, product4.ManageAdd{})
	key = consts.ValidatorPrefixPlatform + "ProductManageEdit"
	containers.Set(key, product4.ManageEdit{})
	key = consts.ValidatorPrefixPlatform + "ProductManageDelete"
	containers.Set(key, product4.ManageDelete{})

	//软件
	key = consts.ValidatorPrefixPlatform + "SoftwarePublicOsTypes"
	containers.Set(key, software4.PublicOsTypes{})
	key = consts.ValidatorPrefixPlatform + "SoftwarePublicAgentConfigSet"
	containers.Set(key, software4.PublicAgentConfigSet{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManagePage"
	containers.Set(key, software4.ManagePage{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManageAdd"
	containers.Set(key, software4.ManageAdd{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManageDownload"
	containers.Set(key, software4.ManageDownload{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManagePublish"
	containers.Set(key, software4.ManagePublish{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManageCancel"
	containers.Set(key, software4.ManageCancel{})
	key = consts.ValidatorPrefixPlatform + "SoftwareManageDelete"
	containers.Set(key, software4.ManageDelete{})

	//威胁情报
	key = consts.ValidatorPrefixPlatform + "ThreatenConfigSave"
	containers.Set(key, threaten4.ConfigSave{})
	key = consts.ValidatorPrefixPlatform + "ThreatenConfigOnline"
	containers.Set(key, threaten4.ConfigOnline{})
	key = consts.ValidatorPrefixPlatform + "ThreatenCachePage"
	containers.Set(key, threaten4.CachePage{})
	key = consts.ValidatorPrefixPlatform + "ThreatenCacheGetAttr"
	containers.Set(key, threaten4.CacheGetAttr{})
	key = consts.ValidatorPrefixPlatform + "ThreatenCacheDelete"
	containers.Set(key, threaten4.CacheDelete{})
	key = consts.ValidatorPrefixPlatform + "ThreatenCacheBatchWhite"
	containers.Set(key, threaten4.PushBatchWhite{})
	key = consts.ValidatorPrefixPlatform + "ThreatenCacheBatchBlack"
	containers.Set(key, threaten4.PushBatchBlack{})
	key = consts.ValidatorPrefixPlatform + "ThreatenPushPage"
	containers.Set(key, threaten4.PushPage{})

	//操作相关
	key = consts.ValidatorPrefixPlatform + "OperationPage"
	containers.Set(key, operation4.Page{})

	//服务器相关
	key = consts.ValidatorPrefixPlatform + "SystemMonitorRecords"
	containers.Set(key, system4.MonitorRecords{})
}
