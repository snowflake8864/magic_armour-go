package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/response"
)

// FileDirController 文件目录相关
type FileDirController struct {
	ControllerBase
}

// GetTypes 类型
func (u *FileDirController) GetTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("search_type", enum_tamper.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("search_type", enum_tamper.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// Page 分页
func (u *FileDirController) Page(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//查询条件
	condition := trust_dir_models.SearchConditions{
		OrgId: loginInfo.OrgId,
		Name:  ctx.GetString(consts.ValidatorPrefixOrg + "name"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := trust_dir_models.CreateGlobalFactory().Page(condition, offset, pageSize, order)

	//获取组织架构
	departmentMaps := map[int64]*department_models.DepartmentModel{}
	if models != nil && len(models) > 0 {
		var deptIds []int64
		for _, model := range models {
			if model.DeptId > 0 {
				deptIds = append(deptIds, model.DeptId)
			}
		}
		departmentMaps = u.getDepartmentMap(loginInfo.OrgId, slice_helper.Int64RemoveRepeatByMap(deptIds))
	}

	response.SuccessReturnData(ctx, (&org_trans.TrustDirTrans{}).TransGlobalPage(total, models, departmentMaps))
}

// HasUpload 判断是否上传
func (u *FileDirController) HasUpload(ctx *gin.Context) {
	response.SuccessSimplify(ctx)
}

// Refresh 刷新
func (u *FileDirController) Refresh(ctx *gin.Context) {
	response.SuccessSimplify(ctx)
}
