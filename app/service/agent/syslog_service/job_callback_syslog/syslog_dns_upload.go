package job_callback_syslog

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/alarm_service/job_callback_alarm"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateDnsUploadFactory() *SyslogDnsUpload {
	return &SyslogDnsUpload{
		masterModel: syslog_models.CreateDnsFactory(),
	}
}

type SyslogDnsUpload struct {
	syslogBase
	masterModel *syslog_models.DnsModel
}

// Add 添加
func (crd *SyslogDnsUpload) Add(model *syslog_models.DnsModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.ProcHash = strings.ToLower(model.ProcHash)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiLogSwitch); hsCli != nil {
		hsData := huanshi_platform.LogDns{
			LogBase: huanshi_platform.LogBase{
				OrgId:     model.OrgId,
				AgentId:   model.AgentId,
				CreatedAt: model.CreatedAt,
			},
			ProcHash:   model.ProcHash,
			ProcId:     model.ProcId,
			ProcDir:    model.ProcDir,
			DomainName: model.DomainName,
			TargetIp:   model.TargetIp,
		}
		if err := hsCli.SyncLogDns(hsData); err != nil {
			variable.ZapLog.Warn(err.Error())
		}
	}

	//外联告警
	if err := job_callback_alarm.CreateOutreachCalcFactory().Calc(model); err != nil {
		variable.ZapLog.Warn(err.Error())
	}
	return nil
}
