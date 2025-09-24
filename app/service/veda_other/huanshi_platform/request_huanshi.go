package huanshi_platform

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"net"
	"strings"
)

const (
	AlarmTypeProcess   int = 1  //告警类型: 进程
	AlarmTypeDll       int = 2  //告警类型: 模块
	AlarmTypeTamper    int = 3  //告警类型: 防篡改
	AlarmTypeBlackmail int = 4  //告警类型: 勒索
	AlarmTypePort      int = 5  //告警类型: 端口虚开
	AlarmTypeUsb       int = 6  //告警类型: 外设
	AlarmTypePlug      int = 7  //告警类型: 动态封堵
	AlarmTypeOutreach  int = 8  //告警类型: 外联
	AlarmTypeLogin     int = 9  //告警类型: 口令破解
	AlarmTypeSelf      int = 10 //告警类型: 自保

	DataTypeAgent string = "zhongduan" //数据类型: 终端
	ApiVer        int    = 2           //回调api版本

	LogTypeDns        int = 1 //日志类型：dns
	LogTypeOuter      int = 2 //日志类型：对外通讯
	LogTypeInner      int = 3 //日志类型：对内通讯
	LogTypeProcess    int = 4 //日志类型：进程启动
	LogTypeLogin      int = 5 //日志类型：登录
	LogTypeOrgOperate int = 6 //日志类型：企业操作

)

func CreateRequestFactory(c RequestConfig) *RequestHuanshi {
	c.HuanshiUrl = strings.TrimRight(c.HuanshiUrl, "/")
	return &RequestHuanshi{
		redisClient:      redis_factory.GetOneRedisClient(),
		cacheKey:         "huanshi_token",
		cacheExpire:      600,
		Config:           c,
		PrivateRsa:       variable.HuanshiPrivateRsa,
		RequestPlatform:  "huanjia",
		Edition:          variable.VersionNumber,
		agentCachePrefix: "huanshi_agent_",
	}
}

type RequestConfig struct {
	HuanshiSwitch      int
	HuanshiCallbackUrl string
	HuanshiCallIp      string
	HuanshiUrl         string
	HuanshiAlarmSwitch int
	HuanshiAlarmIp     string
	HuanshiAlarmPort   string
	HuanshiLogSwitch   int
	HuanshiLogIp       string
	HuanshiLogPort     string
	HuanshiThirdSwitch int
	HuanshiThirdIp     string
	HuanshiThirdPort   string
}

type RequestHuanshi struct {
	redisClient      *redis_factory.RedisClient
	cacheKey         string
	cacheExpire      int //缓存有效期
	Config           RequestConfig
	PrivateRsa       string
	RequestPlatform  string
	Edition          string
	agentCachePrefix string
}

type responseCommon struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
}

// checkCode 校验code
func (aa *RequestHuanshi) checkCode(bodyContext string, err error) error {
	if err != nil {
		return errors.New("body_error: " + err.Error())
	}
	resp := responseCommon{}
	err = json_helper.ToStruct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error())
	}

	//fmt.Println("Type of x:", reflect.TypeOf(resp.Code))
	code := float_helper.Float64(resp.Code)
	//口令过期
	if code == 5004 {
		return errors.New("token fail: " + bodyContext)
	}
	//结果不正确
	if code != 1 {
		return errors.New("result fail: " + bodyContext)
	}
	return nil
}

// Online 是否在线
func (aa *RequestHuanshi) Online() error {
	if aa.Config.HuanshiSwitch == enums.SwitchOff {
		return nil
	}
	uri := aa.Config.HuanshiUrl + "/web/heartbeat"

	bodyContext, err := http_request.GetNoParams(uri, map[string]interface{}{})
	if err1 := aa.checkCode(bodyContext, err); err1 != nil {
		return err1
	}

	return nil
}

// ResponseToken 生成口令接口返回数据
type ResponseToken struct {
	Data struct {
		AccessToken string `json:"token"`
	} `json:"data"`
}

