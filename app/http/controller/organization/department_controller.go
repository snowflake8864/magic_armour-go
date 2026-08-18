package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/tree"
)

// DepartmentController 组织架构相关
type DepartmentController struct {
	ControllerBase
}

// Tree 树形列表
func (u *DepartmentController) Tree(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	conditions := department_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}
	if loginInfo.IsRoot == enums.IsNot {
		conditions.Ids = loginInfo.DeptWithChild
		if conditions.Ids == nil || len(conditions.Ids) == 0 {
			response.SuccessReturnData(ctx, []interface{}{})
			return
		}
	}
	models, err := department_models.CreateDepartmentFactory().GetListByConditions(conditions)
	if err != nil || models == nil {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}

	var nodes []*tree.Node
	for _, model := range models {
		nodes = append(nodes, &tree.Node{
			ID:       model.Id,
			ParentId: model.ParentId,
			Data: map[string]interface{}{
				"name":        model.Name,
				"ip":          model.Ip,
				"full_code":   model.FullCode,
				"level":       model.Level,
				"create_name": model.CreateName,
				"created_at":  model.CreatedAt,
			},
			Child: []*tree.Node{},
		})
	}

	response.SuccessReturnData(ctx, tree.RealTree(nodes))
}

// Add 添加
func (u *DepartmentController) Add(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := department_models.DepartmentModel{
		OrgId:      loginInfo.OrgId,
		ParentId:   int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "parent_id")),
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := department_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Edit 编辑
func (u *DepartmentController) Edit(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	param := department_models.DepartmentModel{
		OrgId:      loginInfo.OrgId,
		Id:         int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		Name:       ctx.GetString(consts.ValidatorPrefixOrg + "name"),
		Ip:         ctx.GetString(consts.ValidatorPrefixOrg + "ip"),
		CreateId:   loginInfo.UserId,
		CreateName: loginInfo.Username,
	}

	if err := department_service.CreateCurdFactory().UpdateById(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// Delete 删除
func (u *DepartmentController) Delete(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id"))

	if !department_service.CreateCurdFactory().DelSelfWithChildById(loginInfo.OrgId, id) {
		response.FailSimplify(ctx, my_errors.CurdDeleteFailed)
		return
	}

	response.SuccessSimplify(ctx)
}
