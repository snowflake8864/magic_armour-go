package huanshi_platform

import (
	"errors"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"net"
	"strings"
)

const (
	AlarmModuleProcess   int = 1  //告警模块: 进程
	AlarmModuleDll       int = 2  //告警模块: 模块
	AlarmModuleTamper    int = 3  //告警模块: 防篡改
	AlarmModuleBlackmail int = 4  //告警模块: 虚开
	AlarmModulePort      int = 5  //告警模块: 端口
	AlarmModuleUsb       int = 6  //告警模块: 外设
	AlarmModulePlug      int = 7  //告警模块: 动态封堵
	AlarmModuleOutreach  int = 8  //告警模块: 外联
	AlarmModuleLogin     int = 9  //告警模块: 口令破解
	AlarmModuleSelf      int = 10 //告警模块: 自保

	LogTypeDns     int = 1 //日志类型：dns
	LogTypeOuter   int = 2 //日志类型：对外通讯
	LogTypeInner   int = 3 //日志类型：对内通讯
	LogTypeProcess int = 4 //日志类型：进程启动
	LogTypeLogin   int = 5 //日志类型：登录
)

func CreateRequestFactory(c RequestConfig) *requestHuanshi {
	c.HuanshiUrl = strings.TrimRight(c.HuanshiUrl, "/")
	return &requestHuanshi{
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
	HuanshiSwitch      int
	HuanshiUrl         string
	HuanshiAlarmSwitch int
	HuanshiAlarmIp     string
	HuanshiAlarmPort   string
	HuanshiLogSwitch   int
	HuanshiLogIp       string
	HuanshiLogPort     string
}

type requestHuanshi struct {
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
func (aa *requestHuanshi) checkCode(bodyContext string, err error) error {
	if err != nil {
		return errors.New("body_error: " + err.Error())
	}
	resp := responseCommon{}
	err = helpers.Json2Struct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error())
	}

	//fmt.Println("Type of x:", reflect.TypeOf(resp.Code))
	code := helpers.Float64(resp.Code)
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
func (aa *requestHuanshi) Online() error {
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
func (aa *requestHuanshi) GenerateToken() error {
	uri := aa.Config.HuanshiUrl + "/web/site_auth/login"
	Timestamp := helpers.TimestampCurr()
	//  $param['rsa'] = md5( $param['url'].'_'.$param['time'].'_'.$hs_private_rsa);
	params := map[string]interface{}{
		"url":  aa.Config.HuanshiCallbackUrl,
		"time": Timestamp,
		"rsa":  md5_encrypt.MD5(aa.Config.HuanshiCallbackUrl + "_" + helpers.Int64ToStr(Timestamp) + "_" + aa.PrivateRsa),
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}

	data := ResponseToken{}
	err = helpers.Json2Struct(bodyContext, &data)
	if err != nil {
		return errors.New("json_error2: " + err.Error())
	}

	//新增缓存
	if isOK, err := aa.redisClient.SetEx(aa.cacheKey, data.Data.AccessToken, aa.cacheExpire); !isOK {
		return errors.New("cache token fail: " + err.Error())
	}
	return nil
}

func (aa *requestHuanshi) GetHeaders() (headers map[string]interface{}) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	if token := aa.redisClient.Get(aa.cacheKey); !helpers.StringIsEmpty(token) {
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
func (aa *requestHuanshi) SyncDeviceByTcp(param map[string]interface{}) error {
	if aa.Config.HuanshiSwitch == enums.SwitchOn {
		return nil
	}

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
func (aa *requestHuanshi) analogSendSyslog(tt string, facility, level int, param map[string]interface{}) error {
	if aa.Config.HuanshiSwitch == enums.SwitchOn {
		return nil
	}

	//发送地址
	address := ""
	if tt == "agent" {
		address = aa.Config.HuanshiAlarmIp + ":" + aa.Config.HuanshiAlarmPort
	} else if tt == "alarm" {
		if aa.Config.HuanshiAlarmSwitch == enums.SwitchOn {
			return nil
		}
		address = aa.Config.HuanshiAlarmIp + ":" + aa.Config.HuanshiAlarmPort
	} else if tt == "log" {
		if aa.Config.HuanshiLogSwitch == enums.SwitchOn {
			return nil
		}
		address = aa.Config.HuanshiLogIp + ":" + aa.Config.HuanshiLogPort
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
	msg += helpers.IntToStr((facility * 8) + level)
	msg += ">"
	msg += helpers.JsonEncode(param)

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return errors.New("net.Write err:" + err.Error())
	}
	return nil
}

// SyncAgentBySyslog 通过SYSLOG同步终端
/*
参数：param {
"agentId" : "" //终端编号
"os" : ""  //系统
"ip" : ""  //终端ip
"netIp" : ""  //通讯ip
"hostName" : "" //计算机名称
"status": 1  //终端上线状态: 0 1
"onlineTime": "终端上线时间"
"version": "终端版本"
"deptCode" : "" //组织架构
"deptName" : //组织架构
"agentType" : "" //终端类型
"aliasName" : "" //别名
"contactName" : "" //管理人
"contactPhone": "" //管理人
"deviceMark": "" //设备标识
"platform":"" //请求平台
"dataType": "zhongduan"
}
*/
func (aa *requestHuanshi) SyncAgentBySyslog(param map[string]interface{}) error {
	//默认值
	param["platform"] = aa.RequestPlatform
	param["dataType"] = "zhongduan"
	//@todo 通道:LOCAL7; 等级:info
	if err := aa.analogSendSyslog("agent", 23, 6, param); err != nil {
		return err
	}
	return nil
}

// SyncAlarmBySyslog 同步告警
// 参数：param 每种告警不相同
// 功能模块: 参考宏定义
func (aa *requestHuanshi) SyncAlarmBySyslog(param map[string]interface{}) error {
	//默认值
	param["platform"] = aa.RequestPlatform
	param["dataType"] = "zhongduan"
	//@todo 通道:LOCAL5; 等级:info
	if err := aa.analogSendSyslog("alarm", 21, 6, param); err != nil {
		return err
	}
	return nil
}

// SyncLogBySyslog 同步日志
// 参数：param 每种日志不相同
// 日志类型: 参考宏定义
func (aa *requestHuanshi) SyncLogBySyslog(param map[string]interface{}) error {
	//默认值
	param["platform"] = aa.RequestPlatform
	param["device_ip"] = aa.Config.HuanshiCallIp
	//@todo 通道:LOCAL3; 等级:warning
	if err := aa.analogSendSyslog("log", 19, 4, param); err != nil {
		return err
	}
	return nil
}

type CacheAgent struct {
	AgentIp     string `json:"agent_ip"`
	NetIp       string `json:"net_ip"`
	DeptName    string `json:"dept_name"`
	ContactName string `json:"contact_name"`
	AgentType   string `json:"agent_type"`
}

func (aa *requestHuanshi) cacheAgentKey(deviceMark, agentId string) string {

	return aa.agentCachePrefix + deviceMark + "_" + agentId
}

// SaveAgent 保存终端信息
func (aa *requestHuanshi) SaveAgent(deviceMark, agentId string, param CacheAgent) error {
	key := aa.cacheAgentKey(deviceMark, agentId)
	data := helpers.JsonEncode(param)
	//新增缓存
	if isOK, err := aa.redisClient.SetEx(key, data, aa.cacheExpire); !isOK {
		return errors.New("cache token fail: " + err.Error())
	}

	return nil
}

// GetAgent 获取终端信息
func (aa *requestHuanshi) GetAgent(deviceMark, agentId string) (data CacheAgent) {

	jsonStr := aa.redisClient.Get(aa.cacheAgentKey(deviceMark, agentId))
	if helpers.StringIsEmpty(jsonStr) {
		return data
	}

	err := helpers.Json2Struct(jsonStr, &data)
	if err != nil {
		return data
	}
	return data
}
