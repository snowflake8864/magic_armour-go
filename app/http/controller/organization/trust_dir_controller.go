package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/trust_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/agent/trust_dir_service"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/response"
)

// TrustDirController 信任目录相关
type TrustDirController struct {
	ControllerBase
}

// PageForGlobal 全局： 分页
func (u *TrustDirController) PageForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := trust_models.SearchConditions{
		OrgId: loginInfo.OrgId,
		Name:  ctx.GetString(consts.ValidatorPrefixOrg + "name"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := trust_models.CreateTrustGlobalFactory().Page(condition, offset, pageSize, order)

	//获取组织架构
	departmentMaps := map[int64]*department_models.DepartmentModel{}
	if models != nil && len(models) > 0 {
		var deptIds []int64
		for _, model := range models {
			if model.DeptId > 0 {
				deptIds = append(deptIds, model.DeptId)
			}
		}
		departmentMaps = u.getDepartmentMap(loginInfo.OrgId, slice_helper.Int64RemoveRepeatByMap(deptIds))
	}

	response.SuccessReturnData(ctx, (&org_trans.TrustDirTrans{}).TransGlobalPage(total, models, departmentMaps))
}

// AddForGlobal 全局： 添加
func (u *TrustDirController) AddForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := trust_models.TrustGlobalModel{
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		Dir:          ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Name:         ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		IsExtend:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend")),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
	}

	if err := u.isBelongUserDepartment(loginInfo.IsRoot, &loginInfo.DeptWithChild, param.DeptId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	if err := trust_dir_service.CreateGlobalCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForGlobal 全局： 编辑
func (u *TrustDirController) EditForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := trust_models.TrustGlobalModel{
		Id:           int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		DeptId:       int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "deptId")),
		Name:         ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		IsExtend:     int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExtend")),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
	}

	if err := u.isBelongUserDepartment(loginInfo.IsRoot, &loginInfo.DeptWithChild, param.DeptId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	if err := trust_dir_service.CreateGlobalCurdFactory().Edit(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetSwitchForGlobal 全局： 设置开关
func (u *TrustDirController) SetSwitchForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := trust_models.TrustGlobalModel{
		Id:           int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		OpenSwitch:   int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch")),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
	}

	if err := trust_dir_service.CreateGlobalCurdFactory().SetSwitch(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForGlobal 全局： 删除
func (u *TrustDirController) DeleteForGlobal(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := trust_models.TrustGlobalModel{
		Id:           int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:        loginInfo.OrgId,
		PlatformType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "platformType")),
		CreateId:     loginInfo.UserId,
		CreateName:   loginInfo.Username,
	}

	if err := trust_dir_service.CreateGlobalCurdFactory().DeleteById(&param, loginInfo.IsRoot); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
