package plug_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
)

func CreatePlugSaveFactory() *AlarmPlugSave {
	return &AlarmPlugSave{
		masterModel: alarm_models.CreatePlugFactory(),
	}
}

type AlarmPlugSave struct {
	alarmBase
	masterModel *alarm_models.PlugModel
}

// Add 新增
func (crd *AlarmPlugSave) Add(plugModel *plug_models.DynamicModel) error {
	alarmType := enum_alarm.TypePlugDynamicFromPort
	if plugModel.Source == enum_plug.SourcePortVirtual {
		alarmType = enum_alarm.TypePlugDynamicFromLogin
	}

	modelAlarm := alarm_models.PlugModel{
		Id:           variable.SnowFlake.GetId(),
		OrgId:        plugModel.OrgId,
		AgentId:      plugModel.AgentId,
		Type:         alarmType,
		Level:        enum_alarm.LevelSerious,
		AgentAt:      plugModel.CreatedAt,
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    plugModel.CreatedAt,
		Ip:           plugModel.Ip,
		PlugType:     plugModel.PlugType,
		Direction:    plugModel.Direction,
		Duration:     plugModel.Duration,
		StartTime:    plugModel.StartAt,
		EndTime:      plugModel.EndAt,
	}
	if !crd.masterModel.InsertItem(&modelAlarm) {
		return errors.New("保存告警失败")
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchPlug); hsCli != nil {
		hsData := huanshi_platform.AlarmPlug{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypePlug,
				OrgId:      modelAlarm.OrgId,
				AgentId:    modelAlarm.AgentId,
				Level:      modelAlarm.Level,
				NoticeType: modelAlarm.Type,
				NoticeTime: modelAlarm.CreatedAt,
			},

			Pip:        modelAlarm.Ip,
			PStartTime: modelAlarm.StartTime,
			PEndTime:   modelAlarm.EndTime,
			PMinute:    modelAlarm.Duration,
			PType:      modelAlarm.PlugType,
			PTypeTxt:   enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", modelAlarm.PlugType, enum_plug.EnumMap),
		}

		if err := hsCli.SyncAlarmPlug(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}
	return nil
}
