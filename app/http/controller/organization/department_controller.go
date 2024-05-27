package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/utils/response"
)

// DepartmentController 组织架构相关
type DepartmentController struct {
	platform.ControllerBase
}

// Tree 树形列表
func (u *DepartmentController) Tree(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Tree")
}

// Add 添加
func (u *DepartmentController) Add(ctx *gin.Context) {
	parentId := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "parentId"))
	name := ctx.GetString(consts.ValidatorPrefixOrg + "name")
	ip := ctx.GetString(consts.ValidatorPrefixOrg + "ip")
	fmt.Println("============:", parentId, name, ip)
	response.SuccessReturnData(ctx, "add")
}

// Edit 编辑
func (u *DepartmentController) Edit(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "edit")
}

// Delete 删除
func (u *DepartmentController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "delete")
}
