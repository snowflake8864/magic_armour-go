package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// AccountController 首页相关
type AccountController struct {
	ControllerBase
}

func (u *AccountController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// GetMenuTree 获取菜单菜单树
func (u *AccountController) GetMenuTree(ctx *gin.Context) {
	response.SuccessSimplify(ctx)
}

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	newPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "newPw"))
	if string_helper.IsEmpty(oldPw) || string_helper.IsEmpty(newPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyPwById(loginInfo.UserId, oldPw, newPw); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyPhone 修改手机
func (u *AccountController) ModifyPhone(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	phone := ctx.GetString(consts.ValidatorPrefixOrg + "phone")
	if string_helper.IsEmpty(oldPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyPhoneById(loginInfo.UserId, oldPw, phone); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyEmail 修改邮箱
func (u *AccountController) ModifyEmail(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	email := ctx.GetString(consts.ValidatorPrefixOrg + "email")
	if string_helper.IsEmpty(oldPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyEmailById(loginInfo.UserId, oldPw, email); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
