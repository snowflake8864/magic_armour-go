package organization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer/org_trans"
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
	_ = u.page(ctx, enum_process.DataFromProcess)
}

// add
func (u *ProcessGlobalController) page(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := proc_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		GlobalType: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type")),
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		Dir:        ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		DataFrom:   dataFrom,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := proc_models.CreateProcGlobalFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.ProcGlobalTrans{}).TransPage(total, models))
	return nil
}

// PageFromDll 分页
func (u *ProcessGlobalController) PageFromDll(ctx *gin.Context) {
	_ = u.page(ctx, enum_process.DataFromDll)

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
	Type := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type"))
	if Type == enum_process.AllStatusGlobalWhite {
		Type = enum_process.GlobalTypeWhite
	} else {
		Type = enum_process.GlobalTypeBlack
	}

	param := proc_models.ProcGlobalModel{
		Id:         ctx.GetString(consts.ValidatorPrefixOrg + "id"),
		OrgId:      loginInfo.OrgId,
		Dir:        ctx.GetString(consts.ValidatorPrefixOrg + "dir"),
		Type:       Type,
		DataFrom:   dataFrom,
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
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
		UpdateId:   loginInfo.UserId,
		UpdateName: loginInfo.Username,
	}

	return proc_service.CreateGlobalCurdFactory().BatchAdd(&param, &list)
}

// ClearFromProcess 清空
func (u *ProcessGlobalController) ClearFromProcess(ctx *gin.Context) {
	if err := u.clear(ctx, enum_process.DataFromProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// clear
func (u *ProcessGlobalController) clear(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)

	return proc_service.CreateGlobalCurdFactory().DeleteByDataFrom(loginInfo.OrgId, dataFrom)
}

// ClearFromDll 清空
func (u *ProcessGlobalController) ClearFromDll(ctx *gin.Context) {
	if err := u.clear(ctx, enum_process.DataFromDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
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

// SetSwitchFromProcess 设置开关来源进程
func (u *ProcessGlobalController) SetSwitchFromProcess(ctx *gin.Context) {
	if err := u.setSwitch(ctx, enum_process.DataFromProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// setSwitch
func (u *ProcessGlobalController) setSwitch(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)
	openSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	id := ctx.GetString(consts.ValidatorPrefixOrg + "id")

	return proc_service.CreateGlobalCurdFactory().SetSwitch(loginInfo.OrgId, id, openSwitch, dataFrom)
}

// SetSwitchFromDll 设置开关来源模块
func (u *ProcessGlobalController) SetSwitchFromDll(ctx *gin.Context) {
	if err := u.setSwitch(ctx, enum_process.DataFromDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteFromProcess 删除来源进程
func (u *ProcessGlobalController) DeleteFromProcess(ctx *gin.Context) {
	if err := u.delete(ctx, enum_process.DataFromProcess); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// delete
func (u *ProcessGlobalController) delete(ctx *gin.Context, dataFrom int) error {
	loginInfo := u.getLogin(ctx)
	id := ctx.GetString(consts.ValidatorPrefixOrg + "id")

	return proc_service.CreateGlobalCurdFactory().DeleteById(loginInfo.OrgId, id, dataFrom)
}

// DeleteFromDll 删除来源模块
func (u *ProcessGlobalController) DeleteFromDll(ctx *gin.Context) {
	if err := u.delete(ctx, enum_process.DataFromDll); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
