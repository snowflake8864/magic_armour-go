package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/agent/agent"
	"goskeleton/app/http/validator/agent/auth"
	"goskeleton/app/http/validator/agent/task"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func AgentRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//登录验证
	key = consts.ValidatorPrefix + "AuthLogin"
	containers.Set(key, auth.Login{})

	//终端相关
	key = consts.ValidatorPrefix + "AgentConfigGet"
	containers.Set(key, agent.ConfigGet{})
	key = consts.ValidatorPrefix + "AgentInfoGet"
	containers.Set(key, agent.InfoGet{})
	key = consts.ValidatorPrefix + "AgentUpgrade"
	containers.Set(key, agent.Upgrade{})
	key = consts.ValidatorPrefix + "AgentUninstall"
	containers.Set(key, agent.Uninstall{})
	key = consts.ValidatorPrefix + "TrustDirGet"
	containers.Set(key, agent.GetTrustDir{})
	key = consts.ValidatorPrefix + "ErrorFileUpload"
	containers.Set(key, agent.UploadErrorFile{})
	key = consts.ValidatorPrefix + "HardwareMonitor"
	containers.Set(key, agent.Uninstall{})

	//任务相关
	key = consts.ValidatorPrefix + "TaskGet"
	containers.Set(key, task.Get{})
	key = consts.ValidatorPrefix + "TaskClose"
	containers.Set(key, task.Close{})
}
