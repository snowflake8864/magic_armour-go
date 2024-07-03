package sysglog_cmd

import (
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCmdFactory() *SyslogCmd {
	return &SyslogCmd{}
}

type SyslogCmd struct {
}

// Delete  删除
func (crd *SyslogCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete start")

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentSyslogExpire)
	expire := string_helper.ToInt(val)
	if expire == 1 {
		syslog_models.CreateDnsFactory().TruncateTable()
		syslog_models.CreateOuterFactory().TruncateTable()
		syslog_models.CreateInnerFactory().TruncateTable()
		syslog_models.CreateProcessFactory().TruncateTable()
		syslog_models.CreateLoginFactory().TruncateTable()
	} else {
		_, _, dayEnd, _ := date_helper.CalcByDay(-expire)
		variable.ZapLog.Sugar().Debug("Delete dayEnd", dayEnd)
		/*syslog_models.CreateDnsFactory().DeleteByConditions(syslog_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		syslog_models.CreateOuterFactory().DeleteByConditions(syslog_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		syslog_models.CreateInnerFactory().DeleteByConditions(syslog_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		syslog_models.CreateProcessFactory().DeleteByConditions(syslog_models.SearchConditions{
			CreateAtLte: dayEnd,
		})
		syslog_models.CreateLoginFactory().DeleteByConditions(syslog_models.SearchConditions{
			CreateAtLte: dayEnd,
		})*/
	}

	variable.ZapLog.Sugar().Debug("Delete end")
}
