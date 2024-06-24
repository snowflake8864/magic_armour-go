package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/platform/org_menu_models"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/tree"
)

// AccountController 首页相关
type AccountController struct {
	ControllerBase
}

func (u *AccountController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// GetMenuTree 获取菜单菜单树
func (u *AccountController) GetMenuTree(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	conditions := org_menu_models.SearchConditions{
		UsageVerList: []int{
			0,
			cur_app.GetAppUsageVersion(),
		},
	}
	//不是管理员， 只能看非管理员菜单
	if loginInfo.IsRoot == enums.IsNot {
		IsNot := enums.IsNot
		conditions.IsRoot = &IsNot
	}
	models := org_menu_models.CreateOrgMenuFactory().GetListByConditions(conditions)
	if models == nil {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}
	fmt.Println(json_helper.Encode(models))

	var nodes []*tree.Node
	for _, model := range models {
		nodes = append(nodes, &tree.Node{
			ID:       model.Id,
			ParentId: model.ParentId,
			Data: map[string]interface{}{
				"name":       model.Name,
				"type":       model.Type,
				"code":       model.Code,
				"sort":       model.Sort,
				"is_show":    model.IsShow,
				"is_root":    model.IsRoot,
				"usage_ver":  model.UsageVer,
				"remark":     model.Remark,
				"api":        model.Api,
				"level":      model.Level,
				"created_at": model.CreatedAt,
			},
			Child: []*tree.Node{},
		})
	}
	fmt.Println(json_helper.Encode(nodes))
	response.SuccessReturnData(ctx, tree.RealTree(nodes))
}

// ModifyPw 修改密码
func (u *AccountController) ModifyPw(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	newPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "newPw"))
	if string_helper.IsEmpty(oldPw) || string_helper.IsEmpty(newPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyPwById(loginInfo.UserId, oldPw, newPw); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyPhone 修改手机
func (u *AccountController) ModifyPhone(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	phone := ctx.GetString(consts.ValidatorPrefixOrg + "phone")
	if string_helper.IsEmpty(oldPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyPhoneById(loginInfo.UserId, oldPw, phone); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyEmail 修改邮箱
func (u *AccountController) ModifyEmail(ctx *gin.Context) {
	oldPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "oldPw"))
	email := ctx.GetString(consts.ValidatorPrefixOrg + "email")
	if string_helper.IsEmpty(oldPw) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	loginInfo := u.getLogin(ctx)
	if errCode, err := user_service.CreateCurdFactory().ModifyEmailById(loginInfo.UserId, oldPw, email); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
