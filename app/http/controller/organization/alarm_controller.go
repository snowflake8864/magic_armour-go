package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/response"
)

// AlarmController  告警相关
type AlarmController struct {
	ControllerBase
}

// HandleStatus 处理状态
func (u *AlarmController) HandleStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("handle_status", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("handle_status", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// Levels 告警等级
func (u *AlarmController) Levels(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("level", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("level", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// TypeProcess  进程告警类型
func (u *AlarmController) TypeProcess(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_process", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// TypeDll  模块告警类型
func (u *AlarmController) TypeDll(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_dll", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// TypeBlackmail  勒索告警类型
func (u *AlarmController) TypeBlackmail(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_blackmail", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// TypeTamper  防篡改告警类型
func (u *AlarmController) TypeTamper(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_tamper", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// TypeUsb  外设告警类型
func (u *AlarmController) TypeUsb(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_usb", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

func (u *AlarmController) TypeSelf(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type_self", enum_alarm.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_alarm.EnumMapInt)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// PageForManage 分页
func (u *AlarmController) PageForManage(ctx *gin.Context) {

	//查询条件
	conditions, err := u.getAgentSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_agent.OrderFieldDefault)

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, order)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, nil))
}

// ClearForManage 管理页: 清空
func (u *AlarmController) ClearForManage(ctx *gin.Context) {
	agentIds := []string{
		u.getAgentIdFromGin(ctx),
	}
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))
	var param = agent_service.ManualColumn{
		LoginSolveSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(u.getLogin(ctx).OrgId, agentIds, &param, "口令破解管理-开关")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// BatchClearForManage 管理页: 批量清空
func (u *AlarmController) BatchClearForManage(ctx *gin.Context) {
	orgId := u.getLogin(ctx).OrgId
	agentIds := u.getAgentIdsFromGin(ctx)
	Switch := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "switch"))

	var param = agent_service.ManualColumn{
		LoginSolveSwitch: &Switch,
	}
	errCode, err := agent_service.CreateCurdFactory().BatchUpdate(orgId, agentIds, &param, "口令破解管理-批量")
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
