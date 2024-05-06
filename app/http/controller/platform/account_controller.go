package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// AccountController 首页相关
type AccountController struct {
	ControllerBase
}

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ModifyPw")
}
