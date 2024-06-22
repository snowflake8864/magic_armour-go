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
	"goskeleton/app/service/agent/alarm_service"
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
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	_ = alarm_service.CreateProcessCurdFactory().ClearByAgent(&alarm_models.ProcessModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateDllCurdFactory().ClearByAgent(&alarm_models.DllModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateTamperCurdFactory().ClearByAgent(&alarm_models.TamperModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateBlackmailCurdFactory().ClearByAgent(&alarm_models.BlackmailModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreatePortCurdFactory().ClearByAgent(&alarm_models.PortModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreatePlugCurdFactory().ClearByAgent(&alarm_models.PlugModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateUsbCurdFactory().ClearByAgent(&alarm_models.UsbModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateOutreachCurdFactory().ClearByAgent(&alarm_models.OutreachModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateLoginCurdFactory().ClearByAgent(&alarm_models.LoginModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	_ = alarm_service.CreateSelfCurdFactory().ClearByAgent(&alarm_models.SelfModel{
		OrgId:   loginInfo.OrgId,
		AgentId: agentId,
	})

	response.SuccessSimplify(ctx)
}

// BatchClearForManage 管理页: 批量清空
func (u *AlarmController) BatchClearForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentIds := u.getAgentIdsFromGin(ctx)

	_ = alarm_service.CreateProcessCurdFactory().ClearByAgents(&alarm_models.ProcessModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateDllCurdFactory().ClearByAgents(&alarm_models.DllModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateTamperCurdFactory().ClearByAgents(&alarm_models.TamperModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateBlackmailCurdFactory().ClearByAgents(&alarm_models.BlackmailModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreatePortCurdFactory().ClearByAgents(&alarm_models.PortModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreatePlugCurdFactory().ClearByAgents(&alarm_models.PlugModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateUsbCurdFactory().ClearByAgents(&alarm_models.UsbModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateOutreachCurdFactory().ClearByAgents(&alarm_models.OutreachModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateLoginCurdFactory().ClearByAgents(&alarm_models.LoginModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	_ = alarm_service.CreateSelfCurdFactory().ClearByAgents(&alarm_models.SelfModel{
		OrgId: loginInfo.OrgId,
	}, agentIds)

	response.SuccessSimplify(ctx)
}

// PageForAlarm 告警: 分页
func (u *AlarmController) PageForAlarm(ctx *gin.Context, module int) {
	loginInfo := u.getLogin(ctx)
	condition := alarm_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		CreateAtGte:  ctx.GetString(consts.ValidatorPrefixOrg + "startTime"),
		CreateAtLte:  ctx.GetString(consts.ValidatorPrefixOrg + "endTime"),
		HandleStatus: int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "handleStatus")),
	}

	//用户管理的终端
	err, userAgentIds := u.getCurrUserAgentIds(loginInfo)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	//查询终端比较用户管理的终端
	agentId := u.getAgentIdFromGin(ctx)
	if !string_helper.IsEmpty(agentId) && len(userAgentIds) > 0 {
		if slice_helper.StringIn(agentId, userAgentIds) {
			condition.AgentId = agentId
		} else {
			response.FailSimplify(ctx, "数据为空")
			return
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
		condition.ProcHash = ctx.GetString(consts.ValidatorPrefixOrg + "hash")
		condition.Type = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "alarmType"))
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
	loginInfo := u.getLogin(ctx)

	var err error
	switch module {
	case enum_alarm.ModuleProcess:
		err = alarm_service.CreateProcessCurdFactory().Handle(&alarm_models.ProcessModel{
			OrgId:    loginInfo.OrgId,
			ProcHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		})
		break
	case enum_alarm.ModuleDll:
		err = alarm_service.CreateDllCurdFactory().Handle(&alarm_models.DllModel{
			OrgId:   loginInfo.OrgId,
			DllHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		})
		break
	case enum_alarm.ModuleTamper:
		err = alarm_service.CreateTamperCurdFactory().Handle(&alarm_models.TamperModel{
			OrgId: loginInfo.OrgId,
			Id:    int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		})
		break
	case enum_alarm.ModuleBlackmail:
		err = alarm_service.CreateBlackmailCurdFactory().Handle(&alarm_models.BlackmailModel{
			OrgId: loginInfo.OrgId,
			Id:    int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		})
		break
	case enum_alarm.ModulePort:
		err = alarm_service.CreatePortCurdFactory().Handle(&alarm_models.PortModel{
			OrgId:    loginInfo.OrgId,
			AttackIp: ctx.GetString(consts.ValidatorPrefixOrg + "attackIp"),
		})
		break
	case enum_alarm.ModulePlug:
		err = alarm_service.CreatePlugCurdFactory().Handle(&alarm_models.PlugModel{
			OrgId: loginInfo.OrgId,
			Ip:    ctx.GetString(consts.ValidatorPrefixOrg + "plugIp"),
		})
		break
	case enum_alarm.ModuleUsb:
		err = alarm_service.CreateUsbCurdFactory().Handle(&alarm_models.UsbModel{
			OrgId:   loginInfo.OrgId,
			UsbHash: ctx.GetString(consts.ValidatorPrefixOrg + "hash"),
		})
		break
	case enum_alarm.ModuleOutreach:
		err = alarm_service.CreateOutreachCurdFactory().Handle(&alarm_models.OutreachModel{
			OrgId:    loginInfo.OrgId,
			TargetIp: ctx.GetString(consts.ValidatorPrefixOrg + "targetIp"),
		})
		break
	case enum_alarm.ModuleLogin:
		err = alarm_service.CreateLoginCurdFactory().Handle(&alarm_models.LoginModel{
			OrgId:    loginInfo.OrgId,
			SourceIp: ctx.GetString(consts.ValidatorPrefixOrg + "sourceIp"),
		})
		break
	case enum_alarm.ModuleSelf:
		err = alarm_service.CreateSelfCurdFactory().Handle(&alarm_models.SelfModel{
			OrgId: loginInfo.OrgId,
			Id:    int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "id")),
		})
		break
	}

	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// HandleAllForAlarm 告警: 处理全部
func (u *AlarmController) HandleAllForAlarm(ctx *gin.Context, module int) {

}

// ClearForAlarm 告警: 清空
func (u *AlarmController) ClearForAlarm(ctx *gin.Context, module int) {

}
