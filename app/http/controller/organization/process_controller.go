package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/http/transformer"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// SearchStatus 查询状态
func (u *ProcessController) SearchStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("proc_status", enum_process.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("all_status", enum_process.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// ChangeStatus 可改变状态
func (u *ProcessController) ChangeStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("change_status", enum_process.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("all_status", enum_process.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// PageForManage 管理页: 分页
func (u *ProcessController) PageForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForOwn")
}

// AgentProcStatForManage 管理页: 终端进程统计
func (u *ProcessController) AgentProcStatForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AgentProcStatForManage")
}

// SetSwitchForManage 管理页: 设置开关
func (u *ProcessController) SetSwitchForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForManage")
}

// SetModeForManage 管理页: 设置模式
func (u *ProcessController) SetModeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetModeForManage")
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *ProcessController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchModifyForManage")
}

// BatchSetWhiteForManage 管理页: 批量把终端的未设置进程变成白名单
func (u *ProcessController) BatchSetWhiteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchSetWhiteForManage")
}

// BatchRefreshForManage 管理页: 批量刷新(运行中进程)
func (u *ProcessController) BatchRefreshForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchRefreshForManage")
}

// BatchCollectForManage 管理页: 批量获取全部进程
func (u *ProcessController) BatchCollectForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "batchCollectForManage")
}

// BatchThreatenWhiteForManage 管理页: 批量把终端的进程推送到威胁情报白名单
func (u *ProcessController) BatchThreatenWhiteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchThreatenWhiteForManage")
}

// PageAllForAgent 单终端: 全部进程
func (u *ProcessController) PageAllForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageAllForAgent")
}

// PageRunningForAgent 单终端: 运行中进程
func (u *ProcessController) PageRunningForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageRunningForAgent")
}

// SearchWhileForAgent 单终端: 查询白名单
func (u *ProcessController) SearchWhileForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetStatusForAgent")
}

// SetStatusForAgent 单终端: 设置状态
func (u *ProcessController) SetStatusForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetStatusForAgent")
}

// BatchWhiteForAgent 单终端: 批量白名单
func (u *ProcessController) BatchWhiteForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhiteForAgent")
}

// CancelForAgent 单终端: 取消
func (u *ProcessController) CancelForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CancelForAgent")
}

// ClearAllForAgent 单终端: 清空全部
func (u *ProcessController) ClearAllForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearAllForAgent")
}

// ClearRunningForAgent 单终端: 清空运行中
func (u *ProcessController) ClearRunningForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearRunningForAgent")
}

// RefreshForAgent 单终端: 刷新(运行中进程)
func (u *ProcessController) RefreshForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RefreshForAgent")
}

// CollectForAgent 单终端: 获取全部进程
func (u *ProcessController) CollectForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CollectForAgent")
}

// RepeatStatForAgent 单终端: 重复进程统计
func (u *ProcessController) RepeatStatForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RepeatHashForAgent")
}

// RepeatPageForAgent 单终端: 重复进程分页
func (u *ProcessController) RepeatPageForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RepeatHashForAgent")
}

// RunningDllExtractForAgent 单终端: 提取运行中模块
func (u *ProcessController) RunningDllExtractForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RunningDllExtractForAgent")
}

// RunningDllShowForAgent 单终端: 展示运行中模块
func (u *ProcessController) RunningDllShowForAgent(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "RunningDllShowForAgent")
}

// PageForSearch 搜索: 分页
func (u *ProcessController) PageForSearch(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForSearch")
}
