package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
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

// AgentStatForManage 管理页: 终端告警统计
func (u *AlarmController) AgentStatForManage(ctx *gin.Context) {
	conditions := alarm_models.SearchConditions{
		OrgId:    u.getLogin(ctx).OrgId,
		AgentIds: u.getAgentIdsFromGin(ctx),
	}
	module := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "module"))

	var calcList []*alarm_models.AlarmCalc
	switch module {
	case enum_alarm.ModuleProcess:
		fact := alarm_models.CreateProcessFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleDll:
		fact := alarm_models.CreateDllFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleTamper:
		fact := alarm_models.CreateTamperFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleBlackmail:
		fact := alarm_models.CreateBlackmailFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModulePort:
		fact := alarm_models.CreatePortFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModulePlug:
		fact := alarm_models.CreatePlugFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleUsb:
		fact := alarm_models.CreateUsbFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleOutreach:
		fact := alarm_models.CreateOutreachFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleLogin:
		fact := alarm_models.CreateLoginFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	case enum_alarm.ModuleSelf:
		fact := alarm_models.CreateSelfFactory()
		calcList = fact.CalcByConditionsGroupAgent(conditions, fact.ColumnCalc(false, true, true, false))
		break
	}
	response.SuccessReturnData(ctx, calcList)
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

// PageForAlarm 告警: 分页
func (u *AlarmController) PageForAlarm(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)
	condition := alarm_models.SearchConditions{
		OrgId: loginInfo.OrgId,
	}

	//用户管理的终端
	err, userAgentIds := u.getCurrUserAgentIds(loginInfo)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
	}
	//查询终端比较用户管理的终端
	agentId := u.getAgentIdFromGin(ctx)
	if !string_helper.IsEmpty(agentId) {
		if slice_helper.StringIn(agentId, userAgentIds) {
			condition.AgentId = agentId
		} else {
			response.FailSimplify(ctx, "数据为空")
		}
	} else {
		condition.AgentIds = userAgentIds
	}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	var agentMap = make(map[string]*agent_models.AgentModel)
	var agentIds []string
	switch module {
	case enum_alarm.ModuleProcess:
		fact := alarm_models.CreateProcessFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransProcessPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleDll:
		fact := alarm_models.CreateDllFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransDllPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleTamper:
		fact := alarm_models.CreateTamperFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransTamperPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleBlackmail:
		fact := alarm_models.CreateBlackmailFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransBlackmailPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModulePort:
		fact := alarm_models.CreatePortFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransPortPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModulePlug:
		fact := alarm_models.CreatePlugFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransPlugPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleUsb:
		fact := alarm_models.CreateUsbFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransUsbPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleOutreach:
		fact := alarm_models.CreateOutreachFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransOutreachPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleLogin:
		fact := alarm_models.CreateLoginFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransLoginPage(total, models, agentMap, nil))
		break
	case enum_alarm.ModuleSelf:
		fact := alarm_models.CreateSelfFactory()
		total, models := fact.Page(condition, offset, pageSize, order)

		if models != nil || len(models) > 0 {
			for _, model := range models {
				agentIds = append(agentIds, model.AgentId)
			}
			agentMap = u.getAgentMap(loginInfo.OrgId, agentIds)
		}
		response.SuccessReturnData(ctx, (&org_trans.AlarmTrans{}).TransSelfPage(total, models, agentMap, nil))
		break
	default:
		response.FailSimplify(ctx, fmt.Sprintf("module(%d) not exists", module))
	}
}

// HandleForAlarm 告警: 处理
func (u *AlarmController) HandleForAlarm(ctx *gin.Context, module int) {

}

// HandleAllForAlarm 告警: 处理全部
func (u *AlarmController) HandleAllForAlarm(ctx *gin.Context, module int) {

}

// ClearForAlarm 告警: 清空
func (u *AlarmController) ClearForAlarm(ctx *gin.Context, module int) {

}
