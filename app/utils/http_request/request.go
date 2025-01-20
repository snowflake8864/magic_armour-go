package http_request

import (
	"fmt"
	"github.com/qifengzhang007/goCurl"
	"goskeleton/app/global/variable"
)

func formatResp(resp *goCurl.Response, err error) (string, error) {
	msg := ""
	if err != nil {
		if resp == nil {
			msg = fmt.Sprintf("response fail[resp is nil]：%s\n", err.Error())
		} else {
			msg = fmt.Sprintf("response fail[resp not nil]：%s\n", err.Error())
		}
		variable.ZapLog.Error(msg)
		return "", err
	}

	if txt, err2 := resp.GetContents(); err2 == nil {
		return txt, nil
	} else {
		msg = fmt.Sprintf("content fail：%s\n", err2.Error())
		variable.ZapLog.Error(msg)
		return "", err2
	}
}

// GetNoParams GET方法, 无参数
func GetNoParams(uri string, headers map[string]interface{}) (string, error) {
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	cli := goCurl.CreateHttpClient(goCurl.Options{
		Headers:       headers,
		SetResCharset: "utf-8",
		Timeout:       30,
	})
	return formatResp(cli.Get(uri))
}

// PostWithForm POST方法, Content-Type默认为application/x-www-form-urlencoded
func PostWithForm(uri string, headers, params map[string]interface{}) (string, error) {
	cli := goCurl.CreateHttpClient()

	//默认
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := cli.Post(uri, goCurl.Options{
		Headers:       headers,
		FormParams:    params,
		SetResCharset: "utf-8",
		Timeout:       30,
	})
	return formatResp(resp, err)
}

// PostWithJson POST方法, Content-Type默认为application/json
func PostWithJson(uri string, headers map[string]interface{}, params interface{}) (string, error) {
	cli := goCurl.CreateHttpClient()

	//默认
	headers["Content-Type"] = "application/json"
	resp, err := cli.Post(uri, goCurl.Options{
		Headers:       headers,
		JSON:          params,
		SetResCharset: "utf-8",
		Timeout:       30,
	})
	return formatResp(resp, err)
}
