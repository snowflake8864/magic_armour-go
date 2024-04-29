package threaten_service

import (
	"errors"
	"fmt"
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

type responseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// dealResponse 处理请求结果
func (aa *requestThreaten) dealResponse(bodyContext string, err error) (error, interface{}) {
	if err != nil {
		return err, nil
	}
	resp := responseData{}
	err = helpers.Json2Struct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_last_error: " + err.Error()), nil
	}

	//口令过期
	if resp.Code == 3007 {
		return errors.New(fmt.Sprintf("token expired[code=%d]", resp.Code)), nil
	}
	//结果不正确
	if resp.Code == 1 {
		return errors.New(fmt.Sprintf("result fail[code=%d]", resp.Code)), nil
	}
	return nil, resp.Data
}

// Online 是否在线
func (aa *requestThreaten) Online() error {
	uri := aa.Config.Domain + "/center/heartbeat"
	err, _ := aa.dealResponse(http_request.GetNoParams(uri, nil))
	return err
}

// GenerateToken 生成口令
func (aa *requestThreaten) GenerateToken() error {
	uri := aa.Config.Domain + "/center/auth/token"
	params := map[string]interface{}{
		"appId":     aa.Config.AppId,
		"secretKey": aa.Config.AppSecret,
	}
	err, resp := aa.dealResponse(http_request.PostWithForm(uri, nil, params))
	fmt.Println(resp)
	fmt.Sprintf("resp %T \n", resp)

	return err
}
