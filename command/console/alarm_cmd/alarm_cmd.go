package alarm_cmd

import (
	"fmt"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/alarm_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCmdFactory() *AlarmCmd {
	return &AlarmCmd{}
}

type AlarmCmd struct {
}

func (crd *AlarmCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *AlarmCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// TodayStatistics  今天统计
func (crd *AlarmCmd) TodayStatistics() {
	variable.ZapLog.Sugar().Debug("TodayStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(0)
	crd.handleStat(statDay, startAt, endAt, statMonth)

	variable.ZapLog.Sugar().Debug("TodayStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

// YesterdayStatistics  昨天统计
func (crd *AlarmCmd) YesterdayStatistics() {
	variable.ZapLog.Sugar().Debug("DailyStatistics start")

	statDay, startAt, endAt, statMonth := date_helper.CalcByDay(-1)
	crd.handleStat(statDay, startAt, endAt, statMonth)

	variable.ZapLog.Sugar().Debug("DailyStatistics end startAt:" + startAt + ", endAt:" + endAt)
}

func (crd *AlarmCmd) handleStat(statDay, startAt, endAt string, statMonth int) {
	orgModels := crd.getOrgModels()

	if len(orgModels) > 0 {
		now := date_helper.CurDatetime()
		for _, orgModel := range orgModels {
			if agentModels := crd.getAgentModels(orgModel.Id); len(agentModels) > 0 {
				for _, agentModel := range agentModels {
					param := alarm_stat_models.StatDailyModel{
						OrgId:     orgModel.Id,
						AgentId:   agentModel.AgentId,
						StatDay:   statDay,
						StatMonth: statMonth,
						CreatedAt: now,
					}

					//计算告警
					crd.calcAlarm(&param, startAt, endAt)

					alarm_service.CreateStatDailyFactory().Add(&param)
				}
			}
		}
	}
}

func (crd *AlarmCmd) calcAlarm(param *alarm_stat_models.StatDailyModel, startAt, endAt string) {
	conditions1 := alarm_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
	}
	conditions2 := alarm_models.SearchConditions{
		OrgId:       param.OrgId,
		AgentId:     param.AgentId,
		CreateAtGte: startAt,
		CreateAtLte: endAt,
	}

	column1 := alarm_models.CreateProcessFactory().ColumnCalc(false, false, false, true)

	//进程告警
	ProcessFact := alarm_models.CreateProcessFactory()
	fmt.Println("ProcessFact=================")
	calc1 := ProcessFact.CalcByConditions(conditions1, column1)
	calc2 := ProcessFact.CalcByConditions(conditions2, column1)
	param.ProcessTotal = calc1.TotalNumber
	param.ProcessHandled = calc1.HandledNumber
	param.ProcessAddTotal = calc2.TotalNumber
	param.ProcessAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//模块告警
	DllFact := alarm_models.CreateDllFactory()
	calc1 = DllFact.CalcByConditions(conditions1, column1)
	calc2 = DllFact.CalcByConditions(conditions2, column1)
	param.DllTotal = calc1.TotalNumber
	param.DllHandled = calc1.HandledNumber
	param.DllAddTotal = calc2.TotalNumber
	param.DllAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//防篡改告警
	TamperFact := alarm_models.CreateTamperFactory()
	calc1 = TamperFact.CalcByConditions(conditions1, column1)
	calc2 = TamperFact.CalcByConditions(conditions2, column1)
	param.TamperTotal = calc1.TotalNumber
	param.TamperHandled = calc1.HandledNumber
	param.TamperAddTotal = calc2.TotalNumber
	param.TamperAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//勒索告警
	BlackmailFact := alarm_models.CreateBlackmailFactory()
	calc1 = BlackmailFact.CalcByConditions(conditions1, column1)
	calc2 = BlackmailFact.CalcByConditions(conditions2, column1)
	param.BlackmailTotal = calc1.TotalNumber
	param.BlackmailHandled = calc1.HandledNumber
	param.BlackmailAddTotal = calc2.TotalNumber
	param.BlackmailAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//端口虚开告警
	PortFact := alarm_models.CreatePortFactory()
	calc1 = PortFact.CalcByConditions(conditions1, column1)
	calc2 = PortFact.CalcByConditions(conditions2, column1)
	param.PortTotal = calc1.TotalNumber
	param.PortHandled = calc1.HandledNumber
	param.PortAddTotal = calc2.TotalNumber
	param.PortAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//封堵告警
	PlugFact := alarm_models.CreatePlugFactory()
	calc1 = PlugFact.CalcByConditions(conditions1, column1)
	calc2 = PlugFact.CalcByConditions(conditions2, column1)
	param.PlugTotal = calc1.TotalNumber
	param.PlugHandled = calc1.HandledNumber
	param.PlugAddTotal = calc2.TotalNumber
	param.PlugAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//外设告警
	UsbFact := alarm_models.CreateUsbFactory()
	calc1 = UsbFact.CalcByConditions(conditions1, column1)
	calc2 = UsbFact.CalcByConditions(conditions2, column1)
	param.UsbTotal = calc1.TotalNumber
	param.UsbHandled = calc1.HandledNumber
	param.UsbAddTotal = calc2.TotalNumber
	param.UsbAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//外联告警
	OutreachFact := alarm_models.CreateOutreachFactory()
	calc1 = OutreachFact.CalcByConditions(conditions1, column1)
	calc2 = OutreachFact.CalcByConditions(conditions2, column1)
	param.OutreachTotal = calc1.TotalNumber
	param.OutreachHandled = calc1.HandledNumber
	param.OutreachAddTotal = calc2.TotalNumber
	param.OutreachAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//口令破解告警
	LoginFact := alarm_models.CreateLoginFactory()
	calc1 = LoginFact.CalcByConditions(conditions1, column1)
	calc2 = LoginFact.CalcByConditions(conditions2, column1)
	param.LoginTotal = calc1.TotalNumber
	param.LoginHandled = calc1.HandledNumber
	param.LoginAddTotal = calc2.TotalNumber
	param.LoginAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

	//自保告警
	SelfFact := alarm_models.CreateSelfFactory()
	calc1 = SelfFact.CalcByConditions(conditions1, column1)
	calc2 = SelfFact.CalcByConditions(conditions2, column1)
	param.SelfTotal = calc1.TotalNumber
	param.SelfHandled = calc1.HandledNumber
	param.SelfAddTotal = calc2.TotalNumber
	param.SelfAddHandled = calc2.HandledNumber
	param.Total += calc1.TotalNumber
	param.Handled += calc1.HandledNumber
	param.AddTotal += calc2.TotalNumber
	param.AddHandled += calc2.HandledNumber

}

// Delete  删除
func (crd *AlarmCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete start")

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentAlarmExpire)
	expire := string_helper.ToInt(val)
	variable.ZapLog.Sugar().Debug("Delete expire:", expire)
	if expire == 1 {
		alarm_models.CreateProcessFactory().TruncateTable()
		alarm_models.CreateDllFactory().TruncateTable()
		alarm_models.CreateTamperFactory().TruncateTable()
		alarm_models.CreateBlackmailFactory().TruncateTable()
		alarm_models.CreatePortFactory().TruncateTable()
		alarm_models.CreatePlugFactory().TruncateTable()
		alarm_models.CreateUsbFactory().TruncateTable()
		alarm_models.CreateOutreachFactory().TruncateTable()
		alarm_models.CreateLoginFactory().TruncateTable()
		alarm_models.CreateSelfFactory().TruncateTable()
	} else {
		_, _, dayEnd, _ := date_helper.CalcByDay(-expire)
		variable.ZapLog.Sugar().Debug("Delete dayEnd:", dayEnd)
		alarm_models.CreateProcessFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateDllFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateTamperFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateBlackmailFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreatePortFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreatePlugFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateUsbFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateOutreachFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateLoginFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		alarm_models.CreateSelfFactory().DeleteByConditions(alarm_models.SearchConditions{
			CreateAtLte: dayEnd,
		})

	}
	variable.ZapLog.Sugar().Debug("Delete end")
}

// EmailNotice  邮件通知
func (crd *AlarmCmd) EmailNotice() {
	variable.ZapLog.Sugar().Debug("EmailNotice")
}
