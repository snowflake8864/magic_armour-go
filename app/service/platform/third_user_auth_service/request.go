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
	if resp.Code != "0" {
		return errors.New("result fail: " + bodyContext), nil
	}
	return nil, &resp
}

func (aa *RequestThird) generateSign(syncSequence int64) map[string]string {
	//签名原始数据
	params := map[string]string{
		"once":       md5_encrypt.MD5(int_helper.IntRandom2Str(10000, 99999)),
		"appKey":     aa.Config.AppKey,
		"signMethod": "SHA-256",
		"ts":         fmt.Sprintf("%d", date_helper.CurTimestampMilli()),

		//url上的参数都需要拼接到 字符串里做签名
		"syncSequence": fmt.Sprintf("%d", syncSequence),
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
		"appKey":     signParams["appKey"],
		"once":       signParams["once"],
		"signMethod": signParams["signMethod"],
		"ts":         signParams["ts"],
		"signData":   signParams["signData"],
	}
	return headers
}

type ResponseOrg struct {
	OrgId   string `json:"orgId"`   // 组织ID
	OrgName string `json:"orgName"` // 机构名称
	//ShortName    string `json:"shortName"`    // 机构简称
	OrgNo   string `json:"orgNo"`   // 机构编号
	OrgType int    `json:"orgType"` // 机构类型（1 单位；2 部门）
	//OrgLevelCode string `json:"orgLevelCode"` // 机构级别编码
	//OrgLevelName string `json:"orgLevelName"` // 机构级别名称
	ShowOrder    int    `json:"showOrder"`    // 同级排序号
	ParentId     string `json:"parentId"`     // 上级组织ID
	ParentOrgNo  string `json:"parentOrgNo"`  // 上级机构编号
	CreateTime   string `json:"createTime"`   // 创建时间
	UpdateTime   string `json:"updateTime"`   // 更新时间
	Status       int    `json:"status"`       // 状态（2000 正常；3000 停用；-1000 注销）
	SyncStatus   int    `json:"syncStatus"`   // 同步状态（2000 正常；800 撤回）
	SyncSequence string `json:"syncSequence"` // 同步序列
	//FullCode     string `json:"fullCode"`     // 机构完整编号
	//FullName     string `json:"fullName"`     // 机构完整名称
}

// PullOrg 拉取组织结构
func (aa *RequestThird) PullOrg() (error, []*ResponseOrg) {
	//拼接地址
	uri := fmt.Sprintf("%s/id-sync/v5/pull/orgs", aa.Config.Addr)
	//参数
	params := map[string]interface{}{
		"syncSequence": aa.Config.OrgSeq,
	}
	//生成签名参数
	signParams := aa.generateSign(aa.Config.OrgSeq)
	//头部数据
	headers := aa.getHeaders(signParams)

	//fmt.Println("uri--------------:")
	//fmt.Println(signParams)
	//fmt.Println(headers)
	//fmt.Println(uri)
	//fmt.Println("uri--------------:")
	//
	bodyContext, err := http_request.GetWithQuery(uri, headers, params)
	err2, respComm := aa.checkCode(bodyContext, err)
	//fmt.Println("respComm--------------:")
	//fmt.Println(json_helper.Encode(respComm))
	//fmt.Println("respComm--------------:")
	if err2 != nil {
		variable.ZapLog.Info(bodyContext)
		return err2, nil
	}

	var result []*ResponseOrg
	switch data := respComm.Result.(type) {
	case []interface{}: //有返回数据
		if err1 := json_helper.ToAny(json_helper.Encode(data), &result); err1 != nil {
			return err1, nil
		}
	}

	return nil, result
}

