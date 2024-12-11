package middleware_auth_platform

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/json_helper"
	"io"
	"strings"
)

// RequestParams 将请求体的数据流保存下来并重新设置到Request.Body中，以便后续的绑定操作可以重新读取请求体的数据
// (参考https://blog.csdn.net/qq_39272466/article/details/131735319)
func RequestParams(ctx *gin.Context) {
	bindParams := make(map[string]interface{})
	if ctx.Request.Method == "POST" {
		contextType := strings.ToLower(ctx.Request.Header.Get("Content-Type"))
		requestBody, err := ctx.GetRawData()
		if err != nil {
			variable.ZapLog.Error("PlatformRequestParams GetRawData,err: " + err.Error())
			return
		}
		// 重置请求体，以便后续的中间件和处理程序可以正常读取
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		defer func() {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}()
		// 绑定参数
		if strings.Contains(contextType, "application/json") {
			err1 := ctx.ShouldBind(&bindParams)
			if err1 != nil { // 报错
				variable.ZapLog.Error("PlatformRequestParams 获取参数异常: " + err1.Error())
				return
			}
		} else {
			_ = ctx.Request.ParseMultipartForm(32 << 20)
			if len(ctx.Request.PostForm) > 0 {
				for k, v := range ctx.Request.PostForm {
					bindParams[k] = v[0]
				}
			}
		}
	}
	ctx.Set("PlatformRequestParams", bindParams)
}

// JoinParamsStr 获取请求参数,返回 key=val格式
func JoinParamsStr(ctx *gin.Context) string {
	params := ctx.GetStringMap("PlatformRequestParams")
	var newParams []string
	if len(params) > 0 {
		for k, v := range params {
			newParams = append(newParams, fmt.Sprintf("%s=%v", k, v))
		}
	}
	return strings.Join(newParams, "&")
}

// GetRequestParams 获取请求参数
func GetRequestParams(ctx *gin.Context) string {
	return json_helper.Encode(ctx.GetStringMap("PlatformRequestParams"))
}