// GenerateToken 生成口令
func (aa *RequestHuanshi) GenerateToken() error {
	if aa.Config.HuanshiSwitch == enums.SwitchOff {
		return nil
	}

	uri := aa.Config.HuanshiUrl + "/web/site_auth/login"
	Timestamp := date_helper.CurTimestamp()
	//  $param['rsa'] = md5( $param['url'].'_'.$param['time'].'_'.$hs_private_rsa);
	params := map[string]interface{}{
		"url":  aa.Config.HuanshiCallbackUrl,
		"time": Timestamp,
		"rsa":  md5_encrypt.MD5(aa.Config.HuanshiCallbackUrl + "_" + int_helper.Int64ToStr(Timestamp) + "_" + aa.PrivateRsa),
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err1 := aa.checkCode(bodyContext, err); err1 != nil {
		return err1
	}

	data := ResponseToken{}
	err = json_helper.ToStruct(bodyContext, &data)
	if err != nil {
		return errors.New("json_error2: " + err.Error())
	}

	//新增缓存
	if isOK, _ := aa.redisClient.SetEx(aa.cacheKey, data.Data.AccessToken, aa.cacheExpire); !isOK {
		return errors.New("cache token fail")
	}
	return nil
}

func (aa *RequestHuanshi) GetHeaders() (headers map[string]interface{}) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	if token := aa.redisClient.Get(aa.cacheKey); !string_helper.IsEmpty(token) {
		headers["Authorization"] = "Bearer " + token
	}
	return headers
}

type DeviceData struct {
	Edition  string `json:"edition"`  //产品版本号
	Platform string `json:"platform"` //对应平台
	IP       string `json:"ip"`       //回调IP
	Url      string `json:"url"`      //回调地址
	ApiVer   int    `json:"apiVer"`   //回调接口版本

	DeviceMark string `json:"deviceMark"` //设备标识
	DeviceRoot string `json:"deviceRoot"` //设备管理员

	AgentNum       int `json:"agentNum"`       //终端总数
	OnlineAgentNum int `json:"onlineAgentNum"` //终端在线数
}

// SyncDevices 同步设备
func (aa *RequestHuanshi) SyncDevices(params *[]DeviceData) error {
	if aa.Config.HuanshiSwitch == enums.SwitchOff {
		return nil
	}
	if err := aa.Online(); err != nil {
		return errors.New("网络连接： " + err.Error())
	}
	if err := aa.GenerateToken(); err != nil {
		return errors.New("生成口令： " + err.Error())
	}

	uri := aa.Config.HuanshiUrl + "/web/device/add"

	for _, param := range *params {
		param.Platform = aa.RequestPlatform
		param.Edition = aa.Edition
		param.IP = aa.Config.HuanshiCallIp
		param.Url = aa.Config.HuanshiCallbackUrl
		param.ApiVer = ApiVer
		bodyContext, err := http_request.PostWithJson(uri, aa.GetHeaders(), param)
		if err1 := aa.checkCode(bodyContext, err); err1 != nil {
			return errors.New("同步设备： " + err1.Error())
		}
	}

	return nil
}

// SYSLOG发送
func (aa *RequestHuanshi) sendBySyslog(tt string, paramJson string) error {
	var facility, level int
	if tt == "agent" {
		//@todo 通道:LOCAL7; 等级:info
		facility = 23
		level = 6
	} else if tt == "alarm" {
		//@todo 通道:LOCAL5; 等级:info
		facility = 21
		level = 5
	} else if tt == "log" {
		//@todo 通道:LOCAL3; 等级:warning
		facility = 19
		level = 4
	} else {
		return errors.New("type not in [agent alarm log]")
	}
	//地址
	var addresses []string

	//幻势平台
	if aa.Config.HuanshiSwitch == enums.SwitchOn {
		if tt == "agent" {
			addresses = append(addresses, aa.Config.HuanshiAlarmIp+":"+aa.Config.HuanshiAlarmPort)
		} else if tt == "alarm" && aa.Config.HuanshiAlarmSwitch == enums.SwitchOn {
			addresses = append(addresses, aa.Config.HuanshiAlarmIp+":"+aa.Config.HuanshiAlarmPort)
		} else if tt == "log" && aa.Config.HuanshiLogSwitch == enums.SwitchOn {
			addresses = append(addresses, aa.Config.HuanshiLogIp+":"+aa.Config.HuanshiLogPort)
		}
	}
	//友商平台
	if aa.Config.HuanshiThirdSwitch == enums.SwitchOn {
		addresses = append(addresses, aa.Config.HuanshiThirdIp+":"+aa.Config.HuanshiThirdPort)
	}
	if len(addresses) > 0 {
		for _, addr := range addresses {
			if err1 := aa.analogSyslog(addr, facility, level, paramJson); err1 != nil {
				variable.ZapLog.Sugar().Error("analogSyslog: " + err1.Error())
			}
		}
	}

	return nil
}

