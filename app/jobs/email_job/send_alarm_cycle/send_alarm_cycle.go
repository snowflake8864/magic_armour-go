package send_alarm_cycle

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_safety_report"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_stat_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/gomail_factory"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"time"
)

func Send(orgId int64) {
	variable.ZapLog.Sugar().Info(fmt.Sprintf("企业orgId=%d", orgId))

	//企业扩展
	OrgExtendFact := organize_service.CreateOrgExtendFactory()
	extendModel := OrgExtendFact.GetItemByOrg(orgId)
	if extendModel == nil || extendModel.Id == 0 {
		variable.ZapLog.Sugar().Debug("企业扩展不存在, orgId=", orgId)
		return
	}
	noticeConfig := OrgExtendFact.FormatAlarmNoticeConfig(extendModel.AlarmNoticeConfig)
	if noticeConfig.EmailSwitch == enums.SwitchOff {
		variable.ZapLog.Sugar().Debug("开关关闭, orgId=", orgId)
		return
	}
	if string_helper.IsEmpty(noticeConfig.Email) {
		variable.ZapLog.Sugar().Debug("邮箱为空, orgId=", orgId)
		return
	}
	if len(noticeConfig.CycleMod) == 0 {
		variable.ZapLog.Sugar().Debug("没有选中告警模块, orgId=", orgId)
		return
	}

	//邮件参数
	mailParam := gomail_factory.AlarmCycle{}
	//统计告警数据条件
	conditions := alarm_stat_models.SearchConditions{
		OrgId: orgId,
	}

	//按照发送周期处理
	if noticeConfig.CycleInterval == enum_safety_report.AlarmNoticeCycleDay {
		yesterday, dayStart, dayEnd, month := date_helper.CalcByDay(-1)
		conditions.StatMonth = month
		conditions.StatDay = yesterday

		mailParam.DateStart = dayStart
		mailParam.DateEnd = dayEnd

	} else if noticeConfig.CycleInterval == enum_safety_report.AlarmNoticeCycleWeek {
		if d := date_helper.WeekdayCurr(); d != time.Monday {
			variable.ZapLog.Sugar().Debug("今天不是周一, orgId=", orgId)
			return
		}
		//上周一
		preMonday, _, _, _ := date_helper.CalcByDay(-7)
		//上周日
		preSunday, _, _, _ := date_helper.CalcByDay(-1)

		conditions.StatDayGte = preMonday
		conditions.StatDayLte = preSunday

		mailParam.DateStart = preMonday
		mailParam.DateEnd = preSunday

	} else if noticeConfig.CycleInterval == enum_safety_report.AlarmNoticeCycleMonth {
		if d := date_helper.DayInCurrMonth(); d != 1 {
			variable.ZapLog.Sugar().Debug("今天不是1号, orgId=", orgId)
			//return
		}
		//上个月
		month, monthFirstDay, monthLastDay := date_helper.CalcByMonth(-1)
		conditions.StatMonth = month

		mailParam.DateStart = monthFirstDay
		mailParam.DateEnd = monthLastDay
	} else {
		variable.ZapLog.Sugar().Debug("发送周期错误, orgId=", orgId)
		return
	}

	//统计
	alarmFact := alarm_stat_models.CreateStatDailyFactory()
	calc := alarmFact.CalcByConditions(conditions, alarmFact.StatColumnOfAllType(false, false, true, false, false))

	//模块和字段映射
	moduleMap := map[int]int{
		enum_alarm.ModuleProcess:   calc.ProcessAddTotal,
		enum_alarm.ModuleDll:       calc.DllAddTotal,
		enum_alarm.ModuleTamper:    calc.TamperAddTotal,
		enum_alarm.ModuleBlackmail: calc.BlackmailAddTotal,
		enum_alarm.ModulePort:      calc.PortAddTotal,
		enum_alarm.ModulePlug:      calc.PlugAddTotal,
		enum_alarm.ModuleUsb:       calc.UsbAddTotal,
		enum_alarm.ModuleOutreach:  calc.OutreachAddTotal,
		enum_alarm.ModuleLogin:     calc.LoginAddTotal,
		enum_alarm.ModuleSelf:      calc.SelfAddTotal,
	}

	for _, module := range noticeConfig.CycleMod {
		fmt.Println("module: ", module)
		if val, isOk := moduleMap[module]; isOk {
			mailParam.AlarmTotal += val
			mailParam.AlarmList = append(mailParam.AlarmList, gomail_factory.AlarmCycleInfo{
				Name: enums.Get1stIsStrAnd2ndIsIntReturnStr("module", module, enum_alarm.EnumMapInt),
				Add:  val,
			})
		}
	}

	//发送邮件
	subject, context := gomail_factory.CreateTemplateFactory().GetAlarmCycle(&mailParam)
	if err := gomail_factory.CreateMailerFactory().SendHtmlMsg(noticeConfig.Email, subject, context); err != nil {
		variable.ZapLog.Sugar().Debug("发送邮件错误, orgId=", orgId, "【error】", err.Error())
		return
	}
}
