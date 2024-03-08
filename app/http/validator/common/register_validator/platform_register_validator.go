package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/platform/auth"
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
	containers.Set(key, auth.Login{})
	key = consts.ValidatorPrefix + "AuthLogout"
	containers.Set(key, auth.Logout{})
	key = consts.ValidatorPrefix + "AuthHsLogin"
	containers.Set(key, auth.HsLogin{})
}
