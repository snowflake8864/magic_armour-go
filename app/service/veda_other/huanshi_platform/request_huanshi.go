package huanshi_platform

import (
	"errors"
	"goskeleton/app/global/consts"
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
	AlarmTypeBlackmail int = 4  //告警类型: 虚开
	AlarmTypePort      int = 5  //告警类型: 端口
	AlarmTypeUsb       int = 6  //告警类型: 外设
	AlarmTypePlug      int = 7  //告警类型: 动态封堵
	AlarmTypeOutreach  int = 8  //告警类型: 外联
	AlarmTypeLogin     int = 9  //告警类型: 口令破解
	AlarmTypeSelf      int = 10 //告警类型: 自保

	DataTypeAgent string = "zhongduan" //数据类型: 终端

	LogTypeDns     int = 1 //日志类型：dns
	LogTypeOuter   int = 2 //日志类型：对外通讯
	LogTypeInner   int = 3 //日志类型：对内通讯
	LogTypeProcess int = 4 //日志类型：进程启动
	LogTypeLogin   int = 5 //日志类型：登录
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
		Edition:          consts.VersionNumber,
		agentCachePrefix: "huanshi_agent_",
	}
}

type RequestConfig struct {
	HuanshiCallbackUrl string
	HuanshiCallIp      string
	HuanshiUrl         string
	HuanshiAlarmIp     string
	HuanshiAlarmPort   string
	HuanshiLogIp       string
	HuanshiLogPort     string
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
	uri := aa.Config.HuanshiUrl + "/web/heartbeat"

	bodyContext, err := http_request.GetNoParams(uri, map[string]interface{}{})
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
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
	uri := aa.Config.HuanshiUrl + "/web/site_auth/login"
	Timestamp := date_helper.CurTimestamp()
	//  $param['rsa'] = md5( $param['url'].'_'.$param['time'].'_'.$hs_private_rsa);
	params := map[string]interface{}{
		"url":  aa.Config.HuanshiCallbackUrl,
		"time": Timestamp,
		"rsa":  md5_encrypt.MD5(aa.Config.HuanshiCallbackUrl + "_" + int_helper.Int64ToStr(Timestamp) + "_" + aa.PrivateRsa),
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}

	data := ResponseToken{}
	err = json_helper.ToStruct(bodyContext, &data)
	if err != nil {
		return errors.New("json_error2: " + err.Error())
	}

	//新增缓存
	if isOK, err := aa.redisClient.SetEx(aa.cacheKey, data.Data.AccessToken, aa.cacheExpire); !isOK {
		return errors.New("cache token fail: " + err.Error())
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

// SyncDeviceByTcp 通过TCO同步设备
/*
参数：param {
"deviceMark": "" //设备标识
"platform": "" //请求平台
"edition":  "" //产品版本号
"ip": "" //回调IP
"url": "" //回调地址
"agentNum": "" //终端总数
"onlineAgentNum": "" 在线终端数
}
*/
func (aa *RequestHuanshi) SyncDeviceByTcp(param map[string]interface{}) error {
	uri := aa.Config.HuanshiUrl + "/web/device/add"

	//默认值
	param["platform"] = aa.RequestPlatform
	param["edition"] = aa.Edition
	param["ip"] = aa.Config.HuanshiCallIp
	param["url"] = aa.Config.HuanshiCallbackUrl

	bodyContext, err := http_request.PostWithForm(uri, aa.GetHeaders(), param)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}
	return nil
}

// 模拟发送SYSLOG
func (aa *RequestHuanshi) analogSendSyslog(tt string, paramJson string) error {
	//发送地址
	var address string
	var facility, level int
	if tt == "agent" {
		address = aa.Config.HuanshiAlarmIp + ":" + aa.Config.HuanshiAlarmPort
		//@todo 通道:LOCAL7; 等级:info
		facility = 23
		level = 6
	} else if tt == "alarm" {
		address = aa.Config.HuanshiAlarmIp + ":" + aa.Config.HuanshiAlarmPort
		//@todo 通道:LOCAL5; 等级:info
		facility = 21
		level = 5
	} else if tt == "log" {
		address = aa.Config.HuanshiLogIp + ":" + aa.Config.HuanshiLogPort
		//@todo 通道:LOCAL3; 等级:warning
		facility = 19
		level = 4
	} else {
		return errors.New("type not in [agent alarm log]")
	}

	//请求地址
	conn, err := net.Dial("udp", address)
	if err != nil {
		return errors.New("net.Dial err:" + err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			variable.ZapLog.Sugar().Error("net.Close err:" + err.Error())
		}
	}(conn)

	//模拟SYSLOG信息
	msg := "<"
	msg += int_helper.IntToStr((facility * 8) + level)
	msg += ">"
	msg += paramJson

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return errors.New("net.Write err:" + err.Error())
	}
	return nil
}

