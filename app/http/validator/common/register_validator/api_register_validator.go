package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/api/huanshi"
	"goskeleton/app/http/validator/api/software"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func ApiRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//软件
	key = consts.ValidatorPrefix + "SoftwarePublished"
	containers.Set(key, software.Published{})
	key = consts.ValidatorPrefix + "SoftwareDownload"
	containers.Set(key, software.Download{})

	//幻势
	key = consts.ValidatorPrefix + "HuanshiAuth"
	containers.Set(key, huanshi.Auth{})
}
