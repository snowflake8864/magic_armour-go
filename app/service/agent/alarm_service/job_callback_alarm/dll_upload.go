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

func CreateDllUploadFactory() *DllUpload {
	return &DllUpload{
		masterModel: alarm_models.CreateDllFactory(),
	}
}

type DllUpload struct {
	alarmBase
	masterModel *alarm_models.DllModel
}

// Add 添加
func (crd *DllUpload) Add(model *alarm_models.DllModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.DllHash = strings.ToLower(model.DllHash)
	model.DllDir = path_helper.ToStandard(model.DllDir)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchDll); hsCli != nil {
		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypeDll,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			Dir:              model.DllDir,
			ExceptionProcess: model.ProcDir,
		}
		if err := hsCli.SyncAlarmComm(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}

// BatchAdd 新增
func (crd *DllUpload) BatchAdd(list *[]alarm_models.DllModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []alarm_models.DllModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		model.DllHash = strings.ToLower(model.DllHash)
		model.DllDir = path_helper.ToStandard(model.DllDir)
		model.ProcDir = path_helper.ToStandard(model.ProcDir)
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchDll); hsCli != nil {
		for _, model := range inserts {
			hsData := huanshi_platform.AlarmComm{
				AlarmBase: huanshi_platform.AlarmBase{
					Type:       huanshi_platform.AlarmTypeDll,
					OrgId:      model.OrgId,
					AgentId:    model.AgentId,
					Level:      model.Level,
					NoticeType: model.Type,
					NoticeTime: model.CreatedAt,
				},
				Dir:              model.DllDir,
				ExceptionProcess: model.ProcDir,
			}
			if err := hsCli.SyncAlarmComm(hsData); err != nil {
				variable.ZapLog.Sugar().Debug(err.Error())
			}
		}
	}

	return nil
}
