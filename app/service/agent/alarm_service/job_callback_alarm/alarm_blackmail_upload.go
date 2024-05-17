package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateBlackmailUploadFactory() *AlarmBlackmailUpload {
	return &AlarmBlackmailUpload{
		masterModel: alarm_models.CreateAlarmBlackmailFactory(),
	}
}

type AlarmBlackmailUpload struct {
	alarmBase
	masterModel *alarm_models.AlarmBlackmailModel
}

// Add 添加
func (crd *AlarmBlackmailUpload) Add(model *alarm_models.AlarmBlackmailModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.ProcHash = strings.ToLower(model.ProcHash)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)
	model.FileDir = path_helper.ToStandard(model.FileDir)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.GetHuanshiCli(); hsCli != nil {
		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypeBlackmail,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Level,
				NoticeTime: model.CreatedAt,
			},
			Dir:              model.FileDir,
			ExceptionProcess: model.ProcDir,
		}
		if err := hsCli.SyncAlarmComm(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}
