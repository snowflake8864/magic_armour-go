package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/utils/validator_translation"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ReturnJson(Context *gin.Context, httpCode int, dataCode interface{}, msg string, data interface{}) {

	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

// ReturnJsonFromString 将json字符窜以标准json格式返回（例如，从redis读取json格式的字符串，返回给浏览器json格式）
func ReturnJsonFromString(Context *gin.Context, httpCode int, jsonStr string) {
	Context.Header("Content-Type", "application/json; charset=utf-8")
	Context.String(httpCode, jsonStr)
}

// 语法糖函数封装

// Success 直接返回成功
func Success(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, msg, data)
}

// SuccessSimplify 直接返回成功: 简化
func SuccessSimplify(c *gin.Context) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, "")
}

// SuccessReturnData 直接返回成功: 有数据
func SuccessReturnData(c *gin.Context, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, consts.CurdStatusOkMsg, data)
}

// Fail 失败的业务逻辑
func Fail(c *gin.Context, dataCode int, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, dataCode, msg, data)
	c.Abort()
}

// FailSimplify 失败的业务逻辑：简化
func FailSimplify(c *gin.Context, msg string) {
	ReturnJson(c, http.StatusOK, consts.CurdBaseCode, msg, "")
	c.Abort()
}

// FailRoutine 失败的业务逻辑： 日常使用
func FailRoutine(c *gin.Context, dataCode interface{}, msg string) {
	ReturnJson(c, http.StatusOK, dataCode, msg, "")
	c.Abort()
}

// FailData 失败的业务逻辑：有错误数据
func FailData(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdBaseCode, msg, data)
	c.Abort()
}

// ErrorTokenBaseInfo token 基本的格式错误
func ErrorTokenBaseInfo(c *gin.Context) {
	ReturnJson(c, http.StatusBadRequest, http.StatusBadRequest, my_errors.ErrorsTokenBaseInfo, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// ErrorTokenFail 口令错误
func ErrorTokenFail(c *gin.Context, dataCode, wrongParam interface{}) {
	ReturnJson(c, http.StatusUnauthorized, dataCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

// ErrorTokenAuthFail token 权限校验失败
func ErrorTokenAuthFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsNoAuthorization, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// ErrorTokenRefreshFail token不符合刷新条件
func ErrorTokenRefreshFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsRefreshTokenFail, "")
	//终止可能已经被加载的其他回调函数的执行
	c.Abort()
}

// TokenErrorParam token 参数校验错误
func TokenErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusUnauthorized, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

// ErrorCasbinAuthFail 鉴权失败，返回 405 方法不允许访问
func ErrorCasbinAuthFail(c *gin.Context, msg interface{}) {
	ReturnJson(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, my_errors.ErrorsCasbinNoAuthorization, msg)
	c.Abort()
}

// ErrorParam 参数校验错误
func ErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

// ErrorSystem 系统执行代码错误
func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.ServerOccurredErrorCode, consts.ServerOccurredErrorMsg+msg, data)
	c.Abort()
}

// ValidatorError 翻译表单参数验证器出现的校验错误
func ValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		wrongParam := validator_translation.RemoveTopStruct(errs.Translate(validator_translation.Trans))
		ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	} else {
		errStr := err.Error()
		// multipart:nextpart:eof 错误表示验证器需要一些参数，但是调用者没有提交任何参数
		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": my_errors.ErrorNotAllParamsIsBlank})
		} else {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": errStr})
		}
	}
	c.Abort()
}

// FileDownload 文件下载
func FileDownload(c *gin.Context, fullDir string) {
	//绝对文件路径
	fullPath, err := filepath.Abs(fullDir)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 通过 http.ServeFile 将文件内容传输给客户端
	http.ServeFile(c.Writer, c.Request, fullPath)

}

// FileDownloadAndRename 文件下载且修改文件名称
// @newName 中文会转码，不然会被截断
func FileDownloadAndRename(c *gin.Context, fullDir, newName string) {
	//绝对文件路径
	fullPath, err := filepath.Abs(fullDir)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 准备文件记录
	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	newName = url.PathEscape(newName)
	// 设置Header信息
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+newName)
	//fmt.Println(c.Writer.Header().Get("Content-Disposition"))
	// 服务文件
	http.ServeContent(c.Writer, c.Request, newName, time.Now(), file)
}
