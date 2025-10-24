package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	agent1 "goskeleton/app/http/validator/agent/agent"
	agentSyslog1 "goskeleton/app/http/validator/agent/agent_syslog"
	alarm1 "goskeleton/app/http/validator/agent/alarm"
	auth1 "goskeleton/app/http/validator/agent/auth"
	backup1 "goskeleton/app/http/validator/agent/backup"
	dll1 "goskeleton/app/http/validator/agent/dll"
	jump1 "goskeleton/app/http/validator/agent/jump"
	plug1 "goskeleton/app/http/validator/agent/plug"
	port1 "goskeleton/app/http/validator/agent/port"
	procSample1 "goskeleton/app/http/validator/agent/proc_sample"
	process1 "goskeleton/app/http/validator/agent/process"
	tamper1 "goskeleton/app/http/validator/agent/tamper"
	task1 "goskeleton/app/http/validator/agent/task"
	usb1 "goskeleton/app/http/validator/agent/usb"
)

// AgentRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func AgentRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	//登录验证
	key = consts.ValidatorPrefixAgent + "AuthMakeUserid"
	containers.Set(key, auth1.MakeUserid{})
	key = consts.ValidatorPrefixAgent + "AuthLogin"
	containers.Set(key, auth1.Login{})

	//终端相关
	key = consts.ValidatorPrefixAgent + "ErrorFileUpload"
	containers.Set(key, agent1.UploadErrorFile{})
	key = consts.ValidatorPrefixAgent + "HardwareMonitor"
	containers.Set(key, agent1.MonitorHardware{})
	key = consts.ValidatorPrefixAgent + "IpMonitor"
	containers.Set(key, agent1.MonitorIp{})

	//任务相关
	key = consts.ValidatorPrefixAgent + "TaskClose"
	containers.Set(key, task1.Close{})

	//进程相关
	key = consts.ValidatorPrefixAgent + "ProcessAllPut"
	containers.Set(key, process1.AllPut{})
	key = consts.ValidatorPrefixAgent + "ProcessRunningPut"
	containers.Set(key, process1.RunningPut{})
	key = consts.ValidatorPrefixAgent + "ProcessAutoPut"
	containers.Set(key, process1.AutoPut{})
	key = consts.ValidatorPrefixAgent + "ProcessAssignPut"
	containers.Set(key, process1.AssignPut{})

	//模块相关
	key = consts.ValidatorPrefixAgent + "DllAllPut"
	containers.Set(key, dll1.AllPut{})
	key = consts.ValidatorPrefixAgent + "DllRunningPut"
	containers.Set(key, dll1.RunningPut{})
	key = consts.ValidatorPrefixAgent + "DllAutoPut"
	containers.Set(key, dll1.AutoPut{})

	//进程或模块样本
	key = consts.ValidatorPrefixAgent + "ProcSampleUpload"
	containers.Set(key, procSample1.Upload{})

	//防篡改相关
	key = consts.ValidatorPrefixAgent + "TamperDirPut"
	containers.Set(key, tamper1.DirPut{})

	//勒索相关

	//端口相关
	key = consts.ValidatorPrefixAgent + "PortWorkingPut"
	containers.Set(key, port1.WorkingPut{})

	// 封堵相关
	key = consts.ValidatorPrefixAgent + "PlugDynamicPut"
	containers.Set(key, plug1.DynamicPut{})

	//外设相关
	key = consts.ValidatorPrefixAgent + "UsbPut"
	containers.Set(key, usb1.Put{})

	//告警相关
	key = consts.ValidatorPrefixAgent + "AlarmMasterPut"
	containers.Set(key, alarm1.MasterPut{})
	key = consts.ValidatorPrefixAgent + "AlarmPortPut"
	containers.Set(key, alarm1.PortPut{})

	//通讯日志相关
	key = consts.ValidatorPrefixAgent + "SyslogPut"
	containers.Set(key, agentSyslog1.Put{})

	//备份相关
	key = consts.ValidatorPrefixAgent + "BackupPut"
	containers.Set(key, backup1.PutBackup{})
	key = consts.ValidatorPrefixAgent + "RollbackPut"
	containers.Set(key, backup1.PutRollback{})

	//跳变
	key = consts.ValidatorPrefixAgent + "PwJumpPut"
	containers.Set(key, jump1.PwJumpPut{})
	key = consts.ValidatorPrefixAgent + "IpJumpPut"
	containers.Set(key, jump1.IpJumpPut{})
}
