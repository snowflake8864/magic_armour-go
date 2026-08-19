package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/organization/org_operate_models"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// OperationController  操作日志相关
type OperationController struct {
	ControllerBase
}

// Modules  功能模块
func (u *OperationController) Modules(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfStr("module", enum_org_operate.WebMapString)
	labels := enums.Get1stIsStrReturnMapStrOfStr("module", enum_org_operate.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatStringList(keys, labels, true))
}

// PageForLog 日志：分页
func (u *OperationController) PageForLog(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := org_operate_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		Module: ctx.GetString(consts.ValidatorPrefixOrg + "module"),
		NetIp:  ctx.GetString(consts.ValidatorPrefixOrg + "netIp"),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := org_operate_models.CreateLogFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Operation{Lang: cur_userinfo.GetLocales(ctx)}).TransPage(total, models))
}
