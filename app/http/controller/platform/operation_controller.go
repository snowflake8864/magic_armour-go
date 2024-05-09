package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/operate_log_service"
	"goskeleton/app/utils/response"
)

// OperationController 首页相关
type OperationController struct {
	ControllerBase
}

// Modules 日志模块
func (u *OperationController) Modules(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfStr("module", enum_platform_operate.WebMapString)
	labels := enums.Get1stIsStrReturnMapStrOfStr("module", enum_platform_operate.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatStringList(keys, labels, true))
}

// Page  分页
func (u *OperationController) Page(ctx *gin.Context) {
	//页面上传
	conditions := operate_log_service.SearchConditions{}
	conditions.Module = ctx.GetString(consts.ValidatorPrefixPlatform + "module")
	conditions.NetIp = ctx.GetString(consts.ValidatorPrefixPlatform + "netIp")
	conditions.CreateName = ctx.GetString(consts.ValidatorPrefixPlatform + "createName")

	//分页参数和排序参数
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, "")

	//产品分页
	total, models := operate_log_service.CreateCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)

	response.SuccessReturnData(ctx, (&platform_trans.Operation{}).TransPage(total, models))
}
