package threaten_platform

import (
	"errors"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

const ReqType int = 1 //约定

func CreateRequestFactory(c RequestConfig) *RequestThreaten {
	c.Domain = strings.TrimRight(c.Domain, "/")
	return &RequestThreaten{
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "threaten_token",
		cacheExpire: 600,
		Config:      c,
	}
}

type RequestConfig struct {
	Domain    string
	AppId     string
	AppSecret string
}

type RequestThreaten struct {
	redisClient *redis_factory.RedisClient
	cacheKey    string
	cacheExpire int //缓存有效期
	Config      RequestConfig
}

type responseCommon struct {
	Code interface{} `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// checkCode 校验code
func (aa *RequestThreaten) checkCode(bodyContext string, err error) (error, *responseCommon) {
	if err != nil {
		return errors.New("body_error: " + err.Error()), nil
	}
	resp := responseCommon{}
	err = json_helper.ToStruct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error()), nil
	}

	//fmt.Println("Type of x:", reflect.TypeOf(resp.Code))
	code := float_helper.Float64(resp.Code)
	//口令过期
	if code == 3007 {
		return errors.New("token expired: " + bodyContext), nil
	}
	//结果不正确
	if code != 1 {
		return errors.New("result fail: " + bodyContext), nil
	}
	return nil, &resp
}

// Online 是否在线
func (aa *RequestThreaten) Online() error {
	uri := aa.Config.Domain + "/center/heartbeat"

	bodyContext, err := http_request.GetNoParams(uri, map[string]interface{}{})
	if err1, _ := aa.checkCode(bodyContext, err); err1 != nil {
		return err1
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
func (aa *RequestThreaten) GenerateToken() error {
	uri := aa.Config.Domain + "/center/auth/token"
	params := map[string]interface{}{
		"appId":     aa.Config.AppId,
		"secretKey": aa.Config.AppSecret,
	}
	bodyContext, err := http_request.PostWithForm(uri, map[string]interface{}{}, params)
	if err1, _ := aa.checkCode(bodyContext, err); err1 != nil {
		return err1
	}

	data := ResponseToken{}
	if err1 := json_helper.ToStruct(bodyContext, &data); err1 != nil {
		return errors.New("json_error2: " + err1.Error())
	}

	//新增缓存
	cacheExpire := aa.cacheExpire
	if data.Data.ExpireIn > 0 {
		cacheExpire = data.Data.ExpireIn - aa.cacheExpire
	}
	if isOK, err1 := aa.redisClient.SetEx(aa.cacheKey, data.Data.AccessToken, cacheExpire); !isOK {
		return errors.New("cache token fail: " + err1.Error())
	}
	return nil
}

func (aa *RequestThreaten) GetHeaders() (headers map[string]interface{}) {
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
	Hash  string `json:"name"` //hash, 大写
	Level int    `json:"type"` //等级
	Attr  string `json:"attr"` //属性
}

type CheckHashes struct {
	ReqType    int    `json:"req_type"`   //固定：1
	Key        string `json:"key"`        //hash
	Filename   string `json:"filename"`   //文件名
	Dir        string `json:"dir"`        //文件路径
	Os         string `json:"os"`         //终端系统
	ParentName string `json:"parentName"` //父文件路径
}

// CheckHashes 验证hash集合,若存在返回
func (aa *RequestThreaten) CheckHashes(list *[]CheckHashes) (error, *[]ResponseCheckHashes) {
	uri := aa.Config.Domain + "/center/threaten/checkthreaten"
	params := map[string]interface{}{
		"keys": json_helper.Encode(&list),
	}
	bodyContext, err := http_request.PostWithJson(uri, aa.GetHeaders(), params)

	err2, respComm := aa.checkCode(bodyContext, err)
	if err2 != nil {
		variable.ZapLog.Info(bodyContext)
		return err2, nil
	}

	var result []ResponseCheckHashes
	switch data := respComm.Data.(type) {
	case []interface{}: //没有返回数据
		return nil, nil
	case map[string]interface{}: //有返回数据
		for _, mapStr := range data {
			tmp := ResponseCheckHashes{}
			jsonStr := json_helper.Encode(mapStr)
			if err1 := json_helper.ToStruct(jsonStr, &tmp); err1 == nil {
				result = append(result, tmp)
			}
		}
		break
	}
	return nil, &result
}

func (aa *RequestThreaten) TestCheckData(bodyContext string) *[]ResponseCheckHashes {
	fmt.Println("====================================")
	fmt.Println(bodyContext)
	err, respComm := aa.checkCode(bodyContext, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(json_helper.Encode(respComm))
	fmt.Println(fmt.Sprintf("data type:%T ", respComm.Data))

	var result []ResponseCheckHashes
	switch data := respComm.Data.(type) {
	case []interface{}:
		return nil
	case map[string]interface{}:
		for hash, mapStr := range data {
			fmt.Println("hash======11111111111111: ", hash)
			fmt.Println(mapStr)
			fmt.Println(fmt.Sprintf("mapStr type:%T ", mapStr))

			tmp := ResponseCheckHashes{}
			jsonStr := json_helper.Encode(mapStr)
			if err1 := json_helper.ToStruct(jsonStr, &tmp); err1 == nil {
				result = append(result, tmp)
			}

		}
		break
	}
	fmt.Println("result===================")
	fmt.Println(json_helper.Encode(result))
	return nil
}

type PushHashes struct {
	Hash      string `json:"threaten_name"` //hash
	Type      int    `json:"type"`          ////约定-1
	Filename  string `json:"filename"`      //文件名
	Dir       string `json:"dir"`           //文件路径
	Os        string `json:"os"`            //终端系统
	Uname     string `json:"uname"`         //推送用户名
	ExecDir   string `json:"execdir"`       //父文件路径
	Attribute string `json:"attribute"`     //属性
	Remark    string `json:"remark"`        //描述
	Copyright string `json:"copyright"`     //版权
}

// PushHashes 推送hash集合
func (aa *RequestThreaten) PushHashes(list *[]PushHashes) error {
	uri := aa.Config.Domain + "/center/threaten/platadd"

	params := map[string]interface{}{
		"keys": json_helper.Encode(&list),
	}
	bodyContext, err := http_request.PostWithJson(uri, aa.GetHeaders(), params)
	if err1, _ := aa.checkCode(bodyContext, err); err1 != nil {
		return err1
	}
	return nil
}
