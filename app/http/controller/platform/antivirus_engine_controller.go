package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/antivirus_engine_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/response"
)

// AntivirusEngineController 杀毒引擎相关
type AntivirusEngineController struct {
	ControllerBase
}

// PageForManage 管理: 分页
func (u *AntivirusEngineController) PageForManage(ctx *gin.Context) {

	//页面上传
	PlatformType := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))
	IsPublish := enums.IsNot
	condition := software_models.SearchConditions{
		PlatformType: PlatformType,
		IsPublish:    &IsPublish,
	}
	if PlatformType == enum_software.PlatformTypeWindows || PlatformType == enum_software.PlatformTypeLinux {
		condition.PlatformType = PlatformType
		IsPublish = enums.IsYes
		condition.IsPublish = &IsPublish
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//软件分页
	total, models := software_models.CreateSoftwareFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Software{}).TransPage(total, models))
}

// AddForManage 管理:  添加
func (u *AntivirusEngineController) AddForManage(ctx *gin.Context) {
	var param = antivirus_engine_models.EngineModel{}
	param.IsNewest = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isNewest"))
	param.PlatformType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))
	param.UploadId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "uploadId"))
	param.Ver = ctx.GetString(consts.ValidatorPrefixPlatform + "ver")
	param.EngineType = enum_antivirus_engine_task.EngineTypeClamav
	param.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")

	loginInfo := u.getLogin(ctx)
	param.CreateId = loginInfo.AdminId
	param.CreateName = loginInfo.Username

	if errCode, err := antivirus_engine_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// ClearForManage 管理: 清空
func (u *AntivirusEngineController) ClearForManage(ctx *gin.Context) {
	if !software_service.CreateCurdFactory().Clear() {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// DownloadForManage 管理: 下载
func (u *AntivirusEngineController) DownloadForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	model := software_service.CreateCurdFactory().GetItemById(softId)
	if model == nil || model.Id == 0 {
		response.FailSimplify(ctx, "软件不存在")
		return
	}

	response.FileDownloadAndRename(ctx, model.FullPath, enum_software.DownloadName(model.PlatformType, model.OsType, model.Ver))
}

// PublishForManage 管理: 发布
func (u *AntivirusEngineController) PublishForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().PublishById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// CancelForManage 管理: 取消
func (u *AntivirusEngineController) CancelForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().CancelById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理: 删除
func (u *AntivirusEngineController) DeleteForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().DeleteById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
