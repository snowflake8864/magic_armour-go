package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/platform/antivirus_engine_models"
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

	condition := antivirus_engine_models.SearchConditions{}

	condition.Ver = ctx.GetString(consts.ValidatorPrefixPlatform + "ver")
	condition.CreateAtGte = ctx.GetString(consts.ValidatorPrefixPlatform + "startTime")
	condition.CreateAtLte = ctx.GetString(consts.ValidatorPrefixPlatform + "endTime")
	condition.PlatformType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))

	//分页参数和排序参数
	offset, pageSize, _ := u.getPageAndSort(ctx, "")

	//软件分页
	total, models := antivirus_engine_models.CreateEngineFactory().Page(condition, offset, pageSize, "")

	response.SuccessReturnData(ctx, (&platform_trans.AntivirusEngine{}).TransPage(total, models))
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

// PublishForManage 管理: 发布
func (u *AntivirusEngineController) PublishForManage(ctx *gin.Context) {
	engineId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "engineId"))
	isPublish := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isPublish"))

	if isPublish == enums.IsYes {
		if errCode, err := antivirus_engine_service.CreateCurdFactory().PublishById(engineId); err != nil {
			response.FailRoutine(ctx, errCode, err.Error())
			return
		}
	} else {
		if errCode, err := antivirus_engine_service.CreateCurdFactory().CancelById(engineId); err != nil {
			response.FailRoutine(ctx, errCode, err.Error())
			return
		}
	}

	response.SuccessSimplify(ctx)
}

// NewestForManage 管理: 设置最新
func (u *AntivirusEngineController) NewestForManage(ctx *gin.Context) {
	engineId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "engineId"))
	if errCode, err := antivirus_engine_service.CreateCurdFactory().NewestById(engineId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理: 删除
func (u *AntivirusEngineController) DeleteForManage(ctx *gin.Context) {
	engineId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "engineId"))
	if errCode, err := antivirus_engine_service.CreateCurdFactory().DeleteById(engineId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
