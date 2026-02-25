package alarm_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/service/agent/alarm_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func StatInStartAtAndEndAt(jsonStr string, startAt, endAt string) {
	model := alarm_stat_models.StatDailyModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Warn(err.Error())
		variable.ZapLog.Info(fmt.Sprintf("startAt=%s, endAt=%s", startAt, endAt))
		variable.ZapLog.Info(jsonStr)
		return
	}
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d,agentId=%s,startAt=%s,endAt=%s", model.OrgId, model.AgentId, startAt, endAt))

	conditions1 := alarm_models.SearchConditions{
		OrgId:   model.OrgId,
		AgentId: model.AgentId,
	}
	conditions2 := alarm_models.SearchConditions{
		OrgId:       model.OrgId,
		AgentId:     model.AgentId,
		CreateAtGte: startAt,
		CreateAtLte: endAt,
	}

	column1 := alarm_models.CreateProcessFactory().ColumnCalc(false, false, false, true)

	//进程告警
	ProcessFact := alarm_models.CreateProcessFactory()
	calc1 := ProcessFact.CalcByConditions(conditions1, column1)
	calc2 := ProcessFact.CalcByConditions(conditions2, column1)
	model.ProcessTotal = calc1.TotalNumber
	model.ProcessHandled = calc1.HandledNumber
	model.ProcessAddTotal = calc2.TotalNumber
	model.ProcessAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//模块告警
	DllFact := alarm_models.CreateDllFactory()
	calc1 = DllFact.CalcByConditions(conditions1, column1)
	calc2 = DllFact.CalcByConditions(conditions2, column1)
	model.DllTotal = calc1.TotalNumber
	model.DllHandled = calc1.HandledNumber
	model.DllAddTotal = calc2.TotalNumber
	model.DllAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//防篡改告警
	TamperFact := alarm_models.CreateTamperFactory()
	calc1 = TamperFact.CalcByConditions(conditions1, column1)
	calc2 = TamperFact.CalcByConditions(conditions2, column1)
	model.TamperTotal = calc1.TotalNumber
	model.TamperHandled = calc1.HandledNumber
	model.TamperAddTotal = calc2.TotalNumber
	model.TamperAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//勒索告警
	BlackmailFact := alarm_models.CreateBlackmailFactory()
	calc1 = BlackmailFact.CalcByConditions(conditions1, column1)
	calc2 = BlackmailFact.CalcByConditions(conditions2, column1)
	model.BlackmailTotal = calc1.TotalNumber
	model.BlackmailHandled = calc1.HandledNumber
	model.BlackmailAddTotal = calc2.TotalNumber
	model.BlackmailAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//端口虚开告警
	PortFact := alarm_models.CreatePortFactory()
	calc1 = PortFact.CalcByConditions(conditions1, column1)
	calc2 = PortFact.CalcByConditions(conditions2, column1)
	model.PortTotal = calc1.TotalNumber
	model.PortHandled = calc1.HandledNumber
	model.PortAddTotal = calc2.TotalNumber
	model.PortAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//封堵告警
	PlugFact := alarm_models.CreatePlugFactory()
	calc1 = PlugFact.CalcByConditions(conditions1, column1)
	calc2 = PlugFact.CalcByConditions(conditions2, column1)
	model.PlugTotal = calc1.TotalNumber
	model.PlugHandled = calc1.HandledNumber
	model.PlugAddTotal = calc2.TotalNumber
	model.PlugAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//外设告警
	UsbFact := alarm_models.CreateUsbFactory()
	calc1 = UsbFact.CalcByConditions(conditions1, column1)
	calc2 = UsbFact.CalcByConditions(conditions2, column1)
	model.UsbTotal = calc1.TotalNumber
	model.UsbHandled = calc1.HandledNumber
	model.UsbAddTotal = calc2.TotalNumber
	model.UsbAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//外联告警
	OutreachFact := alarm_models.CreateOutreachFactory()
	calc1 = OutreachFact.CalcByConditions(conditions1, column1)
	calc2 = OutreachFact.CalcByConditions(conditions2, column1)
	model.OutreachTotal = calc1.TotalNumber
	model.OutreachHandled = calc1.HandledNumber
	model.OutreachAddTotal = calc2.TotalNumber
	model.OutreachAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//口令破解告警
	LoginFact := alarm_models.CreateLoginFactory()
	calc1 = LoginFact.CalcByConditions(conditions1, column1)
	calc2 = LoginFact.CalcByConditions(conditions2, column1)
	model.LoginTotal = calc1.TotalNumber
	model.LoginHandled = calc1.HandledNumber
	model.LoginAddTotal = calc2.TotalNumber
	model.LoginAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//自保告警
	SelfFact := alarm_models.CreateSelfFactory()
	calc1 = SelfFact.CalcByConditions(conditions1, column1)
	calc2 = SelfFact.CalcByConditions(conditions2, column1)
	model.SelfTotal = calc1.TotalNumber
	model.SelfHandled = calc1.HandledNumber
	model.SelfAddTotal = calc2.TotalNumber
	model.SelfAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	//病毒查杀告警
	VirusFact := alarm_models.CreateVirusFactory()
	calc1 = VirusFact.CalcByConditions(conditions1, column1)
	calc2 = VirusFact.CalcByConditions(conditions2, column1)
	model.VirusTotal = calc1.TotalNumber
	model.VirusHandled = calc1.HandledNumber
	model.VirusAddTotal = calc2.TotalNumber
	model.VirusAddHandled = calc2.HandledNumber
	model.Total += calc1.TotalNumber
	model.Handled += calc1.HandledNumber
	model.AddTotal += calc2.TotalNumber
	model.AddHandled += calc2.HandledNumber

	if !alarm_service.CreateStatDailyFactory().Add(&model) {
		variable.ZapLog.Warn("fail")
		return
	}
}