type AgentData struct {
	DeviceMark   string `json:"deviceMark"`   //设备标识
	Platform     string `json:"platform"`     //来源
	DataType     string `json:"dataType"`     //数据类型
	OrgId        int64  `json:"org_id"`       //企业编号
	AgentId      string `json:"agentId"`      //终端编号
	Os           string `json:"os"`           //系统
	Ip           string `json:"ip"`           //终端ip
	NetIp        string `json:"netIp"`        //通讯ip
	HostName     string `json:"hostName"`     //计算机名称
	Status       int    `json:"status"`       //终端上线状态
	OnlineTime   string `json:"onlineTime"`   //终端上线时间
	Version      string `json:"version"`      //终端版本
	DeptCode     string `json:"deptCode"`     //组织架构
	DeptName     string `json:"deptName"`     //组织架构
	AgentType    string `json:"agentType"`    //终端类型
	AliasName    string `json:"aliasName"`    //别名
	ContactName  string `json:"contactName"`  //管理人
	ContactPhone string `json:"contactPhone"` //管理人
}

// SyncAgentBySyslog 通过SYSLOG同步终端
func (aa *RequestHuanshi) SyncAgentBySyslog(param AgentData) error {
	//默认值
	param.Platform = aa.RequestPlatform
	param.DataType = DataTypeAgent
	if err := aa.analogSendSyslog("agent", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type AlarmBase struct {
	DeviceMark  string `json:"deviceMark"`  //设备标识
	Platform    string `json:"platform"`    //来源
	Type        int    `json:"type"`        //类型模块
	OrgId       int64  `json:"org_id"`      //企业编号
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
	} else if cache.OrgHsSwitch == enums.SwitchOff {
		return nil
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.analogSendSyslog("alarm", json_helper.Encode(param)); err != nil {
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
	} else if cache.OrgHsSwitch == enums.SwitchOff {
		return nil
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
		param.ContactName = cache.ContactName
		param.DeptName = cache.DeptName
	}

	if err := aa.analogSendSyslog("alarm", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type LogBase struct {
	DeviceMark string `json:"device_mark"` //设备标识
	DeviceIp   string `json:"device_ip"`   //设备IP
	Platform   string `json:"platform"`    //来源
	LogType    int    `json:"log_type"`    //日志类型
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`   //终端标识
	AgentIp    string `json:"agent_ip"`   //终端ip
	CreatedAt  string `json:"created_at"` //终端ip
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

// LogInner 日志数据结构： 对内通讯日志
type LogInner struct {
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
}

// LogOuter 日志数据结构： 对外通讯日志
type LogOuter struct {
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
}

// LogLogin 日志数据结构： 登录日志
type LogLogin struct {
	LogBase
	SourceIp string `json:"source_ip"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Status   int    `json:"status"`
}

// SyncLogProcess 同步日志：进程
func (aa *RequestHuanshi) SyncLogProcess(param LogProcess) error {
	param.DeviceIp = aa.Config.HuanshiCallIp
	param.Platform = aa.RequestPlatform
	param.LogType = LogTypeProcess
	if err, cache := aa.GetAgent(param.OrgId, param.AgentId); err != nil {
		return err
	} else if cache.OrgHsSwitch == enums.SwitchOff {
		return nil
	} else {
		param.DeviceMark = cache.OrgCode
		param.AgentIp = cache.AgentIp
	}
	if err := aa.analogSendSyslog("log", json_helper.Encode(param)); err != nil {
		return err
	}
	return nil
}

type CacheAgent struct {
	OrgCode     string `json:"org_code"`
	OrgHsSwitch int    `json:"org_hs_switch"`
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
func (aa *RequestHuanshi) SaveAgent(orgId int64, agentId string, param CacheAgent) error {
	key := aa.cacheAgentKey(orgId, agentId)
	data := json_helper.Encode(json_helper.Encode(param))
	//新增缓存
	if isOK, err := aa.redisClient.SetEx(key, data, aa.cacheExpire); !isOK {
		return errors.New("cache token fail: " + err.Error())
	}

	return nil
}

// GetAgent 获取终端信息
func (aa *RequestHuanshi) GetAgent(orgId int64, agentId string) (err error, data CacheAgent) {

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
