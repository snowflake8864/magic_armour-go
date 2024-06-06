package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/agent/trust_models"
	"goskeleton/app/service/agent/trust_dir_service"
	"goskeleton/app/utils/response"
)

// TrustDirController 信任目录相关
type TrustDirController struct {
	ControllerBase
}

// PageForGlobal 全局： 分页
func (u *TrustDirController) PageForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
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
	response.SuccessReturnData(ctx, "FilePowers")
}

// SetSwitchForGlobal 全局： 设置开关
func (u *TrustDirController) SetSwitchForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// DeleteForGlobal 全局： 删除
func (u *TrustDirController) DeleteForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}
