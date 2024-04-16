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
	"goskeleton/app/http/validator/agent/usb"
)

// AgentRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func AgentRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//登录验证
	key = consts.ValidatorPrefixAgent + "AuthMakeUserid"
	containers.Set(key, auth.MakeUserid{})
	key = consts.ValidatorPrefixAgent + "AuthLogin"
	containers.Set(key, auth.Login{})

	//终端相关
	key = consts.ValidatorPrefixAgent + "AgentConfigGet"
	containers.Set(key, agent.ConfigGet{})
	key = consts.ValidatorPrefixAgent + "AgentInfoGet"
	containers.Set(key, agent.InfoGet{})
	key = consts.ValidatorPrefixAgent + "AgentUpgrade"
	containers.Set(key, agent.Upgrade{})
	key = consts.ValidatorPrefixAgent + "AgentUninstall"
	containers.Set(key, agent.Uninstall{})
	key = consts.ValidatorPrefixAgent + "TrustDirGet"
	containers.Set(key, agent.GetTrustDir{})
	key = consts.ValidatorPrefixAgent + "ErrorFileUpload"
	containers.Set(key, agent.UploadErrorFile{})
	key = consts.ValidatorPrefixAgent + "HardwareMonitor"
	containers.Set(key, agent.Uninstall{})

	//任务相关
	key = consts.ValidatorPrefixAgent + "TaskGet"
	containers.Set(key, task.Get{})
	key = consts.ValidatorPrefixAgent + "TaskClose"
	containers.Set(key, task.Close{})

	//进程相关
	key = consts.ValidatorPrefixAgent + "ProcessAllPut"
	containers.Set(key, process.AllPut{})
	key = consts.ValidatorPrefixAgent + "ProcessRunningPut"
	containers.Set(key, process.RunningPut{})
	key = consts.ValidatorPrefixAgent + "ProcessAutoPut"
	containers.Set(key, process.AutoPut{})
	key = consts.ValidatorPrefixAgent + "ProcessWhiteGet"
	containers.Set(key, process.WhiteGet{})
	key = consts.ValidatorPrefixAgent + "ProcessBlackGet"
	containers.Set(key, process.BlackGet{})
	key = consts.ValidatorPrefixAgent + "ProcessAssignGet"
	containers.Set(key, process.AssignGet{})
	key = consts.ValidatorPrefixAgent + "ProcessAssignPut"
	containers.Set(key, process.AssignPut{})

	//模块相关
	key = consts.ValidatorPrefixAgent + "DllAllPut"
	containers.Set(key, dll.AllPut{})
	key = consts.ValidatorPrefixAgent + "DllRunningPut"
	containers.Set(key, dll.RunningPut{})
	key = consts.ValidatorPrefixAgent + "DllAutoPut"
	containers.Set(key, dll.AutoPut{})
	key = consts.ValidatorPrefixAgent + "DllWhiteGet"
	containers.Set(key, dll.WhiteGet{})
	key = consts.ValidatorPrefixAgent + "DllBlackGet"
	containers.Set(key, dll.BlackGet{})

	//进程或模块样本
	key = consts.ValidatorPrefixAgent + "ProcSampleGet"
	containers.Set(key, proc_sample.Get{})
	key = consts.ValidatorPrefixAgent + "ProcSampleUpload"
	containers.Set(key, proc_sample.Upload{})

	//防篡改相关
	key = consts.ValidatorPrefixAgent + "TamperDirGet"
	containers.Set(key, tamper.DirGet{})
	key = consts.ValidatorPrefixAgent + "TamperDirPut"
	containers.Set(key, tamper.DirPut{})
	key = consts.ValidatorPrefixAgent + "TamperPolicyGet"
	containers.Set(key, tamper.PolicyGet{})

	//勒索相关
	key = consts.ValidatorPrefixAgent + "BlackmailPolicyGet"
	containers.Set(key, blackmail.PolicyGet{})

	//端口相关
	key = consts.ValidatorPrefixAgent + "PortWorkingPut"
	containers.Set(key, port.WorkingPut{})
	key = consts.ValidatorPrefixAgent + "PortVirtualGet"
	containers.Set(key, port.VirtualGet{})

	// 封堵相关
	key = consts.ValidatorPrefixAgent + "PlugDynamicGet"
	containers.Set(key, plug.DynamicGet{})
	key = consts.ValidatorPrefixAgent + "PlugDynamicPut"
	containers.Set(key, plug.DynamicPut{})
	key = consts.ValidatorPrefixAgent + "PlugIpBlackGet"
	containers.Set(key, plug.IpBlackGet{})
	key = consts.ValidatorPrefixAgent + "PlugIpWhiteGet"
	containers.Set(key, plug.IpWhiteGet{})

	//外设相关
	key = consts.ValidatorPrefixAgent + "UsbPut"
	containers.Set(key, usb.Put{})
	key = consts.ValidatorPrefixAgent + "UsbBlackGet"
	containers.Set(key, usb.BlackGet{})
	key = consts.ValidatorPrefixAgent + "UsbWhiteGet"
	containers.Set(key, usb.WhiteGet{})

	//告警相关
	key = consts.ValidatorPrefixAgent + "AlarmMasterPut"
	containers.Set(key, alarm.MasterPut{})
	key = consts.ValidatorPrefixAgent + "AlarmPortPut"
	containers.Set(key, alarm.PortPut{})

	//通讯日志相关
	key = consts.ValidatorPrefixAgent + "SyslogPut"
	containers.Set(key, Syslog.Put{})
}
