package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// AgentController 角色相关
type AgentController struct {
	ControllerBase
}

// OnlineStatus 在线状态
func (u *AgentController) OnlineStatus(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "OnlineStatus")
}

// SafeStatus 安全状态
func (u *AgentController) SafeStatus(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SafeStatus")
}

// CountOnline 统计在线和离线数量
func (u *AgentController) CountOnline(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "CountOnline")
}

// TypeList 类型列表
func (u *AgentController) TypeList(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "TypeList")
}

// PageForManage 管理: 分页
func (u *AgentController) PageForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForManage")
}

// DeleteForManage 管理: 删除
func (u *AgentController) DeleteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForManage")
}

// UninstallForManage 管理: 卸载
func (u *AgentController) UninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UninstallForManage")
}

// UpgradeForManage 管理: 升级
func (u *AgentController) UpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "UpgradeForManage")
}

// BatchDeleteForManage 管理: 批量删除
func (u *AgentController) BatchDeleteForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchDeleteForManage")
}

// BatchUninstallForManage 管理: 批量卸载
func (u *AgentController) BatchUninstallForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUninstallForManage")
}

// BatchUpgradeForManage 管理: 批量升级
func (u *AgentController) BatchUpgradeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchUpgradeForManage")
}

// BatchSetConfigForManage 管理: 批量设置配置
func (u *AgentController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchSetConfigForManage")
}

// GetForBasic 基本信息：获取
func (u *AgentController) GetForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForBasic")
}

// SaveForBasic 基本信息：保存
func (u *AgentController) SaveForBasic(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForBasic")
}

// GetForConfig 配置信息：获取
func (u *AgentController) GetForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetForConfig")
}

// SaveForConfig 配置信息：保存
func (u *AgentController) SaveForConfig(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SaveForConfig")
}

// ListForMonitor 监控：列表
func (u *AgentController) ListForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ListForMonitor")
}

// InfoForMonitor 监控：信息
func (u *AgentController) InfoForMonitor(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "InfoForMonitor")
}

// PageForError 错误文件：分页
func (u *AgentController) PageForError(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForError")
}
