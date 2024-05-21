package threaten_platform

import (
	"errors"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
	err = json_helper.ToStruct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error())
	}

	//fmt.Println("Type of x:", reflect.TypeOf(resp.Code))
	code := float_helper.Float64(resp.Code)
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
	err = json_helper.ToStruct(bodyContext, &data)
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
	if headers == nil {
		headers = make(map[string]interface{})
	}

	if token := aa.redisClient.Get(aa.cacheKey); !string_helper.IsEmpty(token) {
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

type CheckHashes struct {
	ReqType    int    `json:"req_Type"`   //固定：1
	Key        string `json:"key"`        //hash
	Filename   string `json:"filename"`   //文件名
	Dir        string `json:"dir"`        //文件路径
	Os         string `json:"os"`         //终端系统
	ParentName int    `json:"parentName"` //父文件路径
}

// CheckHashes 验证hash集合,若存在返回
func (aa *requestThreaten) CheckHashes(list *[]CheckHashes) (error error, result []*ResponseCheckHashes) {
	uri := aa.Config.Domain + "/center/threaten/checkthreaten"
	params := map[string]interface{}{
		"keys": json_helper.Encode(&list),
	}
	bodyContext, err := http_request.PostWithForm(uri, aa.GetHeaders(), params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err, result
	}

	err = json_helper.ToStruct(bodyContext, &result)
	if err != nil {
		return errors.New("json_error2: " + err.Error()), result
	}

	return nil, result
}

type PushHashes struct {
	ThreatenName string `json:"threaten_name"` //hash
	Type         int    `json:"req_Type"`      ////约定 -1
	Filename     string `json:"filename"`      //文件名
	Dir          string `json:"dir"`           //文件路径
	Os           string `json:"os"`            //终端系统
	Uname        string `json:"uname"`         //推送用户名
	ExecDir      string `json:"execdir"`       //父文件路径
	Attribute    string `json:"attribute"`     //属性
	Remark       string `json:"remark"`        //描述
	Copyright    string `json:"copyright"`     //版权
}

// PushHashes 推送hash集合
func (aa *requestThreaten) PushHashes(list *[]PushHashes) error {
	uri := aa.Config.Domain + "/center/threaten/platadd"
	params := map[string]interface{}{
		"keys": json_helper.Encode(&list),
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err := aa.checkCode(bodyContext, err); err != nil {
		return err
	}
	return nil
}
