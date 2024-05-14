package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_menu"
	"goskeleton/app/model/organization/org_menu_models"
	"goskeleton/app/service/platform/org_menu_service"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/tree"
)

// OrgMenuController 企业相关
type OrgMenuController struct {
	ControllerBase
}

// Tree 树形列表
func (u *OrgMenuController) Tree(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	models := org_menu_service.CreateCurdFactory().GetAll(loginInfo.IsRoot)
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
				"usage_ver":  model.UsageVer,
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
func (u *OrgMenuController) Add(ctx *gin.Context) {
	var param = org_menu_models.OrgMenuModel{
		Position: enum_menu.PositionLeft,
		ParentId: int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "parentId")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "type")),
		Name:     ctx.GetString(consts.ValidatorPrefixPlatform + "name"),
		Code:     ctx.GetString(consts.ValidatorPrefixPlatform + "code"),
		Sort:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "sort")),
		IsShow:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isShow")),
		IsRoot:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isRoot")),
		UsageVer: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "usageVer")),
		Remark:   ctx.GetString(consts.ValidatorPrefixPlatform + "remark"),
		Api:      "",
	}

	errCode, err := org_menu_service.CreateCurdFactory().Add(&param)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// Edit 编辑
func (u *OrgMenuController) Edit(ctx *gin.Context) {
	var param = org_menu_models.OrgMenuModel{
		Id:       int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id")),
		Type:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "type")),
		Name:     ctx.GetString(consts.ValidatorPrefixPlatform + "name"),
		Code:     ctx.GetString(consts.ValidatorPrefixPlatform + "code"),
		Sort:     int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "sort")),
		IsShow:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isShow")),
		IsRoot:   int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isRoot")),
		UsageVer: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "usageVer")),
		Remark:   ctx.GetString(consts.ValidatorPrefixPlatform + "remark"),
	}
	errCode, err := org_menu_service.CreateCurdFactory().UpdateById(&param)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// Sort 排序
func (u *OrgMenuController) Sort(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "sort")
}

// Delete 删除
func (u *OrgMenuController) Delete(ctx *gin.Context) {
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id"))
	errCode, err := org_menu_service.CreateCurdFactory().DeleteById(id)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
