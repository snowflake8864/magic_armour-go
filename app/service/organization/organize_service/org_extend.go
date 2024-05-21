package organize_service

import (
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_safety_report"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

func CreateOrgExtendFactory() *OrgExtend {
	return &OrgExtend{masterModel: organization_models.CreateOrgExtendFactory()}
}

type OrgExtend struct {
	masterModel *organization_models.OrgExtendModel
}

// InitDefault 初始化默认数据
func (crd *OrgExtend) InitDefault(orgId int64) organization_models.OrgExtendModel {
	plugPortConfig := enum_plug.PlugPortConfig{}
	plugPortConfig.Format()
	plugLoginConfig := enum_plug.PlugLoginConfig{}
	plugLoginConfig.Format()
	loginSolveConfig := enum_login_solve.LoginSolveConfig{}
	loginSolveConfig.Format()
	alarmNoticeConfig := enum_safety_report.AlarmNoticeConfig{}
	alarmNoticeConfig.Format()
	hardwareNoticeConfig := enum_safety_report.HardwareNoticeConfig{}
	hardwareNoticeConfig.Format()

	return organization_models.OrgExtendModel{
		Id:                   orgId,
		PlugPortConfig:       json_helper.Encode(plugPortConfig),
		PlugLoginConfig:      json_helper.Encode(plugLoginConfig),
		LoginSolveConfig:     json_helper.Encode(loginSolveConfig),
		AlarmNoticeConfig:    json_helper.Encode(loginSolveConfig),
		HardwareNoticeConfig: json_helper.Encode(hardwareNoticeConfig),
		UpdatedAt:            date_helper.CurDatetime(),
	}
}

// DeleteByOrg 通过企业编号删除
func (crd *OrgExtend) DeleteByOrg(orgId int64) bool {
	return crd.masterModel.DeleteById(orgId)
}

// GetItemByOrg 通过企业编号获取
func (crd *OrgExtend) GetItemByOrg(orgId int64) (model *organization_models.OrgExtendModel) {
	model, _ = crd.masterModel.GetItemById(orgId)
	return model
}

// FormatPlugPortConfig 格式化：封堵-端口虚开配置
func (crd *OrgExtend) FormatPlugPortConfig(configJson string) (data enum_plug.PlugPortConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatLoginSolveConfig 格式化：封堵-口令破解配置
func (crd *OrgExtend) FormatLoginSolveConfig(configJson string) (data enum_login_solve.LoginSolveConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatPlugLoginConfig 格式化：封堵-口令破解配置
func (crd *OrgExtend) FormatPlugLoginConfig(configJson string) (data enum_plug.PlugLoginConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatAlarmNoticeConfig 格式化：告警通知配置
func (crd *OrgExtend) FormatAlarmNoticeConfig(configJson string) (data enum_safety_report.AlarmNoticeConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatHardwareNoticeConfig 格式化： 硬件资源通知配置
func (crd *OrgExtend) FormatHardwareNoticeConfig(configJson string) (data enum_safety_report.HardwareNoticeConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}
