package upload_files

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/organization"
	"goskeleton/app/utils/response"
)

type OrganizationUpFiles struct {
	UploadBase
}

// 文件上传公共模块表单参数验证器
func (u OrganizationUpFiles) CheckParams(context *gin.Context) {
	errCode, errMsg, data := u.CommCheck(context)
	if errCode != 0 {
		response.Fail(context, errCode, errMsg, data)
		return
	}
	(&organization.Upload{}).StartUpload(context)
}