// 模拟 Syslog
func (aa *RequestHuanshi) analogSyslog(address string, facility, level int, paramJson string) error {
	//请求地址
	conn, err := net.Dial("udp", address)
	if err != nil {
		return errors.New("net.Dial err:" + err.Error())
	}
	defer func(conn net.Conn) {
		if err1 := conn.Close(); err1 != nil {
			variable.ZapLog.Sugar().Error("net.Close err:" + err.Error())
		}
	}(conn)

	//模拟SYSLOG信息
	msg := "<"
	msg += int_helper.IntToStr((facility * 8) + level)
	msg += ">"
	msg += paramJson
	//variable.ZapLog.Sugar().Debug("============================", address)
	//variable.ZapLog.Sugar().Debug(msg)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return errors.New("net.Write err:" + err.Error())
	}
	return nil
}

type AgentData struct {
	DeviceMark   string  `json:"deviceMark"`   //设备标识
	Platform     string  `json:"platform"`     //来源
	DataType     string  `json:"dataType"`     //数据类型
	OrgId        int64   `json:"orgId"`        //企业编号
	AgentId      string  `json:"agentId"`      //终端编号
	Os           string  `json:"os"`           //系统
	Ip           string  `json:"ip"`           //终端ip
	Mac          string  `json:"mac"`          //终端Mac
	NetIp        string  `json:"netIp"`        //通讯ip
	HostName     string  `json:"hostName"`     //计算机名称
	Status       int     `json:"status"`       //终端上线状态
	OnlineTime   string  `json:"onlineTime"`   //终端上线时间
	Version      string  `json:"version"`      //终端版本
	DeptCode     string  `json:"deptCode"`     //组织架构
	DeptName     string  `json:"deptName"`     //组织架构
	AgentType    string  `json:"agentType"`    //终端类型
	AliasName    string  `json:"aliasName"`    //别名
	ContactName  string  `json:"contactName"`  //管理人
	ContactPhone string  `json:"contactPhone"` //管理人
	CpuInfo      string  `json:"cpuInfo"`      //终端硬件信息
	MemoryInfo   string  `json:"memoryInfo"`   //终端硬件信息
	DiskInfo     string  `json:"diskInfo"`     //终端硬件信息
	NetworkCard  string  `json:"networkCard"`  //终端硬件信息
	CpuLoad      float64 `json:"cpuLoad"`      //终端在线时，最近资源监控数据
	MemoryLoad   float64 `json:"memoryLoad"`   //终端在线时，最近资源监控数据
	DiskLoad     float64 `json:"diskLoad"`     //终端在线时，最近资源监控数据

}

