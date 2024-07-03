package alarm_cmd

import (
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCmdFactory() *AlarmCmd {
	return &AlarmCmd{}
}

type AlarmCmd struct {
}

// TodayStatistics  今天统计
func (crd *AlarmCmd) TodayStatistics() {
	variable.ZapLog.Sugar().Debug("TodayStatistics")
}

// YesterdayStatistics  昨天统计
func (crd *AlarmCmd) YesterdayStatistics() {
	variable.ZapLog.Sugar().Debug("DailyStatistics")
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
