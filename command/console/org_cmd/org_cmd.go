package org_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *OrgCmd {
	return &OrgCmd{}
}

type OrgCmd struct {
}

// TodayStatistics  今天统计
func (crd *OrgCmd) TodayStatistics() {
	variable.ZapLog.Sugar().Debug("TodayStatistics")
}

// YesterdayStatistics  昨天统计
func (crd *OrgCmd) YesterdayStatistics() {
	variable.ZapLog.Sugar().Debug("DailyStatistics")
}
