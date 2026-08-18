package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/service/organization/role_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/org_menu_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
	"strings"
)

// AccountController 首页相关
type AccountController struct {
	ControllerBase
}

func (u *AccountController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// GetUser 获取用户信息
func (u *AccountController) GetUser(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	model := user_service.CreateCurdFactory().GetItemByOrgAndUser(loginInfo.OrgId, loginInfo.UserId)

	info := (&org_trans.User{}).TransInfo(model, nil, nil)
	info["org_name"] = loginInfo.OrgName
	info["org_code"] = loginInfo.OrgCode
	response.SuccessReturnData(ctx, info)
}

// ModifyUser 修改用户信息
func (u *AccountController) ModifyUser(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := user_models.UserModel{
		Id:         loginInfo.UserId,
		OrgId:      loginInfo.OrgId,
		Position:   ctx.GetString(consts.ValidatorPrefixOrg + "position"),
		Intro:      ctx.GetString(consts.ValidatorPrefixOrg + "intro"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := user_service.CreateCurdFactory().ModifyInfoById(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetMenuTree 获取菜单菜单树
func (u *AccountController) GetMenuTree(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	IsYes := enums.IsYes
	conditions := org_menu_models.SearchConditions{
		UsageVerList: []int{
			0, //全部版本
			cur_app.GetAppUsageVersion(),
		},
		IsShow: &IsYes,
	}

	//不是root
	if loginInfo.IsRoot == enums.IsNot {
		//只能看非root菜单
		IsNot := enums.IsNot
		conditions.IsRoot = &IsNot

		//角色菜单
		if roleModel := role_service.CreateCurdFactory().HasId(loginInfo.RoleId, loginInfo.OrgId); roleModel != nil {
			for _, val := range strings.Split(roleModel.Menus, ",") {
				conditions.Ids = append(conditions.Ids, string_helper.ToInt64(val))
			}
		} else {
			conditions.Ids = append(conditions.Ids, 0)
		}
	}

	models := org_menu_models.CreateOrgMenuFactory().GetListByConditions(conditions)
	if models == nil {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}

	response.SuccessReturnData(ctx, org_menu_service.CreateCurdFactory().ToTree(models, nil))
}

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := user_models.UserModel{
		Id:         loginInfo.UserId,
		OrgId:      loginInfo.OrgId,
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw")),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	newPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "newPw"))
	if string_helper.IsEmpty(param.Password) || string_helper.IsEmpty(newPw) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if err := user_service.CreateCurdFactory().ModifyPwById(&param, newPw, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyPhone 修改手机
func (u *AccountController) ModifyPhone(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := user_models.UserModel{
		Id:         loginInfo.UserId,
		OrgId:      loginInfo.OrgId,
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw")),
		Phone:      ctx.GetString(consts.ValidatorPrefixOrg + "phone"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	if string_helper.IsEmpty(param.Password) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if err := user_service.CreateCurdFactory().ModifyPhoneById(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyEmail 修改邮箱
func (u *AccountController) ModifyEmail(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := user_models.UserModel{
		Id:         loginInfo.UserId,
		OrgId:      loginInfo.OrgId,
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw")),
		Email:      ctx.GetString(consts.ValidatorPrefixOrg + "email"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	if string_helper.IsEmpty(param.Password) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if err := user_service.CreateCurdFactory().ModifyEmailById(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
