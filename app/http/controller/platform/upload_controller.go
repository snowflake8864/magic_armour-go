package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/platform/upload_service"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/response"
)

type Upload struct {
	ControllerBase
}

//	文件上传是一个独立模块，给任何业务返回文件上传后的存储路径即可。
//
// 开始上传
func (u *Upload) StartUpload(context *gin.Context) {
	fileType, _ := context.GetPostForm("type")
	allowFileTypes := map[string]struct{}{
		"software": {},
		"upgrade":  {},
		"license":  {},
	}
	if _, isOK := allowFileTypes[fileType]; !isOK {
		response.FailSimplify(context, "type fail:"+fileType)
		return
	}

	savePath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath")
	if r, uploadParam := upload_file.Upload(context, savePath); r == true {
		uploadParam.Type = fileType
		loginInfo := u.getLogin(context)
		uploadParam.OperateSource = enums.OperateSourcePlatform
		uploadParam.OperateId = loginInfo.AdminId
		uploadParam.OperateName = loginInfo.Username
		if !upload_service.CreateCurdFactory().Add(&uploadParam) {
			response.FailSimplify(context, "保存失败")
			return
		}

		response.SuccessReturnData(context, gin.H{
			"id":   uploadParam.Id,
			"path": uploadParam.StoragePath,
		})
	} else {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}
