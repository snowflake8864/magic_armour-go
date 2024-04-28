package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/http/transformer"
	"goskeleton/app/utils/response"
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
	response.SuccessReturnData(ctx, "GetForConfig")
}

// SaveForConfig 配置: 保存
func (u *ThreatenController) SaveForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForConfig")
}

// OnlineForConfig 配置: 验证在线
func (u *ThreatenController) OnlineForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "OnlineForConfig")
}

// PageForCache 缓存: 分页
func (u *ThreatenController) PageForCache(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForCache")
}

// ClearForCache 缓存: 清空
func (u *ThreatenController) ClearForCache(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearForCache")
}

// RefreshForCache 缓存: 刷新
func (u *ThreatenController) RefreshForCache(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RefreshForCache")
}

// DeleteForCache 缓存: 删除
func (u *ThreatenController) DeleteForCache(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForCache")
}

// GetAttrForCache 缓存: 获取属性
func (u *ThreatenController) GetAttrForCache(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetAttrForCache")
}

// PageForPush 推送: 分页
func (u *ThreatenController) PageForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// AllWhiteForPush 推送: 全部白
func (u *ThreatenController) AllWhiteForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AllWhite")
}

// BatchWhiteForPush 推送: 批量白
func (u *ThreatenController) BatchWhiteForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhite")
}

// BatchBlackForPush 推送: 批量黑
func (u *ThreatenController) BatchBlackForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlack")
}

// ClearForPush 推送: 清空
func (u *ThreatenController) ClearForPush(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Clear")
}
