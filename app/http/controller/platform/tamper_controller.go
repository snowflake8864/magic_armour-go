package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
	ControllerBase
}

// FilePower 文件权限
func (u *TamperController) FilePower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// FolderPower 文件夹权限
func (u *TamperController) FolderPower(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FolderPower")
}

// Page 分页
func (u *TamperController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// SetSwitchForManage 管理页: 设置开关
func (u *TamperController) SetSwitchForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForManage")
}

// SetModeForManage 管理页: 设置模式
func (u *TamperController) SetModeForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetModeForManage")
}

// BatchSetConfigForManage 管理页: 批量设置配置(开关和模式)
func (u *TamperController) BatchSetConfigForManage(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchSetConfigForManage")
}

// PageForPolicy 终端策略: 分页
func (u *TamperController) PageForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForPolicy")
}

// SetSwitchForPolicy 终端策略: 设置开关
func (u *TamperController) SetSwitchForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForPolicy")
}

// AddFolderForPolicy 终端策略: 添加文件夹
func (u *TamperController) AddFolderForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFolderForPolicy")
}

// AddFilForPolicy 终端策略: 添加文件
func (u *TamperController) AddFilForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFilForPolicy")
}

// EditFolderForPolicy 终端策略: 编辑文件夹
func (u *TamperController) EditFolderForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFolderForPolicy")
}

// EditFilForPolicy 终端策略: 编辑文件
func (u *TamperController) EditFilForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "EditFilForPolicy")
}

// DeleteForPolicy 终端策略: 删除
func (u *TamperController) DeleteForPolicy(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForPolicy")
}

// PageForGlobal 全局策略: 分页
func (u *TamperController) PageForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageForGlobal")
}

// SetSwitchForGlobal 全局策略: 设置开关
func (u *TamperController) SetSwitchForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "SetSwitchForGlobal")
}

// AddFolderForGlobal 全局策略: 添加文件夹
func (u *TamperController) AddFolderForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFolderForGlobal")
}

// AddFilForGlobal 全局策略: 添加文件
func (u *TamperController) AddFilForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFilForGlobal")
}

// EditFolderForGlobal 全局策略: 编辑文件夹
func (u *TamperController) EditFolderForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFolderForGlobal")
}

// EditFilForGlobal 全局策略: 编辑文件
func (u *TamperController) EditFilForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "EditFilForGlobal")
}

// DeleteForGlobal 全局策略: 删除
func (u *TamperController) DeleteForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "DeleteForGlobal")
}
