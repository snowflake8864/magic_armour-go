package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

// TaskController  终端任务相关
type TaskController struct {
	ControllerBase
}

// TaskTypes  订单状态
func (u *TaskController) TaskTypes(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("type", enum_task.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_task.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, true))
}

// TaskStatus 任务状态
func (u *TaskController) TaskStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("status", enum_task.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("status", enum_task.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, true))
}

// PageForHistory 历史：分页
func (u *TaskController) PageForHistory(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	Type := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "type"))
	condition := task_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
		Type:    &Type,
		Status:  int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "status")),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := task_models.CreateTaskFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.TaskTrans{Lang: cur_userinfo.GetLocales(ctx)}).TransPage(total, models))
}

// ClearForHistory 历史：清空
func (u *TaskController) ClearForHistory(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	if !task_service.CreateCurdFactory().DeleteForAgent(loginInfo.OrgId, agentId) {
		response.FailSimplify(ctx, my_errors.CurdClearFailed)
		return
	}
	response.SuccessSimplify(ctx)
}
