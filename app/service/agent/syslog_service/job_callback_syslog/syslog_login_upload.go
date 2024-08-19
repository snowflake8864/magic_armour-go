package job_callback_syslog

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/alarm_service/job_callback_alarm"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateLoginUploadFactory() *SyslogLoginUpload {
	return &SyslogLoginUpload{
		masterModel: syslog_models.CreateLoginFactory(),
	}
}

type SyslogLoginUpload struct {
	syslogBase
	masterModel *syslog_models.LoginModel
}

// Add 添加
func (crd *SyslogLoginUpload) Add(model *syslog_models.LoginModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(); hsCli != nil {
		if err := hsCli.SyncLogLogin(crd.initHs(model)); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	//登录失败
	if model.Status == enums.IsNot {
		//计算口令破解告警
		if err := job_callback_alarm.CreateLoginCalcFactory().Calc(model.OrgId, model.AgentId, model.SourceIp, model.CreatedAt); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}
	return nil
}

// BatchAddSuccess 批量新增成功
func (crd *SyslogLoginUpload) BatchAddSuccess(list *[]syslog_models.LoginModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []syslog_models.LoginModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(); hsCli != nil {
		for _, model := range inserts {
			if err := hsCli.SyncLogLogin(crd.initHs(&model)); err != nil {
				variable.ZapLog.Sugar().Debug(err.Error())
			}
		}
	}

	return nil
}

func (crd *SyslogLoginUpload) initHs(model *syslog_models.LoginModel) huanshi_platform.LogLogin {

	return huanshi_platform.LogLogin{
		LogBase: huanshi_platform.LogBase{
			OrgId:     model.OrgId,
			AgentId:   model.AgentId,
			CreatedAt: model.CreatedAt,
		},
		SourceIp: model.SourceIp,
		Username: model.Username,
		Type:     model.Type,
		Status:   model.Status,
	}
}
