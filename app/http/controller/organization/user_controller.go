package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/service/organization/user_service"
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
	response.SuccessSimplify(ctx)
}

// Add 添加
func (u *UserController) Add(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Username:   ctx.GetString(consts.ValidatorPrefixPlatform + "userName"),
		OrgId:      loginInfo.OrgId,
		Email:      ctx.GetString(consts.ValidatorPrefixPlatform + "userEmail"),
		Phone:      ctx.GetString(consts.ValidatorPrefixPlatform + "userPhone"),
		Password:   u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "userPw")),
		Realname:   ctx.GetString(consts.ValidatorPrefixPlatform + "realname"),
		Position:   ctx.GetString(consts.ValidatorPrefixPlatform + "position"),
		Intro:      ctx.GetString(consts.ValidatorPrefixPlatform + "intro"),
		Depts:      ctx.GetString(consts.ValidatorPrefixPlatform + "departments"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}
	if string_helper.IsEmpty(param.Password) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	if err := user_service.CreateCurdFactory().AddSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Edit 编辑
func (u *UserController) Edit(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := user_models.UserModel{
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id")),
		OrgId:      loginInfo.OrgId,
		Email:      ctx.GetString(consts.ValidatorPrefixPlatform + "userEmail"),
		Phone:      ctx.GetString(consts.ValidatorPrefixPlatform + "userPhone"),
		Realname:   ctx.GetString(consts.ValidatorPrefixPlatform + "realname"),
		Position:   ctx.GetString(consts.ValidatorPrefixPlatform + "position"),
		Intro:      ctx.GetString(consts.ValidatorPrefixPlatform + "intro"),
		Depts:      ctx.GetString(consts.ValidatorPrefixPlatform + "departments"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := user_service.CreateCurdFactory().UpdateSub(&param, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Delete 删除
func (u *UserController) Delete(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id"))

	if !department_service.CreateCurdFactory().DelSelfWithChildById(loginInfo.OrgId, id) {
		response.FailSimplify(ctx, "删除失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// ModifyPw 修改密码
func (u *UserController) ModifyPw(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id"))

	if !department_service.CreateCurdFactory().DelSelfWithChildById(loginInfo.OrgId, id) {
		response.FailSimplify(ctx, "删除失败")
		return
	}

	response.SuccessSimplify(ctx)
}
