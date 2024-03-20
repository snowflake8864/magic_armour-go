package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// SearchStatus 查询状态
func (u *ProcessController) SearchStatus(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SearchStatus")
}

// ChangeStatus 可改变状态
func (u *ProcessController) ChangeStatus(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ChangeStatus")
}

// CountOnline 统计在线和离线数量
func (u *ProcessController) CountOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CountOnline")
}

// TypeList 类型列表
func (u *ProcessController) TypeList(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "TypeList")
}

// PageForManage 管理: 分页
func (u *ProcessController) PageForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForManage")
}

// DeleteForManage 管理: 删除
func (u *ProcessController) DeleteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForManage")
}

// UninstallForManage 管理: 卸载
func (u *ProcessController) UninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UninstallForManage")
}

// UpgradeForManage 管理: 升级
func (u *ProcessController) UpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UpgradeForManage")
}

// BatchDeleteForManage 管理: 批量删除
func (u *ProcessController) BatchDeleteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchDeleteForManage")
}

// BatchUninstallForManage 管理: 批量卸载
func (u *ProcessController) BatchUninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUninstallForManage")
}

// BatchUpgradeForManage 管理: 批量升级
func (u *ProcessController) BatchUpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUpgradeForManage")
}

// BatchSetConfigForManage 管理: 批量设置配置
func (u *ProcessController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchSetConfigForManage")
}

// GetForBasic 基本信息：获取
func (u *ProcessController) GetForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForBasic")
}

// SaveForBasic 基本信息：保存
func (u *ProcessController) SaveForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForBasic")
}

// GetForConfig 配置信息：获取
func (u *ProcessController) GetForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForConfig")
}

// SaveForConfig 配置信息：保存
func (u *ProcessController) SaveForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForConfig")
}

// ListForMonitor 监控：列表
func (u *ProcessController) ListForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ListForMonitor")
}

// InfoForMonitor 监控：信息
func (u *ProcessController) InfoForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "InfoForMonitor")
}

// PageForError 错误文件：分页
func (u *ProcessController) PageForError(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForError")
}
