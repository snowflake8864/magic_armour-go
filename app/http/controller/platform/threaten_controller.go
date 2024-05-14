package platform

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
	"strings"
)

// ThreatenController 威胁情报推送相关
type ThreatenController struct {
	ControllerBase
}

// Levels 等级
func (u *ThreatenController) Levels(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("level", enum_threaten.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("level", enum_threaten.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// GetForConfig 配置: 获取
func (u *ThreatenController) GetForConfig(ctx *gin.Context) {
	dictIds := []string{
		enum_dictionary.ThreatenSwitch,
		enum_dictionary.ThreatenUrl,
		enum_dictionary.ThreatenAppId,
		enum_dictionary.ThreatenAppKey,
		enum_dictionary.ThreatenCacheExpire,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	response.SuccessReturnData(ctx, vals)
}

// SaveForConfig 配置: 保存
func (u *ThreatenController) SaveForConfig(ctx *gin.Context) {
	idAndVal := map[string]interface{}{
		enum_dictionary.ThreatenSwitch:      ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_switch"),
		enum_dictionary.ThreatenUrl:         strings.TrimRight(ctx.GetString(consts.ValidatorPrefixPlatform+"threaten_url"), "/"),
		enum_dictionary.ThreatenAppId:       ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_app_id"),
		enum_dictionary.ThreatenAppKey:      ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_app_key"),
		enum_dictionary.ThreatenCacheExpire: ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_cache_expire"),
	}

	if !dictionary_service.CreateCurdFactory().BatchUpdate(idAndVal) {
		response.FailSimplify(ctx, "保持失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// OnlineForConfig 配置: 验证在线
func (u *ThreatenController) OnlineForConfig(ctx *gin.Context) {
	config := threaten_platform.RequestConfig{
		Domain:    ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_url"),
		AppId:     ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_app_id"),
		AppSecret: ctx.GetString(consts.ValidatorPrefixPlatform + "threaten_app_key"),
	}
	threatenPlatform := threaten_platform.CreateRequestFactory(config)
	if err := threatenPlatform.Online(); err != nil {
		response.FailData(ctx, "网络连接失败", err.Error())
		return
	}
	if err := threatenPlatform.GenerateToken(); err != nil {
		response.FailData(ctx, "配置不正确", err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForCache 缓存: 分页
func (u *ThreatenController) PageForCache(ctx *gin.Context) {
	//页面上传
	conditions := threaten_service.SearchConditions{}
	conditions.Id = ctx.GetString(consts.ValidatorPrefixPlatform + "id")
	conditions.Dir = ctx.GetString(consts.ValidatorPrefixPlatform + "dir")
	conditions.Level = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "level"))

	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//分页
	total, models := threaten_service.CreateCacheFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Threaten{}).TransCachePage(total, models))
}

// ClearForCache 缓存: 清空
func (u *ThreatenController) ClearForCache(ctx *gin.Context) {
	if !threaten_service.CreateCacheFactory().Clear() {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// RefreshForCache 缓存: 刷新
func (u *ThreatenController) RefreshForCache(ctx *gin.Context) {
	if !threaten_service.CreateCacheFactory().Refresh() {
		response.FailSimplify(ctx, "刷新失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForCache 缓存: 删除
func (u *ThreatenController) DeleteForCache(ctx *gin.Context) {
	id := ctx.GetString(consts.ValidatorPrefixPlatform + "id")

	if !threaten_service.CreateCacheFactory().DeleteById(id) {
		response.FailSimplify(ctx, "删除失败")
		return
	}
	response.SuccessSimplify(ctx)
}

// GetAttrForCache 缓存: 获取属性
func (u *ThreatenController) GetAttrForCache(ctx *gin.Context) {
	id := ctx.GetString(consts.ValidatorPrefixPlatform + "id")
	model := threaten_service.CreateCacheFactory().GetItemById(id)
	if model == nil || helpers.StringIsEmpty(model.Id) {
		response.FailSimplify(ctx, "数据不存在")
		return
	}
	response.SuccessReturnData(ctx, (&platform_trans.Threaten{}).TransThreatenInfo(model))
}

// PageForPush 推送: 分页
func (u *ThreatenController) PageForPush(ctx *gin.Context) {
	//页面上传
	conditions := threaten_service.SearchConditions{}
	conditions.Id = ctx.GetString(consts.ValidatorPrefixPlatform + "id")
	conditions.Dir = ctx.GetString(consts.ValidatorPrefixPlatform + "dir")

	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//分页
	total, models := threaten_service.CreatePushFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Threaten{}).TransPushPage(total, models))
}

// AllWhiteForPush 推送: 全部白
func (u *ThreatenController) AllWhiteForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AllWhite")
}

// BatchWhiteForPush 推送: 批量白
func (u *ThreatenController) BatchWhiteForPush(ctx *gin.Context) {
	if errCode, err := u.handelBatchForPush(ctx, enum_threaten.LevelBlack); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchBlackForPush 推送: 批量黑
func (u *ThreatenController) BatchBlackForPush(ctx *gin.Context) {
	if errCode, err := u.handelBatchForPush(ctx, enum_threaten.LevelBlack); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchBlackForPush 推送: 批量黑
func (u *ThreatenController) handelBatchForPush(ctx *gin.Context, level int) (errCode int, err error) {
	list := ctx.GetString(consts.ValidatorPrefixPlatform + "list")
	var params []threaten_models.ThreatenPushModel
	err = helpers.Json2Any(list, &params)
	if err != nil {
		return -1, errors.New("参数格式错误")
	}
	loginInfo := u.getLogin(ctx)
	return threaten_service.CreatePushFactory().BatchAdd(params, level, loginInfo.AdminId, loginInfo.Username)
}

// ClearForPush 推送: 清空
func (u *ThreatenController) ClearForPush(ctx *gin.Context) {
	if !threaten_service.CreatePushFactory().Clear() {
		response.FailSimplify(ctx, "清空失败")
		return
	}
	response.SuccessSimplify(ctx)
}
