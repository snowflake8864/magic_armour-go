package third_user_auth_service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/third_user_auth_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/http_request"
	"goskeleton/app/utils/md5_encrypt"
	"sort"
	"strings"
)

func CreateRequestFactory(param *third_user_auth_models.AuthModel) *RequestThird {
	param.Addr = strings.TrimRight(param.Addr, "/")
	return &RequestThird{
		Config: param,
	}
}

type RequestThird struct {
	Config *third_user_auth_models.AuthModel
}

type ResponseCommon struct {
	Code      string      `json:"code"`
	Timestamp string      `json:"timestamp"`
	Error     string      `json:"error"`
	Result    interface{} `json:"result"`
}

// checkCode 校验code
func (aa *RequestThird) checkCode(bodyContext string, err error) (error, *ResponseCommon) {
	if err != nil {
		return errors.New("body_error: " + err.Error()), nil
	}
	resp := ResponseCommon{}
	err = json_helper.ToStruct(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error()), nil
	}

	//fmt.Println("Type of x:", reflect.TypeOf(resp.Code))
	//code := float_helper.Float64(resp.Code)
	//结果不正确
	if resp.Code != "200" {
		return errors.New("result fail: " + bodyContext), nil
	}
	return nil, &resp
}

func (aa *RequestThird) generateSign() map[string]string {
	//签名原始数据
	params := map[string]string{
		"appKey":     aa.Config.AppKey,
		"once":       md5_encrypt.MD5(int_helper.IntRandom2Str(10000, 99999)),
		"signMethod": "SHA-256",
		"startTime":  date_helper.SpecialDatetime("2006-01-02T15:04:05.000Z"),
		"ts":         fmt.Sprintf("%d", date_helper.CurTimestampMilli()),
	}

	params["signData"] = aa.signData(params, aa.Config.AppSecret)

	return params
}

func (aa *RequestThird) signData(params map[string]string, appSecret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}

	source := appSecret + ":" + sb.String() + ":" + appSecret
	sum := sha256.Sum256([]byte(source))
	return hex.EncodeToString(sum[:])
}

func (aa *RequestThird) getHeaders(signParams map[string]string) map[string]interface{} {

	headers := map[string]interface{}{
		"appKey":       signParams["appKey"],
		"once":         signParams["once"],
		"signMethod":   signParams["signMethod"],
		"ts":           signParams["ts"],
		"signData":     signParams["signData"],
		"Content-Type": "application/json;charset=utf-8",
	}

	fmt.Println("headers--------------:")
	fmt.Println(json_helper.Encode(headers))
	fmt.Println("headers--------------:")
	return headers
}

type ResponseOrg struct {
}

// PullOrg 验证hash集合,若存在返回
func (aa *RequestThird) PullOrg() (error, []*ResponseOrg) {
	//生成签名参数
	signParams := aa.generateSign()

	//拼接地址
	//uri := fmt.Sprintf("%s/ums/service/V1/res/orgs?startTime=%s&syncSequence=%d", aa.Config.Addr, signParams["startTime"], aa.Config.OrgSeq)
	uri := fmt.Sprintf("%s/id-sync/v5/pull/orgs?startTime=%s&syncSequence=%d", aa.Config.Addr, signParams["startTime"], aa.Config.OrgSeq)
	fmt.Println("uri--------------:")
	fmt.Println(uri)
	fmt.Println("uri--------------:")
	bodyContext, err := http_request.GetNoParams(uri, aa.getHeaders(signParams))

	err2, respComm := aa.checkCode(bodyContext, err)
	if err2 != nil {
		variable.ZapLog.Info(bodyContext)
		return err2, nil
	}
	fmt.Println("respComm.Result--------------:")
	fmt.Println(respComm.Result)
	fmt.Println("respComm.Result--------------:")
	var result []*ResponseOrg
	switch data := respComm.Result.(type) {
	case []interface{}: //没有返回数据
		return nil, nil
	case map[string]interface{}: //有返回数据
		for _, mapStr := range data {
			tmp := ResponseOrg{}
			jsonStr := json_helper.Encode(mapStr)
			if err1 := json_helper.ToStruct(jsonStr, &tmp); err1 == nil {
				result = append(result, &tmp)
			}
		}
		break
	}
	return nil, result
}
