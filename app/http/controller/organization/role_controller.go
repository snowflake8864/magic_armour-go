package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/service/organization/role_service"
	"goskeleton/app/utils/response"
)

// RoleController 角色相关
type RoleController struct {
	ControllerBase
}

// AllForManage 管理页：全部角色
func (u *RoleController) AllForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := role_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	models := role_models.CreateRoleFactory().GetListByConditions(condition)
	response.SuccessReturnData(ctx, (&org_trans.Role{}).TransList(models))
}

// AddForManage 管理页：添加
func (u *RoleController) AddForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := role_models.RoleModel{
		OrgId:      loginInfo.OrgId,
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := role_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// EditForManage 管理页：编辑
func (u *RoleController) EditForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := role_models.RoleModel{
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:      loginInfo.OrgId,
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := role_service.CreateCurdFactory().UpdateById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理页：删除
func (u *RoleController) DeleteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForManage")
}

// LoginMenuTreeForPower 权限：用户登录后的菜单树
func (u *RoleController) LoginMenuTreeForPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "LoginMenuTreeForPower")
}

// EditableMenuTreeForPower 权限：可编辑的菜单树
func (u *RoleController) EditableMenuTreeForPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "EditableMenuTreeForPower")
}

// GetSelectedMenusForPower 权限：获取已选中的菜单集合
func (u *RoleController) GetSelectedMenusForPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetSelectedMenusForPower")
}

// SetSelectedMenusForPower 权限：保存已选中的菜单集合
func (u *RoleController) SetSelectedMenusForPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSelectedMenusForPower")
}
