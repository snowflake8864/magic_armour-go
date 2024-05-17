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

func CreateSelfUploadFactory() *AlarmSelfUpload {
	return &AlarmSelfUpload{
		masterModel: alarm_models.CreateAlarmSelfFactory(),
	}
}

type AlarmSelfUpload struct {
	alarmBase
	masterModel *alarm_models.AlarmSelfModel
}

// Add 添加
func (crd *AlarmSelfUpload) Add(model *alarm_models.AlarmSelfModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.ProcHash = strings.ToLower(model.ProcHash)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)
	model.FileDir = path_helper.ToStandard(model.FileDir)
	model.TargetDir = path_helper.ToStandard(model.TargetDir)
	model.ProcParam = path_helper.ToStandard(model.ProcParam)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.GetHuanshiCli(); hsCli != nil {
		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypeSelf,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
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
