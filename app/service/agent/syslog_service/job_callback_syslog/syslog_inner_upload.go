package job_callback_syslog

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateInnerUploadFactory() *SyslogInnerUpload {
	return &SyslogInnerUpload{
		masterModel: syslog_models.CreateInnerFactory(),
	}
}

type SyslogInnerUpload struct {
	syslogBase
	masterModel *syslog_models.InnerModel
}

// Add 添加
func (crd *SyslogInnerUpload) Add(model *syslog_models.InnerModel) error {
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
	if hsCli := crd.getHuanshiCli(); hsCli != nil {
		hsData := huanshi_platform.LogNetwork{
			LogBase: huanshi_platform.LogBase{
				OrgId:     model.OrgId,
				AgentId:   model.AgentId,
				CreatedAt: model.CreatedAt,
				LogType:   huanshi_platform.LogTypeInner,
			},
			ProcHash:         model.ProcHash,
			ProcId:           model.ProcId,
			ProcDir:          model.ProcDir,
			TargetIp:         model.ProcDir,
			TargetPort:       model.TargetPort,
			TargetPortStatus: model.TargetPortStatus,
			SourceIp:         model.SourceIp,
			SourcePort:       model.SourcePort,
			Proto:            model.Proto,
		}
		if err := hsCli.SyncLogNetwork(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}
