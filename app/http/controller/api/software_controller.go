package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/response"
)

// SoftwareController 软件相关
type SoftwareController struct {
}

// Download 下载
func (u *SoftwareController) Download(ctx *gin.Context) {
	token := ctx.GetString(consts.ValidatorPrefixApi + "token")

	downloadFact := software_service.CreateDownloadFactory()
	loginInfo, err := downloadFact.DecryptConfig(token)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	param := software_models.DownloadModel{
		OrgId:      loginInfo.OrgId,
		SoftwareId: loginInfo.SoftId,
		Remark:     "Offline",
		CreateId:   loginInfo.UserId,
	}
	err, downloadModel := software_service.CreateDownloadFactory().GetDownloadInfo(&param)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.FileDownload(ctx, downloadModel.FullPath)
}
