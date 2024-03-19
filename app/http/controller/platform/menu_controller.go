package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// MenuController 角色相关
type MenuController struct {
	ControllerBase
}

// Tree 树形列表
func (u *MenuController) Tree(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Tree")
}

// Add 添加
func (u *MenuController) Add(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "add")
}

// Edit 编辑
func (u *MenuController) Edit(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "edit")
}

// Sort 排序
func (u *MenuController) Sort(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "sort")
}

// Delete 删除
func (u *MenuController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "delete")
}
