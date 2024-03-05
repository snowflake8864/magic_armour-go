package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// AgentController 终端相关
type AgentController struct {
}

// GetConfig 获取配置信息
func (u *AgentController) GetConfig(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取配置信息")
}

// GetInfo 获取终端信息(调试用)
func (u *AgentController) GetInfo(ctx *gin.Context) {
	currentAgent := cur_userinfo.GetCurrentAgent(ctx)
	response.SuccessReturnData(ctx, currentAgent)
}

// Upgrade 升级
func (u *AgentController) Upgrade(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "升级")
}

// Uninstall 卸载回调
func (u *AgentController) Uninstall(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "卸载回调")
}

// GetTrustDir 获取信任目录
func (u *AgentController) GetTrustDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取信任目录")
}

// UploadErrorFile 上传转储(错误) 文件
func (u *AgentController) UploadErrorFile(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传转储(错误) 文件")
}

// MonitorHardware 监控硬件资源
func (u *AgentController) MonitorHardware(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "监控硬件资源")
}
