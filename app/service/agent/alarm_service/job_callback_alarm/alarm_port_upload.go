package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/plug_service/plug_calc_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreatePortUploadFactory() *AlarmPortUpload {
	return &AlarmPortUpload{
		masterModel: alarm_models.CreateAlarmPortFactory(),
	}
}

type AlarmPortUpload struct {
	alarmBase
	masterModel *alarm_models.AlarmPortModel
}

// Add 添加
func (crd *AlarmPortUpload) Add(model *alarm_models.AlarmPortModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchPort); hsCli != nil {
		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypePort,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			AttackIp:      model.AttackIp,
			DestinationIp: model.TargetIp,
			OpenPort:      model.TargetPort,
			RedirectIp:    model.RedirectIp,
		}
		if err := hsCli.SyncAlarmComm(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	//动态封堵
	sumWeight := model.Weight * model.AttackTotal
	if err := plug_calc_service.CreatePlugCalcFactory().CalcByPort(model.OrgId, model.AgentId, model.AttackIp, sumWeight, model.CreatedAt); err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
	}

	return nil
}
