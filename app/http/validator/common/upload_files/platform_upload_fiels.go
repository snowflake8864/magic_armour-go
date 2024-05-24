package upload_files

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/utils/response"
)

type PlatformUpFiles struct {
	UploadBase
}

// 文件上传公共模块表单参数验证器
func (u PlatformUpFiles) CheckParams(context *gin.Context) {
	errCode, errMsg, data := u.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}
	(&platform.Upload{}).StartUpload(context)
}
