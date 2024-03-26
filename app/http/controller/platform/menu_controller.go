package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	enum_enum "goskeleton/app/global/enums/enum_menu"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/munu_service"
	"goskeleton/app/utils/response"
)

// MenuController 角色相关
type MenuController struct {
	ControllerBase
}

// Tree 树形列表
func (u *MenuController) Tree(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Tree")
}

// Add 添加
func (u *MenuController) Add(ctx *gin.Context) {
	var param = model_platform.MenuModel{
		ParentId:   int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "parentId")),
		Position:   enum_enum.PositionLeft,
		Title:      ctx.GetString(consts.ValidatorPrefixPlatform + "name"),
		WebCode:    ctx.GetString(consts.ValidatorPrefixPlatform + "code"),
		Sort:       int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "sort")),
		IsShow:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isShow")),
		IsRoot:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isRoot")),
		Remark:     ctx.GetString(consts.ValidatorPrefixPlatform + "remark"),
		ServiceApi: "",
	}

	errCode, err := munu_service.CreateCurdFactory().Add(&param)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// Edit 编辑
func (u *MenuController) Edit(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "edit")
}

// Sort 排序
func (u *MenuController) Sort(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "sort")
}

// Delete 删除
func (u *MenuController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "delete")
}
