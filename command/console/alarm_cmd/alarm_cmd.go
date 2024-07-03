package alarm_cmd

import "goskeleton/app/global/variable"

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
	variable.ZapLog.Sugar().Debug("DailyStatistics")
}

// EmailNotice  邮件通知
func (crd *AlarmCmd) EmailNotice() {
	variable.ZapLog.Sugar().Debug("EmailNotice")
}
