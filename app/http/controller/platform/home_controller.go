package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/model/organization/org_stat_models"
	"goskeleton/app/service/agent/alarm_service"
	"goskeleton/app/service/organization/org_log_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/response"
)

// HomeController 首页相关
type HomeController struct {
	ControllerBase
}

// calcRatio 计算比例
func (u *HomeController) calcRatio(divisor, dividend int) (expect, percent string, ratio float64) {
	if divisor > dividend {
		// divisor=3; dividend=1
		expect = "up" //上升
		if dividend == 0 {
			ratio = 0
		} else {
			ratio = float_helper.Div(divisor-dividend, dividend, 2)
		}

	} else if divisor < dividend {
		// divisor=1; dividend=3
		expect = "down" //下降
		if divisor == 0 {
			ratio = float_helper.FloatScale(float_helper.Float64(dividend), 2)
		} else {
			ratio = 1 - float_helper.Div(divisor, dividend, 2)
		}
	} else {
		expect = "same" //不变
		ratio = 0
	}
	percent = float_helper.Float64ToStr(ratio)
	if percent == "0" {
		percent = "0.00"
	}
	percent += "%"
	return expect, percent, ratio
}

// AlarmTodayStat 告警: 今天统计
func (u *HomeController) AlarmTodayStat(ctx *gin.Context) {
	//今天
	today := date_helper.SpecialDatetime(variable.DayFormat)
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	StatDailyFact := alarm_service.CreateStatDailyFactory()
	column := alarm_stat_models.CreateStatDailyFactory().StatColumn(alarm_stat_models.PrefixTotal, false,
		false, true, false, false)
	statToday := StatDailyFact.SumByDay(today, column)
	statYesterday := StatDailyFact.SumByDay(yesterday, column)

	expect, percent, ratio := u.calcRatio(statToday.AddTotal, statYesterday.AddTotal)

	response.SuccessReturnData(ctx, gin.H{
		"total_today":     statToday.AddTotal,
		"total_yesterday": statYesterday.AddTotal,
		"expect":          expect,
		"percent":         percent,
		"ratio":           ratio,
	})
}

// AgentTodayStat 终端: 今天统计
func (u *HomeController) AgentTodayStat(ctx *gin.Context) {
	//今天
	today := date_helper.SpecialDatetime(variable.DayFormat)
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	StatDailyFact := org_stat_models.CreateDailyFactory()
	column := StatDailyFact.ColumnOfAgent(false, true, false, false, false, false,
		false, false, false, false)
	conditions := org_stat_models.SearchConditions{}
	conditions.StatDay = today
	statToday := StatDailyFact.CalcByConditions(conditions, column)

	conditions.StatDay = yesterday
	statYesterday := StatDailyFact.CalcByConditions(conditions, column)

	expect, percent, ratio := u.calcRatio(statToday.AgentOnline, statYesterday.AgentOnline)
	response.SuccessReturnData(ctx, gin.H{
		"total_today":     statToday.AgentOnline,
		"total_yesterday": statYesterday.AgentOnline,
		"expect":          expect,
		"percent":         percent,
		"ratio":           ratio,
	})
}

// AgentDailyOnline 终端: 每天在线
func (u *HomeController) AgentDailyOnline(ctx *gin.Context) {
	conditions := org_stat_models.SearchConditions{
		StatMonth: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "month")),
	}

	StatDailyFact := org_stat_models.CreateDailyFactory()
	column := StatDailyFact.ColumnOfAgent(false, false, false, false, true, false,
		false, true, false, true)
	models := StatDailyFact.CalcByConditionsGroupDay(conditions, column)

	var result []map[string]interface{}
	if models != nil && len(models) > 0 {
		for _, model := range models {
			result = append(result, map[string]interface{}{
				"date":         date_helper.DaysInMonthByTimestamp(date_helper.DatetimeToTimestamp(model.StatDay, variable.DayFormat)),
				"win_online":   model.AgentWinOnline,
				"linux_online": model.AgentLinuxOnline,
			})
		}
	}
	response.SuccessReturnData(ctx, result)
}

// UserTodayStat 用户: 今天统计
func (u *HomeController) UserTodayStat(ctx *gin.Context) {
	//今天
	today := date_helper.SpecialDatetime(variable.DayFormat)
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	StatDailyFact := org_stat_models.CreateDailyFactory()
	column := StatDailyFact.ColumnOfUser(false, true, true, false)
	conditions := org_stat_models.SearchConditions{}
	conditions.StatDay = today
	statToday := StatDailyFact.CalcByConditions(conditions, column)

	conditions.StatDay = yesterday
	statYesterday := StatDailyFact.CalcByConditions(conditions, column)

	//新增
	expect, percent, ratio := u.calcRatio(statToday.UserAdd, statYesterday.UserAdd)
	mapAdd := map[string]interface{}{
		"total_today":     statToday.UserAdd,
		"total_yesterday": statYesterday.UserAdd,
		"expect":          expect,
		"percent":         percent,
		"ratio":           ratio,
	}

	//在线
	expect, percent, ratio = u.calcRatio(statToday.UserOnline, statYesterday.UserOnline)
	mapOnline := map[string]interface{}{
		"total_today":     statToday.UserOnline,
		"total_yesterday": statYesterday.UserOnline,
		"expect":          expect,
		"percent":         percent,
		"ratio":           ratio,
	}

	response.SuccessReturnData(ctx, gin.H{
		"add":    mapAdd,
		"online": mapOnline,
	})
}

// UserDailyOnline 用户: 每天在线
func (u *HomeController) UserDailyOnline(ctx *gin.Context) {
	conditions := org_stat_models.SearchConditions{
		StatMonth: int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "month")),
	}

	StatDailyFact := org_stat_models.CreateDailyFactory()
	column := StatDailyFact.ColumnOfUser(false, true, false, true)
	models := StatDailyFact.CalcByConditionsGroupDay(conditions, column)

	var result []map[string]interface{}
	if models != nil && len(models) > 0 {
		for _, model := range models {
			result = append(result, map[string]interface{}{
				"date":   date_helper.DaysInMonthByTimestamp(date_helper.DatetimeToTimestamp(model.StatDay, variable.DayFormat)),
				"online": model.UserOnline,
			})
		}
	}
	response.SuccessReturnData(ctx, result)
}

// UserNewestLogs 用户: 最新日志
func (u *HomeController) UserNewestLogs(ctx *gin.Context) {
	models := org_log_service.CreateOrgLogFactory().NewestList(10)

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransLogList(models))
}

// GrantUsedStat 授权: 使用统计
func (u *HomeController) GrantUsedStat(ctx *gin.Context) {
	searchConditions := order_models.SearchConditions{}
	searchConditions.OrderStatus = enum_order.OrderStatusCompleted

	CurdFact := order_models.CreateOrderFactory()
	column := CurdFact.ColumnCalcOfUsage()
	model := CurdFact.CalcByConditions(searchConditions, column)
	var result []map[string]interface{}
	result = append(result,
		map[string]interface{}{
			"name":  "已使用",
			"value": model.UsageNumber,
		}, map[string]interface{}{
			"name":  "未使用",
			"value": model.BuyNumber - model.UsageNumber,
		})
	response.SuccessReturnData(ctx, result)
}
