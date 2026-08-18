package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_menu"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/service/platform/org_menu_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// OrgMenuController 企业相关
type OrgMenuController struct {
	ControllerBase
}

// AppVers 应用版本
func (u *OrgMenuController) AppVers(ctx *gin.Context) {
	var list []map[string]interface{}

	list = append(list, map[string]interface{}{
		"val":   0,
		"label": "全部",
	}, map[string]interface{}{
		"val":   variable.AppUsageVersionCloud,
		"label": "云版本",
	}, map[string]interface{}{
		"val":   variable.AppUsageVersionOffline,
		"label": "离线版本",
	})

	response.SuccessReturnData(ctx, list)
}

// Tree 树形列表
func (u *OrgMenuController) Tree(ctx *gin.Context) {
	conditions := org_menu_models.SearchConditions{}
	models := org_menu_models.CreateOrgMenuFactory().GetListByConditions(conditions)
	if models == nil {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}

	response.SuccessReturnData(ctx, org_menu_service.CreateCurdFactory().ToTree(models, nil, cur_userinfo.GetLocales(ctx)))
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
