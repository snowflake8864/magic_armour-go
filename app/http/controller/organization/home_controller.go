package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/response"
)

// HomeController 首页相关
type HomeController struct {
	ControllerBase
}

// AgentStat 终端：实时统计
func (u *HomeController) AgentStat(ctx *gin.Context) {

	loginInfo := u.getLogin(ctx)
	conditions := agent_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.IsOnline = enums.AllDefaultInt
	conditions.IsSafe = enums.AllDefaultInt
	err, agentIds := u.getCurrUserAgentIds(u.getLogin(ctx))

	if err != nil {
		//不是管理员 & 无终端
		response.SuccessReturnData(ctx, agent_models.AgentCalc{})
		return
	} else {
		//不是管理员有终端
		if agentIds != nil {
			conditions.AgentIds = agentIds
		}
	}
	//统计
	agentFact := agent_models.CreateAgentFactory()
	result := agentFact.CalcByConditions(conditions, agentFact.ColumnCalcAgents(false, true, true, true, true))

	response.SuccessReturnData(ctx, agentFact.FillCalcData(result))
}

// AlarmAgentTops 告警：终端排行
func (u *HomeController) AlarmAgentTops(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	alarmConditions := alarm_stat_models.SearchConditions{}
	alarmConditions.OrgId = loginInfo.OrgId
	alarmConditions.StatDay = date_helper.SpecialDatetime(variable.DayFormat) //今天
	err, agentIds := u.getCurrUserAgentIds(u.getLogin(ctx))
	if err != nil {
		//不是管理员 & 无终端
		response.SuccessReturnData(ctx, []interface{}{})
		return
	} else {
		//不是管理员有终端
		if agentIds != nil {
			alarmConditions.AgentIds = agentIds
		}
	}
	//今天新增最多的8个终端
	alarmFact := alarm_stat_models.CreateAlarmStatDailyFactory()
	alarmModels := alarmFact.GetItemsByConditions(alarmConditions, 0, 8, "add_total desc")
	if alarmModels == nil || len(alarmModels) == 0 {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}

	//扩展信息
	var extendInfo = make(map[string]interface{})
	var topIds []string
	for _, model := range alarmModels {
		topIds = append(topIds, model.AgentId)
		extendInfo[model.AgentId] = map[string]interface{}{"number_total": model.AddTotal}
	}

	//终端
	agentFact := agent_models.CreateAgentFactory()
	agentMaps := agentFact.GetMapsByOrgAndAgents(loginInfo.OrgId, topIds, agentFact.ColumnBasic(true, false))
	if agentMaps == nil || len(agentMaps) == 0 {
		response.SuccessReturnData(ctx, []interface{}{})
		return
	}

	//按照排序返回
	var list []*agent_models.AgentModel
	for _, agentId := range topIds {
		if model, isOK := agentMaps[agentId]; isOK {
			list = append(list, model)
		}
	}

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransList(list, nil, extendInfo))
}

// AlarmTodayStat 告警： 今天统计
func (u *HomeController) AlarmTodayStat(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	conditions := alarm_stat_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.StatDay = date_helper.SpecialDatetime(variable.DayFormat)
	err, agentIds := u.getCurrUserAgentIds(u.getLogin(ctx))
	if err != nil {
		//不是管理员 & 无终端
		response.SuccessReturnData(ctx, alarm_stat_models.AlarmStatDailyModel{})
		return
	} else {
		//不是管理员有终端
		if agentIds != nil {
			conditions.AgentIds = agentIds
		}
	}
	//统计
	alarmFact := alarm_stat_models.CreateAlarmStatDailyFactory()
	result := alarmFact.CalcByConditions(conditions, alarmFact.StatColumnOfAllType(true, false, true, false, false))

	response.SuccessReturnData(ctx, result)

}

// AlarmDailyStat 告警： 每天统计
func (u *HomeController) AlarmDailyStat(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	conditions := alarm_stat_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.StatMonth = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "month"))
	err, agentIds := u.getCurrUserAgentIds(u.getLogin(ctx))
	if err != nil {
		//不是管理员 & 无终端
		response.SuccessReturnData(ctx, alarm_stat_models.AlarmStatDailyModel{})
		return
	} else {
		//不是管理员有终端
		if agentIds != nil {
			conditions.AgentIds = agentIds
		}
	}
	//统计
	alarmFact := alarm_stat_models.CreateAlarmStatDailyFactory()
	models := alarmFact.CalcByConditionsGroupDay(conditions, alarmFact.StatColumn(alarm_stat_models.PrefixTotal, true, false, true, false, true))
	var result []map[string]interface{}
	if models != nil && len(models) > 0 {
		for _, model := range models {
			result = append(result, map[string]interface{}{
				"date":      date_helper.DaysInMonthByTimestamp(date_helper.DatetimeToTimestamp(model.StatDay, variable.DayFormat)),
				"total":     model.Total,
				"add_total": model.AddTotal,
			})
		}
	}
	response.SuccessReturnData(ctx, result)
}

// GrantUsedStat 授权： 使用统计
func (u *HomeController) GrantUsedStat(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	conditions := order_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	conditions.OrderStatus = enum_order.OrderStatusCompleted

	orderFact := order_models.CreateOrderFactory()
	models := orderFact.GetItemsByConditions(conditions, 0, 2, "")

	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransList(models))
}
