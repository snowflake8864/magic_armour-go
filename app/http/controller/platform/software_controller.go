package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/platform/software_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/response"
)

// SoftwareController 终端相关
type SoftwareController struct {
	ControllerBase
}

// GetOsTypes  获取系统类型
func (u *SoftwareController) GetOsTypes(ctx *gin.Context) {
	platformType := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))

	osTypes := enums.Get1stIsStrAnd2ndIsIntReturnSliceInt("map_platform_os", platformType, enum_software.EnumSlice)
	if osTypes == nil {
		response.FailSimplify(ctx, "platformType no exist")
		return
	}

	var result []map[string]interface{}
	for _, osType := range osTypes {

		//先获取系统类型保护的版本
		var children []map[string]interface{}
		typeVers := enums.Get1stIsStrAnd2ndIsIntReturnSliceInt("map_os_ver", osType, enum_software.EnumSlice)
		for _, typeVer := range typeVers {
			children = append(children, map[string]interface{}{
				"val":   typeVer,
				"label": enums.Get1stIsStrAnd2ndIsIntReturnStr("os_ver", typeVer, enum_software.EnumMap),
			})
		}

		//系统类型
		result = append(result, map[string]interface{}{
			"val":      osType,
			"label":    enums.Get1stIsStrAnd2ndIsIntReturnStr("os_type", osType, enum_software.EnumMap),
			"children": children,
		})
	}

	response.SuccessReturnData(ctx, result)
}

// AgentConfigGet  终端配置:获取
func (u *SoftwareController) AgentConfigGet(ctx *gin.Context) {
	dictIds := []string{
		enum_dictionary.AgentHttpUrl,
		enum_dictionary.AgentHttpPort,
		enum_dictionary.AgentCompatibleDefaultOrg,
		enum_dictionary.AgentCompatibleDefaultUser,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, gin.H{
		"http_url":     vals[enum_dictionary.AgentHttpUrl],
		"http_port":    vals[enum_dictionary.AgentHttpPort],
		"default_org":  vals[enum_dictionary.AgentCompatibleDefaultOrg],
		"default_user": vals[enum_dictionary.AgentCompatibleDefaultUser],
	})
}

// AgentConfigSet  终端配置:保存
func (u *SoftwareController) AgentConfigSet(ctx *gin.Context) {
	idAndVal := map[string]interface{}{
		enum_dictionary.AgentHttpUrl:               ctx.GetString(consts.ValidatorPrefixPlatform + "http_url"),
		enum_dictionary.AgentHttpPort:              ctx.GetString(consts.ValidatorPrefixPlatform + "http_port"),
		enum_dictionary.AgentCompatibleDefaultOrg:  ctx.GetString(consts.ValidatorPrefixPlatform + "default_org"),
		enum_dictionary.AgentCompatibleDefaultUser: ctx.GetString(consts.ValidatorPrefixPlatform + "default_user"),
	}

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保持失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// PageForManage 管理: 分页
func (u *SoftwareController) PageForManage(ctx *gin.Context) {

	//页面上传
	PlatformType := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))
	condition := software_models.SearchConditions{
		PlatformType: PlatformType,
		IsPublish:    enums.IsYes,
	}
	if PlatformType == enum_software.PlatformTypeWindows || PlatformType == enum_software.PlatformTypeLinux {
		condition.PlatformType = PlatformType
		condition.IsPublish = enums.IsYes
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//软件分页
	total, models := software_models.CreateSoftwareFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Software{}).TransPage(total, models))
}

// AddForManage 管理:  添加
func (u *SoftwareController) AddForManage(ctx *gin.Context) {
	var param = software_models.SoftwareModel{}
	param.PlatformType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "platformType"))
	param.OsType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "osType"))
	param.OsVer = json_helper.Encode(u.getWebSliceOfInt(ctx, "osVer"))
	param.UploadId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "uploadId"))
	param.Ver = ctx.GetString(consts.ValidatorPrefixPlatform + "versionNumber")
	param.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")

	loginInfo := u.getLogin(ctx)
	param.CreateId = loginInfo.AdminId
	param.CreateName = loginInfo.Username

	if errCode, err := software_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// ClearForManage 管理: 清空
func (u *SoftwareController) ClearForManage(ctx *gin.Context) {
	if !software_service.CreateCurdFactory().Clear() {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// DownloadForManage 管理: 下载
func (u *SoftwareController) DownloadForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	model := software_service.CreateCurdFactory().GetItemById(softId)
	if model == nil || model.Id == 0 {
		response.FailSimplify(ctx, "软件不存在")
		return
	}

	response.FileDownloadAndRename(ctx, model.FullPath, enum_software.DownloadName(model.PlatformType, model.OsType, model.Ver))
}

// PublishForManage 管理: 发布
func (u *SoftwareController) PublishForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().PublishById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// CancelForManage 管理: 取消
func (u *SoftwareController) CancelForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().CancelById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理: 删除
func (u *SoftwareController) DeleteForManage(ctx *gin.Context) {
	softId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "softId"))
	if errCode, err := software_service.CreateCurdFactory().DeleteById(softId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
