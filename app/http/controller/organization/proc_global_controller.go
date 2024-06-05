package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/response"
)

// ProcessGlobalController 进程全局名单相关
type ProcessGlobalController struct {
	ControllerBase
}

// PageFromProcess 分页
func (u *ProcessGlobalController) PageFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageFromProcess")
}

// PageFromDll 分页
func (u *ProcessGlobalController) PageFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageFromDll")
}

// AddFromProcess 添加
func (u *ProcessGlobalController) AddFromProcess(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.DataFromProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// AddFromDll 添加
func (u *ProcessGlobalController) AddFromDll(ctx *gin.Context) {
	if err := u.add(ctx, enum_process.DataFromDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// add
func (u *ProcessGlobalController) add(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)

	param := proc_models.ProcGlobalModel{
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		OrgId:      loginInfo.OrgId,
		Dir:        ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Type:       int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		DataFrom:   dataFrom,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	return proc_service.CreateGlobalCurdFactory().Add(&param)
}

// add
func (u *ProcessGlobalController) batchAdd(ctx *gin.Context, dataFrom int, globalType int) error {
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "list")
	var list []proc_models.ProcGlobalModel
	if err := json_helper.ToAny(jsonStr, &list); err != nil {
		return errors.New("参数错误")
	}

	loginInfo := u.getLogin(ctx)
	param := proc_models.ProcGlobalModel{
		OrgId:      loginInfo.OrgId,
		Type:       globalType,
		DataFrom:   dataFrom,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	return proc_service.CreateGlobalCurdFactory().BatchAdd(&param, &list)
}

// ClearFromProcess 清空
func (u *ProcessGlobalController) ClearFromProcess(ctx *gin.Context) {
	if err := u.batchAdd(ctx, enum_process.DataFromProcess, enum_process.GlobalTypeWhite); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ClearFromDll 清空
func (u *ProcessGlobalController) ClearFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearFromDll")
}

// BatchWhiteFromProcess 批量白来源进程
func (u *ProcessGlobalController) BatchWhiteFromProcess(ctx *gin.Context) {
	if err := u.batchAdd(ctx, enum_process.DataFromProcess, enum_process.GlobalTypeWhite); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchBlackFromProcess 批量黑来源进程
func (u *ProcessGlobalController) BatchBlackFromProcess(ctx *gin.Context) {
	if err := u.batchAdd(ctx, enum_process.DataFromProcess, enum_process.GlobalTypeBlack); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchWhiteFromDll 批量白来源模块
func (u *ProcessGlobalController) BatchWhiteFromDll(ctx *gin.Context) {
	if err := u.batchAdd(ctx, enum_process.DataFromDll, enum_process.GlobalTypeWhite); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchBlackFromDll 批量黑来源模块
func (u *ProcessGlobalController) BatchBlackFromDll(ctx *gin.Context) {
	if err := u.batchAdd(ctx, enum_process.DataFromDll, enum_process.GlobalTypeBlack); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetSwitch 编辑
func (u *ProcessGlobalController) SetSwitch(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Edit")
}

// Delete 删除
func (u *ProcessGlobalController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Delete")
}
