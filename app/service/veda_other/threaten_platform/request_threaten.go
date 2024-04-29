package threaten_platform

import (
	"errors"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateRequestFactory(c RequestConfig) *requestThreaten {
	c.Domain = strings.TrimRight(c.Domain, "/")
	return &requestThreaten{
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "threaten_token",
		cacheExpire: 600,
		Config:      c,
	}
}

type RequestConfig struct {
	Switch    int
	Domain    string
	AppId     string
	AppSecret string
}

type requestThreaten struct {
	redisClient *redis_factory.RedisClient
	cacheKey    string
	cacheExpire int //缓存有效期
	Config      RequestConfig
}

type responseCommon struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
}

// checkCode 校验code
func (aa *requestThreaten) checkCode(bodyContext string, err error) error {
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
func (aa *requestThreaten) Online() error {
	uri := aa.Config.Domain + "/center/heartbeat"

	bodyContext, err := http_request.GetNoParams(uri, map[string]interface{}{})
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}

	return nil
}

// ResponseToken 生成口令接口返回数据
type ResponseToken struct {
	Data struct {
		AccessToken string `json:"access_token"`
		ExpireIn    int    `json:"expire_in"`
	} `json:"data"`
}

// GenerateToken 生成口令
func (aa *requestThreaten) GenerateToken() error {
	uri := aa.Config.Domain + "/center/auth/token"
	params := map[string]interface{}{
		"appId":     aa.Config.AppId,
		"secretKey": aa.Config.AppSecret,
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
	cacheExpire := aa.cacheExpire
	if data.Data.ExpireIn > 0 {
		cacheExpire = data.Data.ExpireIn - aa.cacheExpire
	}
	if isOK, err := aa.redisClient.SetEx(aa.cacheKey, data.Data.AccessToken, cacheExpire); !isOK {
		return errors.New("cache token fail: " + err.Error())
	}
	return nil
}

func (aa *requestThreaten) GetHeaders() (headers map[string]interface{}) {
	if token := aa.redisClient.Get(aa.cacheKey); !helpers.StringIsEmpty(token) {
		headers["AuthToken"] = token
	}
	return headers
}

// ResponseCheckHashes 验证hash接口返回数据
type ResponseCheckHashes struct {
	Data struct {
		Hash  string `json:"name"`
		Level int    `json:"type"`
		Attr  string `json:"attr"`
	} `json:"data"`
}

// CheckHashes 验证hash集合,若存在返回
/*
参数：[]{
	"req_type":1
	"key" : "",
	"filename" : "",
	"dir" : "",
	"os" : "",
	"parentName": "",
}
*/
func (aa *requestThreaten) CheckHashes(list []map[string]interface{}) (error error, result []*ResponseCheckHashes) {
	uri := aa.Config.Domain + "/center/threaten/checkthreaten"
	params := map[string]interface{}{
		"keys": helpers.JsonEncode(list),
	}
	bodyContext, err := http_request.PostWithForm(uri, aa.GetHeaders(), params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err, result
	}

	err = helpers.Json2Struct(bodyContext, &result)
	if err != nil {
		return errors.New("json_error2: " + err.Error()), result
	}

	return nil, result
}

// PushHashes 推送hash集合
/*
参数：[]{
	"threaten_name": "",
	"type": 0, //约定-1
	"filename": "",
	"os": "",
	"dir": "",
	"attribute": "",
	"uname": "",
	"execdir": "",
	"remark": "",
	"copyright": ""
}
*/
func (aa *requestThreaten) PushHashes(list []map[string]interface{}) error {
	uri := aa.Config.Domain + "/center/threaten/platadd"
	params := map[string]interface{}{
		"keys": helpers.JsonEncode(list),
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}
	return nil
}