type ResponseUser struct {
	UserId      string `json:"userId"`      // 统一用户ID
	DisplayName string `json:"displayName"` // 姓名
	WorkNo      string `json:"workNo"`      // 工号
	//IdCardNo    string `json:"idCardNo"`    // 身份证号码
	Mobile    string `json:"mobile"`    // 手机号
	UserEmail string `json:"userEmail"` // 工作网邮箱
	//InternalEmail string `json:"internalEmail"` // 内网邮箱
	IsRoot string `json:"isRoot"` // 是否是管理员

	Status       int    `json:"status"`       // 状态（2000 正常；3000 停用；-1000 注销）
	SyncStatus   int    `json:"syncStatus"`   // 同步状态（2000 正常；800 撤回）
	SyncSequence string `json:"syncSequence"` // 同步序列

	UserJobInfo []*ResponseUserJob `json:"userJobInfo"` // 用户单位信息

	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间

}
type ResponseUserJob struct {
	CompanyId string `json:"companyId"` // 单位ID
	//CompanyNo        string `json:"companyNo"`        // 单位编号
	//JobCompanyStatus int    `json:"jobCompanyStatus"`
	//UserJobType      string `json:"userJobType"`
	DepJobInfo []*ResponseUserJobDept `json:"depJobInfo"` // 用户部门信息

}
type ResponseUserJobDept struct {
	DepId    string `json:"depId"` //部门ID
	Position string `json:"position"`
	//Job           interface{} `json:"job"`
	//JobStatus     int         `json:"jobStatus"`
	//JobType       string      `json:"jobType"`
	//OrgName       string      `json:"orgName"`
	//OrgNo         string      `json:"orgNo"`
	//ParentOrgName string      `json:"parentOrgName"`
	//ParentOrgNo   string      `json:"parentOrgNo"`
	//ShowOrder interface{} `json:"showOrder"`
}

// PullUser 拉取用户
func (aa *RequestThird) PullUser() (error, []*ResponseUser) {
	//拼接地址
	uri := fmt.Sprintf("%s/id-sync/v5/pull/users", aa.Config.Addr)
	//参数
	params := map[string]interface{}{
		"syncSequence": aa.Config.UserSeq,
	}
	//生成签名参数
	signParams := aa.generateSign(aa.Config.UserSeq)
	//头部数据
	headers := aa.getHeaders(signParams)

	//fmt.Println("uri--------------:")
	//fmt.Println(signParams)
	//fmt.Println(headers)
	//fmt.Println(uri)
	//fmt.Println("uri--------------:")
	//
	bodyContext, err := http_request.GetWithQuery(uri, headers, params)
	err2, respComm := aa.checkCode(bodyContext, err)
	//fmt.Println("respComm--------------:")
	//fmt.Println(json_helper.Encode(respComm))
	//fmt.Println("respComm--------------:")
	if err2 != nil {
		variable.ZapLog.Info(bodyContext)
		return err2, nil
	}

	var result []*ResponseUser
	switch data := respComm.Result.(type) {
	case []interface{}: //有返回数据
		if err1 := json_helper.ToAny(json_helper.Encode(data), &result); err1 != nil {
			return err1, nil
		}
	}
	return nil, result
}

// GenerateLoginUrl 生成跳转地址
func (aa *RequestThird) GenerateLoginUrl(callback string) string {
	return fmt.Sprintf("%s/authn-api/v5/cas/%s/login?service=%s", aa.Config.Addr, aa.Config.AppId, callback)
}

type ResponseTicket struct {
	ServiceResponse ResponseTicketResult `json:"serviceResponse"`
}

type ResponseTicketResult struct {
	AuthenticationSuccess *ResponseTicketSuccess `json:"authenticationSuccess,omitempty"`
	AuthenticationFailure *ResponseTicketFailure `json:"authenticationFailure,omitempty"`
}

type ResponseTicketSuccess struct {
	User       string                          `json:"user"`
	Attributes ResponseTicketSuccessAttributes `json:"attributes"`
}

type ResponseTicketSuccessAttributes struct {
	Id      string `json:"id"`
	Account string `json:"account"`
	Email   string `json:"email"`
}

type ResponseTicketFailure struct {
	Code string `json:"code"`
	Desc string `json:"description"`
}

// ValidateTicket 验证票据
func (aa *RequestThird) ValidateTicket(ticket, callback string) (error, *ResponseTicketSuccessAttributes) {
	//拼接地址
	uri := fmt.Sprintf("%s/authn-api/v5/cas/%s/p3/serviceValidate",
		aa.Config.Addr,
		aa.Config.AppId,
	)
	params := map[string]interface{}{
		"ticket":  ticket,
		"format":  "json",
		"service": callback,
	}

	bodyContext, err := http_request.GetWithQuery(uri, map[string]interface{}{}, params)
	if err != nil {
		return errors.New("body_error: " + err.Error()), nil
	}
	resp := ResponseTicket{}
	err = json_helper.ToAny(bodyContext, &resp)
	if err != nil {
		return errors.New("json_error: " + err.Error()), nil
	}
	if resp.ServiceResponse.AuthenticationFailure != nil {
		return fmt.Errorf("code: %s, desc: %s", resp.ServiceResponse.AuthenticationFailure.Code, resp.ServiceResponse.AuthenticationFailure.Desc), nil
	}

	return nil, &resp.ServiceResponse.AuthenticationSuccess.Attributes
}
