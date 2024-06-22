package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateProcessUploadFactory() *AlarmProcessUpload {
	return &AlarmProcessUpload{
		masterModel: alarm_models.CreateProcessFactory(),
	}
}

type AlarmProcessUpload struct {
	alarmBase
	masterModel *alarm_models.ProcessModel
}

// Add 添加
func (crd *AlarmProcessUpload) Add(model *alarm_models.ProcessModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.ProcHash = strings.ToLower(model.ProcHash)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)
	model.ProcParam = path_helper.ToStandard(model.ProcParam)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchDll); hsCli != nil {
		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypeProcess,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			Dir: model.ProcDir,
		}
		if err := hsCli.SyncAlarmComm(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}
