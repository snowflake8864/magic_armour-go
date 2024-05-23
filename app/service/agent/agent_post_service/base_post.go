package agent_post_service

import "encoding/json"

const (
	ProcTaskAll     string = "all"
	ProcTaskAuto    string = "auto"
	ProcTaskRefresh string = "refresh"
)

// 进程模块数据
type procData struct {
	ProcId       int       `json:"id"`
	Dir          string    `json:"dir"`
	User         string    `json:"user"`
	Hash         string    `json:"hash"`
	Introduce    string    `json:"introduce"`
	Copyright    string    `json:"copyright"`
	Attribute    string    `json:"attribute"`
	ModuleNumber int       `json:"module_number"`
	ModuleList   []dllData `json:"module"`
	Signature    string    `json:"signature"` //不在使用 @todo bug6830
	Finish       int       `json:"finish"`    //不在使用 @todo 100代表任务已完成
}

type dllData struct {
	Dir       string `json:"name"`
	Hash      string `json:"hash"`
	Attribute string `json:"attribute"`
}

// 告警：除端口虚开告警外
type alarmComm struct {
	Type             int    `json:"type"`
	AgentTime        int64  `json:"time"`
	Level            int    `json:"level"`
	Hash             string `json:"hash"`
	Dir              string `json:"dir"`
	PParam           string `json:"p_param"`
	ExceptionProcess string `json:"exception_process"`
	Pid              int    `json:"id"`

	ProcDir   string `json:"proc_dir"`
	ProcHash  string `json:"proc_hash"`
	ProcParam string `json:"proc_param"`
	FileDir   string `json:"file_dir"`
	TargetDir string `json:"target_dir"`

	PeripheralEid    string `json:"peripheral_eid"`
	PeripheralName   string `json:"peripheral_name"`
	PeripheralRemark string `json:"peripheral_remark"`
}

// 告警：端口虚开
type alarmPort struct {
	Weight        int    `json:"weight"`
	Time          int64  `json:"time"`
	AttackIp      string `json:"attack_ip"`
	DestinationIp string `json:"destination_ip"`
	OpenPort      int    `json:"open_port"`
	RedirectIp    string `json:"redirect_ip"`
	RedirectPort  int    `json:"redirect_port"`
}

// 通讯日志： 公共
type syslogComm struct {
	LogType       json.Number `json:"log_type"`
	Time          json.Number `json:"time"`
	PId           json.Number `json:"p_id"`
	PDir          string      `json:"p_dir"`
	DomainName    string      `json:"domain_name"`
	Hash          string      `json:"hash"`
	Uid           string      `json:"uid"`
	SourceIp      string      `json:"source_ip"`
	SourcePort    json.Number `json:"source_port"`
	ResIp         string      `json:"res_ip"`
	RsPort        json.Number `json:"rs_port"` // 注意字段
	ResPortStatus json.Number `json:"res_port_status"`
	PParam        string      `json:"p_param"`
	PpHash        string      `json:"pp_hash"`
	PpId          json.Number `json:"pp_id"`
	PpDir         string      `json:"pp_dir"`
	PpParam       string      `json:"pp_param"`
	Proto         json.Number `json:"proto"`
	Ip            string      `json:"ip"`
	Username      string      `json:"username"`
	Type          string      `json:"type"`
	Status        json.Number `json:"status"`
}

// 端口： 业务端口
type portWorking struct {
	Type    string      `json:"type"`
	Ip      string      `json:"ip"`
	Port    json.Number `json:"port"`
	Process string      `json:"process"`
	Pid     json.Number `json:"pid"` //linux没有
}

// 终端监控
type agentMonitor struct {
	HdSize  string      `json:"hd_size"`
	HdUsage json.Number `json:"hd_usage"`

	CpuNumber string            `json:"cpu_number"`
	CpuUsage  json.Number       `json:"cpu_usage"`
	CpuTops   []agentMonitorTop `json:"cpu_tops"`

	MemSize  string            `json:"mem_size"`
	MemUsage json.Number       `json:"mem_usage"`
	MemTops  []agentMonitorTop `json:"mem_tops"`
	Self     agentMonitorTop   `json:"self"`
}

// 终端监控top列表
type agentMonitorTop struct {
	Id       json.Number `json:"id"`
	Dir      string      `json:"dir"`
	Hash     string      `json:"hash"`
	CpuUsage json.Number `json:"cpu_usage"`
	MemSize  string      `json:"mem_size"`
	User     string      `json:"user"`
}

// 外设
type usbData struct {
	PeripheralEid   string `json:"peripheral_eid"`
	PeripheralName  string `json:"peripheral_name"`
	PeripheralIntro string `json:"peripheral_intro"`
	PeripheralType  string `json:"peripheral_type"`
}
