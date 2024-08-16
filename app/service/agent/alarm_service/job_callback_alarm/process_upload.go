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

// BatchAdd 新增
func (crd *AlarmProcessUpload) BatchAdd(list *[]alarm_models.ProcessModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []alarm_models.ProcessModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		model.ProcHash = strings.ToLower(model.ProcHash)
		model.ProcDir = path_helper.ToStandard(model.ProcDir)
		model.ProcParam = path_helper.ToStandard(model.ProcParam)

		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchProcess); hsCli != nil {
		for _, model := range inserts {
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
	}

	return nil
}
