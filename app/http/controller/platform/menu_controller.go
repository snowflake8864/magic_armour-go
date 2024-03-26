package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_menu"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/munu_service"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/tree"
)

// MenuController 角色相关
type MenuController struct {
	ControllerBase
}

// Tree 树形列表
func (u *MenuController) Tree(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	models := munu_service.CreateCurdFactory().GetAll(loginInfo.IsRoot)
	if models == nil {
		response.FailSimplify(ctx, "没有数据")
		return
	}

	var nodes []*tree.Node
	for _, model := range models {
		nodes = append(nodes, &tree.Node{
			ID:       model.Id,
			ParentId: model.ParentId,
			Data: map[string]interface{}{
				"name":       model.Name,
				"type":       model.Type,
				"code":       model.Code,
				"sort":       model.Sort,
				"is_show":    model.IsShow,
				"is_root":    model.IsRoot,
				"remark":     model.Remark,
				"api":        model.Api,
				"level":      model.Level,
				"created_at": model.CreatedAt,
			},
			Child: []*tree.Node{},
		})
	}

	response.SuccessReturnData(ctx, tree.RealTree(nodes))
}

// Add 添加
func (u *MenuController) Add(ctx *gin.Context) {
	var param = model_platform.MenuModel{
		Position: enum_menu.PositionLeft,
		ParentId: int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "parentId")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "type")),
		Name:     ctx.GetString(consts.ValidatorPrefixPlatform + "name"),
		Code:     ctx.GetString(consts.ValidatorPrefixPlatform + "code"),
		Sort:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "sort")),
		IsShow:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isShow")),
		IsRoot:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isRoot")),
		Remark:   ctx.GetString(consts.ValidatorPrefixPlatform + "remark"),
		Api:      "",
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
	var param = model_platform.MenuModel{
		Id:     int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id")),
		Type:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "type")),
		Name:   ctx.GetString(consts.ValidatorPrefixPlatform + "name"),
		Code:   ctx.GetString(consts.ValidatorPrefixPlatform + "code"),
		Sort:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "sort")),
		IsShow: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isShow")),
		IsRoot: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isRoot")),
		Remark: ctx.GetString(consts.ValidatorPrefixPlatform + "remark"),
	}
	errCode, err := munu_service.CreateCurdFactory().UpdateById(&param)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// Sort 排序
func (u *MenuController) Sort(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "sort")
}

// Delete 删除
func (u *MenuController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "delete")
}
