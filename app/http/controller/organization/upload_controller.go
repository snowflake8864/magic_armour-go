package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/response"
)

type Upload struct {
}

//	文件上传是一个独立模块，给任何业务返回文件上传后的存储路径即可。
//
// 开始上传
func (u *Upload) StartUpload(context *gin.Context) {
	if r, finnalSavePath := upload_file.Upload(context, upload_file.GetDefaultPath(), []string{}); r == true {
		response.Success(context, my_errors.CurdStatusOkMsg, finnalSavePath)
	} else {
		response.Fail(context, my_errors.FilesUploadFailCode, my_errors.FilesUploadFailMsg, "")
	}
}
