package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateUsbUploadFactory() *AlarmUsbUpload {
	return &AlarmUsbUpload{
		masterModel: alarm_models.CreateAlarmUsbFactory(),
	}
}

type AlarmUsbUpload struct {
	alarmBase
	masterModel *alarm_models.AlarmUsbModel
}

// Add 添加
// @TODO 只保存告警和推送到幻势，不在对外设表做相关处理(原bug7675); linux版本2.0.2.20，在插拔时重新上传外接设备列表
func (crd *AlarmUsbUpload) Add(model *alarm_models.AlarmUsbModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.UsbHash = strings.ToLower(model.UsbHash)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchUsb); hsCli != nil {
		//外接设备别名
		var usbAlias string
		if usbModel := usb_models.CreateUsbFactory().GetItemByOrgAndAgentAndHash(model.OrgId, model.AgentId, model.UsbHash); usbModel != nil {
			usbAlias = usbModel.Alias
		}

		hsData := huanshi_platform.AlarmComm{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypeUsb,
				OrgId:      model.OrgId,
				AgentId:    model.AgentId,
				Level:      model.Level,
				NoticeType: model.Type,
				NoticeTime: model.CreatedAt,
			},
			PeripheralEid:     model.UsbHash,
			PeripheralName:    model.UsbName,
			PeripheralRemark:  model.UsbIntro,
			PeripheralContact: usbAlias,
		}
		if err := hsCli.SyncAlarmComm(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}
