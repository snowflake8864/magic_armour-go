package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/my_errors"
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
		masterModel: alarm_models.CreateUsbFactory(),
	}
}

type AlarmUsbUpload struct {
	alarmBase
	masterModel *alarm_models.UsbModel
}

// BatchAdd 新增
// @TODO 只保存告警和推送到幻势，不在对外设表做相关处理(原bug7675); linux 2.0.2.20版本在插拔时重新上传外接设备列表
func (crd *AlarmUsbUpload) BatchAdd(list *[]alarm_models.UsbModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []alarm_models.UsbModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		model.UsbHash = strings.ToLower(model.UsbHash)
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchUsb); hsCli != nil {

		var aliasMap = make(map[string]string)
		for _, model := range inserts {
			//外接设备别名
			if _, isOk := aliasMap[model.UsbHash]; !isOk {
				if usbModel := usb_models.CreateUsbFactory().GetItemByOrgAndAgentAndHash(model.OrgId, model.AgentId, model.UsbHash); usbModel != nil {
					aliasMap[model.UsbHash] = usbModel.Alias
				}
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
				PeripheralContact: aliasMap[model.UsbHash],
			}
			if err := hsCli.SyncAlarmComm(hsData); err != nil {
				variable.ZapLog.Warn(err.Error())
			}
		}
	}

	return nil
}
