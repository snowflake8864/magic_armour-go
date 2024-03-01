package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/agent/agent"
	"goskeleton/app/http/validator/agent/alarm"
	"goskeleton/app/http/validator/agent/auth"
	"goskeleton/app/http/validator/agent/blackmail"
	"goskeleton/app/http/validator/agent/dll"
	"goskeleton/app/http/validator/agent/plug"
	"goskeleton/app/http/validator/agent/port"
	"goskeleton/app/http/validator/agent/proc_sample"
	"goskeleton/app/http/validator/agent/process"
	Syslog "goskeleton/app/http/validator/agent/syslog"
	"goskeleton/app/http/validator/agent/tamper"
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

	//进程相关
	key = consts.ValidatorPrefix + "ProcessAllPut"
	containers.Set(key, process.AllPut{})
	key = consts.ValidatorPrefix + "ProcessRunningPut"
	containers.Set(key, process.RunningPut{})
	key = consts.ValidatorPrefix + "ProcessAutoPut"
	containers.Set(key, process.AutoPut{})
	key = consts.ValidatorPrefix + "ProcessWhiteGet"
	containers.Set(key, process.WhiteGet{})
	key = consts.ValidatorPrefix + "ProcessBlackGet"
	containers.Set(key, process.BlackGet{})
	key = consts.ValidatorPrefix + "ProcessAssignGet"
	containers.Set(key, process.AssignGet{})
	key = consts.ValidatorPrefix + "ProcessAssignPut"
	containers.Set(key, process.AssignPut{})

	//模块相关
	key = consts.ValidatorPrefix + "DllAllPut"
	containers.Set(key, dll.AllPut{})
	key = consts.ValidatorPrefix + "DllRunningPut"
	containers.Set(key, dll.RunningPut{})
	key = consts.ValidatorPrefix + "DllAutoPut"
	containers.Set(key, dll.AutoPut{})
	key = consts.ValidatorPrefix + "DllWhiteGet"
	containers.Set(key, dll.WhiteGet{})
	key = consts.ValidatorPrefix + "DllBlackGet"
	containers.Set(key, dll.BlackGet{})

	//进程或模块样本
	key = consts.ValidatorPrefix + "ProcSampleGet"
	containers.Set(key, proc_sample.Get{})
	key = consts.ValidatorPrefix + "ProcSampleUpload"
	containers.Set(key, proc_sample.Upload{})

	//防篡改相关
	key = consts.ValidatorPrefix + "TamperDirGet"
	containers.Set(key, tamper.DirGet{})
	key = consts.ValidatorPrefix + "TamperDirPut"
	containers.Set(key, tamper.DirPut{})
	key = consts.ValidatorPrefix + "TamperPolicyGet"
	containers.Set(key, tamper.PolicyGet{})

	//勒索相关
	key = consts.ValidatorPrefix + "BlackmailPolicyGet"
	containers.Set(key, blackmail.PolicyGet{})

	//端口相关
	key = consts.ValidatorPrefix + "PortWorkingPut"
	containers.Set(key, port.WorkingPut{})
	key = consts.ValidatorPrefix + "PortVirtualGet"
	containers.Set(key, port.VirtualGet{})

	// 封堵相关
	key = consts.ValidatorPrefix + "PlugDynamicGet"
	containers.Set(key, plug.DynamicGet{})
	key = consts.ValidatorPrefix + "PlugDynamicPut"
	containers.Set(key, plug.DynamicPut{})
	key = consts.ValidatorPrefix + "PlugIpBlackGet"
	containers.Set(key, plug.IpBlackGet{})
	key = consts.ValidatorPrefix + "PlugIpWhiteGet"
	containers.Set(key, plug.IpWhiteGet{})

	//外设相关
	key = consts.ValidatorPrefix + "UsbPut"
	containers.Set(key, task.Close{})
	key = consts.ValidatorPrefix + "UsbBlackGet"
	containers.Set(key, task.Close{})
	key = consts.ValidatorPrefix + "UsbWhiteGet"
	containers.Set(key, task.Close{})

	//告警相关
	key = consts.ValidatorPrefix + "AlarmMasterPut"
	containers.Set(key, alarm.MasterPut{})
	key = consts.ValidatorPrefix + "AlarmPortPut"
	containers.Set(key, alarm.PortPut{})

	//通讯日志相关
	key = consts.ValidatorPrefix + "SyslogPut"
	containers.Set(key, Syslog.Put{})
}
