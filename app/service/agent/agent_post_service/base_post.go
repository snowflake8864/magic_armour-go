package agent_post_service

const (
	ProcTaskAll     string = "all"
	ProcTaskAuto    string = "auto"
	ProcTaskRefresh string = "refresh"
)

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

type alarmComm struct {
	Type             int    `json:"type"`
	AgentTime        int64  `json:"time"`
	Level            int    `json:"level"`
	Hash             string `json:"hash"`
	Dir              string `json:"dir"`
	PParam           string `json:"p_param"`
	ExceptionProcess string `json:"exception_process"`

	ProcDir   string `json:"proc_dir"`
	ProcHash  string `json:"proc_hash"`
	ProcParam string `json:"proc_param"`
	FileDir   string `json:"file_dir"`
	TargetDir string `json:"target_dir"`

	PeripheralEid    string `json:"peripheral_eid"`
	PeripheralName   string `json:"peripheral_name"`
	PeripheralRemark string `json:"peripheral_remark"`
}

type alarmPort struct {
	Weight        int    `json:"weight"`
	Time          int64  `json:"time"`
	AttackIp      string `json:"attack_ip"`
	DestinationIp string `json:"destination_ip"`
	OpenPort      int    `json:"open_port"`
	RedirectIp    string `json:"redirect_ip"`
	RedirectPort  int    `json:"redirect_port"`
}
