package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers"
)

// ControllerBase 基础
type ControllerBase struct {
}

// getLogin 获取登录数据
func (u *ControllerBase) getLogin(ctx *gin.Context) admin_service.AdminCacheLogin {
	return cur_userinfo.GetCurrentAdmin(ctx)
}

// getOrgSearchConditions 获取企业查询条件
func (u *ControllerBase) getOrgSearchConditions(ctx *gin.Context) (organize_service.SearchConditions, error) {
	//页面上传
	searchConditions := organize_service.SearchConditions{}
	searchConditions.OraName = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	searchConditions.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")
	searchConditions.UserPhone = ctx.GetString(consts.ValidatorPrefixPlatform + "userPhone")
	searchConditions.OrgType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgType"))
	return searchConditions, nil
}

// getPageAndSort 获取分页参数和排序参数
func (u *ControllerBase) getPageAndSort(ctx *gin.Context, defaultField string) (offset, pageSize int, order string) {
	page := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "page"))
	pageSize = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "pageSize"))
	orderField := ctx.GetString(consts.ValidatorPrefixPlatform + "field")
	orderSort := ctx.GetString(consts.ValidatorPrefixPlatform + "sort")
	if page <= 0 {
		page = 1
	}
	if pageSize <= 1 {
		page = 20
	}
	offset = (page - 1) * pageSize

	if helpers.StringIsEmpty(orderField) {
		if helpers.StringIsEmpty(defaultField) {
			orderField = "created_at"
		} else {
			orderField = defaultField
		}
	}
	if helpers.StringIsEmpty(orderSort) {
		orderSort = "DESC"
	}
	return offset, pageSize, orderField + " " + orderSort
}

// getWebSliceOfInt64 获取前端上传数组: int64
func (u *ControllerBase) getWebSliceOfInt64(ctx *gin.Context, field string) []int64 {
	var list []int64
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixPlatform + field); isOk {
		var tmp []interface{}
		tmp, _ = deptIds.([]interface{})
		list = helpers.SliceAny2Int64(tmp)
	}
	return list
}

// getWebSliceOfInt 获取前端上传数组: int64
func (u *ControllerBase) getWebSliceOfInt(ctx *gin.Context, field string) []int {
	var list []int
	if deptIds, isOk := ctx.Get(consts.ValidatorPrefixPlatform + field); isOk {
		var tmp []interface{}
		tmp, _ = deptIds.([]interface{})
		list = helpers.SliceAny2Int(tmp)
	}
	return list
}
