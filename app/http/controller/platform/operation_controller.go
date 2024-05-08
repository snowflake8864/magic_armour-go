package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/http/transformer"
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

	response.SuccessSimplify(ctx)
}
