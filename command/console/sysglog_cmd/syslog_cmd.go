package sysglog_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *SyslogCmd {
	return &SyslogCmd{}
}

type SyslogCmd struct {
}

// Delete  删除
func (crd *SyslogCmd) Delete() {
	variable.ZapLog.Sugar().Debug("DailyStatistics")
}
