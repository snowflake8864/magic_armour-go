package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/string_helper"
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
		Remark:     "离线下载",
		CreateId:   loginInfo.UserId,
	}
	err, downloadModel := software_service.CreateDownloadFactory().GetDownloadInfo(&param)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.FileDownload(ctx, downloadModel.FullPath)
}

// PublishedList 已发布列表
func (u *SoftwareController) PublishedList(ctx *gin.Context) {
	IsPublish := enums.IsYes
	condition := software_models.SearchConditions{
		IsPublish: &IsPublish,
	}

	models := software_models.CreateSoftwareFactory().GetItemsByConditions(condition, 0, 0, "")
	response.SuccessReturnData(ctx, (&org_trans.Software{}).TransList(models))
}

func (u *SoftwareController) getUser() (int64, int64) {
	var orgId, userId int64
	dictIds := []string{
		enum_dictionary.AgentCompatibleDefaultOrg,
		enum_dictionary.AgentCompatibleDefaultUser,
	}

	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if val, isOk := mapDict[enum_dictionary.AgentCompatibleDefaultOrg]; !isOk || string_helper.IsEmpty(val) {
		return 0, 0
	} else {
		orgId = string_helper.ToInt64(val)
	}
	if val, isOk := mapDict[enum_dictionary.AgentCompatibleDefaultUser]; !isOk || string_helper.IsEmpty(val) {
		return 0, 0
	} else {
		userId = string_helper.ToInt64(val)
	}
	return orgId, userId
}

func (u *SoftwareController) OfflineDownload(ctx *gin.Context) {

	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixApi + "softId"))
	model := software_service.CreateCurdFactory().GetItemById(softId)
	if model == nil || model.Id == 0 {
		response.FailSimplify(ctx, "软件不存在")
		return
	}
	downloadFact := software_service.CreateDownloadFactory()
	err, addr, _, _ := downloadFact.GetServiceAddr()
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	orgId, userId := u.getUser()
	verId := downloadFact.EncryptConfig(orgId, userId, softId)

	response.SuccessReturnData(ctx, gin.H{
		"url": enum_software.GetOfflineDownload(model.PlatformType, addr, verId),
	})
}

func (u *SoftwareController) OnlineDownload(ctx *gin.Context) {
	orgId, userId := u.getUser()

	param := software_models.DownloadModel{
		OrgId:      orgId,
		SoftwareId: int64(ctx.GetFloat64(consts.ValidatorPrefixApi + "softId")),
		Remark:     "立即下载",
		CreateId:   userId,
		CreateName: "",
	}

	if err, downloadModel := software_service.CreateDownloadFactory().GetDownloadInfo(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		response.SuccessReturnData(ctx, gin.H{
			"path": downloadModel.StoragePath,
		})
		//response.FileDownload(ctx, downloadModel.FullPath)
	}
}
