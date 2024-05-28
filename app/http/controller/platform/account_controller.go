package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/platform/admin_service"
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

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "oldPw"))
	newPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "newPw"))
	if string_helper.IsEmpty(oldPw) || string_helper.IsEmpty(newPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := admin_service.CreateCurdFactory().ModifyPwById(loginInfo.AdminId, oldPw, newPw); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