// SyncAgentBySyslog 通过SYSLOG同步终端
func (aa *RequestHuanshi) SyncAgentBySyslog(param *AgentData) error {
	//默认值
	param.Platform = aa.RequestPlatform
	param.DataType = DataTypeAgent
	if err := aa.sendBySyslog("agent", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type AlarmBase struct {
	DeviceMark  string `json:"deviceMark"`  //设备标识
	Platform    string `json:"platform"`    //来源
	Type        int    `json:"type"`        //类型模块
	OrgId       int64  `json:"orgId"`       //企业编号
	AgentId     string `json:"uid"`         //终端编号
	AgentIp     string `json:"ip"`          //终端ip
	ContactName string `json:"contactName"` //管理人
	DeptName    string `json:"deptName"`    //组织架构
	Level       int    `json:"level"`       //告警级别
	NoticeType  int    `json:"noticeType"`  //告警类型
	NoticeTime  string `json:"noticeTime"`  //告警时间
}

// AlarmComm 告警数据结构： 公共
type AlarmComm struct {
	AlarmBase
	Dir               string `json:"dir"`               //告警关联路径
	ExceptionProcess  string `json:"exceptionProcess"`  ///异常进程路径
	AttackIp          string `json:"attackIp"`          //虚开端口: 攻击ip
	DestinationIp     string `json:"destinationIp"`     //虚开端口: 目的ip
	OpenPort          int    `json:"openPort"`          //虚开端口: 目的端口号
	RedirectIp        string `json:"redirectIp"`        //虚开端口: 重定向ip
	PeripheralEid     string `json:"peripheralEid"`     //外设:  外设ID
	PeripheralName    string `json:"peripheralName"`    //外设:  设备名称
	PeripheralRemark  string `json:"peripheralRemark"`  //外设: 设备描述
	PeripheralContact string `json:"peripheralContact"` //外设: 责任人
}

// AlarmPlug 告警数据结构： 封堵
type AlarmPlug struct {
	AlarmBase
	Pip        string `json:"pip"`        //封堵 ip
	PStartTime string `json:"pStartTime"` //封堵 开始时间
	PEndTime   string `json:"pEndTime"`   //封堵 结束时间
	PMinute    int    `json:"pMinute"`    //封堵 时长（分钟）
	PType      int    `json:"pType"`      //封堵 类型
	PTypeTxt   string `json:"pTypeTxt"`   //封堵 类型说明
}

// AlarmOutreach 告警数据结构： 外联
type AlarmOutreach struct {
	AlarmBase
	Phash      string `json:"phash"`      //进程 HASH
	Pid        int    `json:"pid"`        //进程 ID
	Ppath      string `json:"ppath"`      //进程 路径
	Domain     string `json:"domain"`     //域名
	AnalysisIp string `json:"analysisIp"` //解析ip
}

// AlarmLogin 告警数据结构：  登录口令破解
type AlarmLogin struct {
	AlarmBase
	SourceIp string `json:"sourceIp"` //来源ip
	Total    int    `json:"total"`    //次数
}

// SyncAlarmComm 同步告警: 公共
func (aa *RequestHuanshi) SyncAlarmComm(param AlarmComm) error {
	param.Platform = aa.RequestPlatform
	//variable.ZapLog.Sugar().Debug(json_helper.Encode(param))
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.sendBySyslog("alarm", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncAlarmOutreach 同步告警: 外联
func (aa *RequestHuanshi) SyncAlarmOutreach(param AlarmOutreach) error {
	param.Platform = aa.RequestPlatform
	param.Type = AlarmTypeOutreach
	//variable.ZapLog.Sugar().Debug(json_helper.Encode(param))
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.sendBySyslog("alarm", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncAlarmLogin 同步告警: 登录
func (aa *RequestHuanshi) SyncAlarmLogin(param AlarmLogin) error {
	param.Platform = aa.RequestPlatform
	param.Type = AlarmTypeLogin
	//variable.ZapLog.Sugar().Debug(json_helper.Encode(param))
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.sendBySyslog("alarm", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncAlarmPlug 同步告警: 封堵
func (aa *RequestHuanshi) SyncAlarmPlug(param AlarmPlug) error {
	param.Platform = aa.RequestPlatform
	param.Type = AlarmTypePlug
	//variable.ZapLog.Sugar().Debug(json_helper.Encode(param))
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.sendBySyslog("alarm", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type LogBase struct {
	DeviceMark string `json:"device_mark"` //设备标识
	DeviceIp   string `json:"device_ip"`   //设备IP
	Platform   string `json:"platform"`    //来源
	LogType    int    `json:"log_type"`    //日志类型
	OrgId      int64  `json:"org_id"`      //企业编号
	AgentId    string `json:"agent_id"`    //终端标识
	AgentIp    string `json:"agent_ip"`    //终端ip
	CreatedAt  string `json:"created_at"`  //终端ip
}

// LogProcess 日志数据结构： 进程启动日志
type LogProcess struct {
	LogBase
	ProcHash    string `json:"proc_hash"`
	ProcId      int    `json:"proc_id"`
	ProcDir     string `json:"proc_dir"`
	ProcParam   string `json:"proc_param"`
	ParentHash  string `json:"parent_hash"`
	ParentId    int    `json:"parent_id"`
	ParentDir   string `json:"parent_dir"`
	ParentParam string `json:"parent_param"`
}

// LogDns 日志数据结构： DNS日志
type LogDns struct {
	LogBase
	ProcHash   string `json:"proc_hash"`
	ProcId     int    `json:"proc_id"`
	ProcDir    string `json:"proc_dir"`
	DomainName string `json:"domain_name"`
	TargetIp   string `json:"target_ip"`
}

// LogNetwork 日志数据结构： 通讯日志
type LogNetwork struct {
	LogBase
	ProcHash         string `json:"proc_hash"`
	ProcId           int    `json:"proc_id"`
	ProcDir          string `json:"proc_dir"`
	TargetIp         string `json:"target_ip"`
	TargetPort       int    `json:"target_port"`
	TargetPortStatus int    `json:"target_port_status"`
	SourceIp         string `json:"source_ip"`
	SourcePort       int    `json:"source_port"`
	Proto            int    `json:"proto"`
	ProtoLabel       string `json:"proto_label"`
}

// LogLogin 日志数据结构： 登录日志
type LogLogin struct {
	LogBase
	SourceIp string `json:"source_ip"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Status   int    `json:"status"`
}

// SyncLogDns 同步日志：DNS
func (aa *RequestHuanshi) SyncLogDns(param LogDns) error {
	param.DeviceIp = aa.Config.HuanshiCallIp
	param.Platform = aa.RequestPlatform
	param.LogType = LogTypeDns
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
	}
	if err := aa.sendBySyslog("log", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncLogNetwork 同步日志：通讯
func (aa *RequestHuanshi) SyncLogNetwork(param LogNetwork) error {
	param.DeviceIp = aa.Config.HuanshiCallIp
	param.Platform = aa.RequestPlatform
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
	}
	if err := aa.sendBySyslog("log", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncLogProcess 同步日志：进程
func (aa *RequestHuanshi) SyncLogProcess(param LogProcess) error {
	param.DeviceIp = aa.Config.HuanshiCallIp
	param.Platform = aa.RequestPlatform
	param.LogType = LogTypeProcess
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
	}
	if err := aa.sendBySyslog("log", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

// SyncLogLogin 同步日志：登录
func (aa *RequestHuanshi) SyncLogLogin(param LogLogin) error {
	param.DeviceIp = aa.Config.HuanshiCallIp
	param.Platform = aa.RequestPlatform
	param.LogType = LogTypeLogin
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
	}
	if err := aa.sendBySyslog("log", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type CacheAgent struct {
	OrgCode     string `json:"org_code"`
	AgentIp     string `json:"agent_ip"`
	NetIp       string `json:"net_ip"`
	DeptName    string `json:"dept_name"`
	ContactName string `json:"contact_name"`
	AgentType   string `json:"agent_type"`
}

func (aa *RequestHuanshi) cacheAgentKey(orgId int64, agentId string) string {

	return aa.agentCachePrefix + int_helper.Int64ToStr(orgId) + "_" + agentId
}

// SaveAgent 保存终端信息
func (aa *RequestHuanshi) SaveAgent(orgId int64, agentId string, param *CacheAgent) error {
	key := aa.cacheAgentKey(orgId, agentId)
	data := json_helper.Encode(param)
	//新增缓存
	if isOK, err := aa.redisClient.SetEx(key, data, aa.cacheExpire); !isOK {
		return errors.New("cache agent fail: " + err.Error())
	}

	return nil
}

// GetAgent 获取终端信息
func (aa *RequestHuanshi) GetAgent(orgId int64, agentId string) (err error, data *CacheAgent) {

	jsonStr := aa.redisClient.Get(aa.cacheAgentKey(orgId, agentId))
	if string_helper.IsEmpty(jsonStr) {
		return errors.New("缓存数据不存在"), data
	}

	err = json_helper.ToStruct(jsonStr, &data)
	if err != nil {
		return errors.New("缓存数据错误"), data
	}
	return nil, data
}
