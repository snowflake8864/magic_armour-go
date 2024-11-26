package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/organization/order_service"
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

		data := gin.H{
			"id":   uploadParam.Id,
			"path": uploadParam.StoragePath,
		}
		//返回授权
		if fileType == "license" {
			if err, resp := order_service.CreateLicenseFactory().Parse(uploadParam.FullPath); err != nil {
				response.FailSimplify(context, err.Error())
				return
			} else {
				data["data"] = resp
			}
		}

		response.SuccessReturnData(context, data)
	} else {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}
