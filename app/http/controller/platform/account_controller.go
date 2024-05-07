package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/response"
)

// AccountController 首页相关
type AccountController struct {
	ControllerBase
}

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	oldPw := ctx.GetString(consts.ValidatorPrefixPlatform + "oldPw")
	newPw := ctx.GetString(consts.ValidatorPrefixPlatform + "newPw")

	loginInfo := u.getLogin(ctx)
	if errCode, err := admin_service.CreateCurdFactory().ModifyPwById(loginInfo.AdminId, oldPw, newPw); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
