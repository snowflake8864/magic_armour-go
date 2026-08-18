package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// UserController 组织架构相关
type UserController struct {
	ControllerBase
}

func (u *UserController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// Page 分页
func (u *UserController) Page(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	isRoot := enums.IsNot
	conditions := user_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		IsRoot: &isRoot,
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")
	//分页
	total, userModels := user_models.CreateUserFactory().Page(conditions, offset, pageSize, order)

	departmentMaps := make(map[int64]*department_models.DepartmentModel)
	roleMaps := make(map[int64]*role_models.RoleModel)
	if userModels != nil && len(userModels) > 0 {
		var deptIds []int64
		var roleIds []int64
		for _, model := range userModels {
			deptIds = append(deptIds, json_helper.ToSliceInt64(model.Depts)...)
			roleIds = append(roleIds, model.RoleId)
		}
		departmentMaps = u.getDepartmentMap(loginInfo.OrgId, deptIds)
		roleMaps = role_models.CreateRoleFactory().GetMapsByOrgAndIds(loginInfo.OrgId, roleIds)
	}

	response.SuccessReturnData(ctx, (&org_trans.User{}).TransPage(total, userModels, departmentMaps, roleMaps))
}

// Add 添加
func (u *UserController) Add(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Username:   ctx.GetString(consts.ValidatorPrefixOrg + "userName"),
		OrgId:      loginInfo.OrgId,
		Email:      ctx.GetString(consts.ValidatorPrefixOrg + "userEmail"),
		Phone:      ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"),
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPw")),
		Realname:   ctx.GetString(consts.ValidatorPrefixOrg + "realname"),
		Position:   ctx.GetString(consts.ValidatorPrefixOrg + "position"),
		Intro:      ctx.GetString(consts.ValidatorPrefixOrg + "intro"),
		Depts:      ctx.GetString(consts.ValidatorPrefixOrg + "departments"),
		RoleId:     int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "roleId")),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	if string_helper.IsEmpty(param.Password) {
		response.FailSimplify(ctx, "my_errors.CurdInvalidParameter")
		return
	}
	if err := user_service.CreateCurdFactory().AddForSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Edit 编辑
func (u *UserController) Edit(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:      loginInfo.OrgId,
		Email:      ctx.GetString(consts.ValidatorPrefixOrg + "userEmail"),
		Phone:      ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"),
		Realname:   ctx.GetString(consts.ValidatorPrefixOrg + "realname"),
		Position:   ctx.GetString(consts.ValidatorPrefixOrg + "position"),
		Intro:      ctx.GetString(consts.ValidatorPrefixOrg + "intro"),
		Depts:      ctx.GetString(consts.ValidatorPrefixOrg + "departments"),
		RoleId:     int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "roleId")),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := user_service.CreateCurdFactory().UpdateForSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Delete 删除
func (u *UserController) Delete(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:      loginInfo.OrgId,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := user_service.CreateCurdFactory().DeleteForSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// ModifyPw 修改密码
func (u *UserController) ModifyPw(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		OrgId:      loginInfo.OrgId,
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPw")),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	if string_helper.IsEmpty(param.Password) {
		response.FailSimplify(ctx, "my_errors.CurdInvalidParameter")
		return
	}

	if err := user_service.CreateCurdFactory().ModifyPwForSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
