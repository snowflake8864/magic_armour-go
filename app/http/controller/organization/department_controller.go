package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/response"
)

// DepartmentController 组织架构相关
type DepartmentController struct {
	ControllerBase
}

// Tree 树形列表
func (u *DepartmentController) Tree(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Tree")
}

// Add 添加
func (u *DepartmentController) Add(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := department_models.DepartmentModel{
		OrgId:      loginInfo.OrgId,
		ParentId:   int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "parentId")),
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := department_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

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
