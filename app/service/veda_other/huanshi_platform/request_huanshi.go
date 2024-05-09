package huanshi_platform

import (
	"errors"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateRequestFactory(c RequestConfig) *requestHuanshi {
	c.HuanshiUrl = strings.TrimRight(c.HuanshiUrl, "/")
	return &requestHuanshi{
		redisClient:     redis_factory.GetOneRedisClient(),
		cacheKey:        "huanshi_token",
		cacheExpire:     600,
		Config:          c,
		PrivateRsa:      variable.HuanshiPrivateRsa,
		RequestPlatform: "huanjia",
		Edition:         consts.VersionNumber,
	}
}

type RequestConfig struct {
	HuanshiCallbackUrl          string
	HuanshiCallIp               string
	HuanshiSwitch               int
	HuanshiUrl                  string
	HuanshiAlarmSwitch          int
	HuanshiAlarmIp              string
	HuanshiAlarmPort            string
	HuanshiLogSwitch            int
	HuanshiLogIp                string
	HuanshiLogPort              string
	HuanshiAlarmSwitchProcess   int
	HuanshiAlarmSwitchDll       int
	HuanshiAlarmSwitchTamper    int
	HuanshiAlarmSwitchBlackmail int
	HuanshiAlarmSwitchPort      int
	HuanshiAlarmSwitchUsb       int
	HuanshiAlarmSwitchPlug      int
	HuanshiAlarmSwitchLogin     int
	HuanshiAlarmSwitchOutreach  int
	HuanshiAlarmSwitchSelf      int
}

type requestHuanshi struct {
	redisClient     *redis_factory.RedisClient
	cacheKey        string
	cacheExpire     int //缓存有效期
	Config          RequestConfig
	PrivateRsa      string
	RequestPlatform string
	Edition         string
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
	if code == 3007 {
		return errors.New("token expired: " + bodyContext)
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
"deviceMark" : "设备标识",
"platform" => "请求平台"
"edition" =>  "产品版本号",
"ip" => "回调IP",
"url" => "回调地址",
"agentNum"=> "终端总数",
"onlineAgentNum"=> "在线终端数",
}
*/
func (aa *requestHuanshi) SyncDeviceByTcp(param map[string]interface{}) error {
	uri := aa.Config.HuanshiUrl + "/web/device/add"

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
