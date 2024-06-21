package agent_post_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/service/agent/usb_service"
	"goskeleton/app/utils/helpers/json_helper"
	"strings"
)

func CreateUsbFactory() *UsbPost {
	return &UsbPost{}
}

type UsbPost struct {
}

// Save 保存
func (crd *UsbPost) Save(orgId int64, agentId, inputJson string) error {
	var params []usbData
	//忽略错误, 上传数据有为空情况
	_ = json_helper.ToAny(inputJson, &params)

	var models []usb_models.UsbModel
	var usbHashList []string
	if len(params) > 0 {
		for _, item := range params {
			model := usb_models.UsbModel{
				Id:    strings.ToLower(item.PeripheralEid),
				Name:  item.PeripheralName,
				Intro: item.PeripheralIntro,
				Type:  item.PeripheralType,
			}
			models = append(models, model)
			usbHashList = append(usbHashList, model.Id)
		}
	}

	if err := usb_service.CreateUsbCurdFactory().BatchAdd(orgId, agentId, usbHashList, models); err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
	}

	return nil
}
