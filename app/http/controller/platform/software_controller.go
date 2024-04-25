package platform

import (
	"github.com/gin-gonic/gin"
)

// SoftwareController 终端相关
type SoftwareController struct {
	ControllerBase
}

// GetOsTypes  获取系统类型
func (u *SoftwareController) GetOsTypes(ctx *gin.Context) {

}

// GetAgentConfig  获取终端配置
func (u *SoftwareController) GetAgentConfig(ctx *gin.Context) {

}

// SetAgentConfig  保存终端配置
func (u *SoftwareController) SetAgentConfig(ctx *gin.Context) {

}

// PageForManage 管理: 分页
func (u *SoftwareController) PageForManage(ctx *gin.Context) {

}

// AddForManage 管理:  添加
func (u *SoftwareController) AddForManage(ctx *gin.Context) {

}

// ClearForManage 管理: 清空
func (u *SoftwareController) ClearForManage(ctx *gin.Context) {

}

// DownloadForManage 管理: 下载
func (u *SoftwareController) DownloadForManage(ctx *gin.Context) {

}

// PublishForManage 管理: 发布
func (u *SoftwareController) PublishForManage(ctx *gin.Context) {

}

// CancelForManage 管理: 取消
func (u *SoftwareController) CancelForManage(ctx *gin.Context) {

}

// DeleteForManage 管理: 删除
func (u *SoftwareController) DeleteForManage(ctx *gin.Context) {

}
