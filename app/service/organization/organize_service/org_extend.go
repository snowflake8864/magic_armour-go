package organize_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_ip_jump"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/global/enums/enum_safety_report"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

func CreateOrgExtendFactory() *OrgExtend {
	return &OrgExtend{masterModel: organization_models.CreateExtendFactory()}
}

type OrgExtend struct {
	masterModel *organization_models.ExtendModel
}

// InitDefault 初始化默认数据
func (crd *OrgExtend) InitDefault(orgId int64) organization_models.ExtendModel {
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

	return organization_models.ExtendModel{
		Id:                   orgId,
		PlugPortConfig:       json_helper.Encode(plugPortConfig),
		PlugLoginConfig:      json_helper.Encode(plugLoginConfig),
		LoginSolveConfig:     json_helper.Encode(loginSolveConfig),
		AlarmNoticeConfig:    json_helper.Encode(alarmNoticeConfig),
		HardwareNoticeConfig: json_helper.Encode(hardwareNoticeConfig),
		UpdatedAt:            date_helper.CurDatetime(),
	}
}

// DeleteByOrg 通过企业编号删除
func (crd *OrgExtend) DeleteByOrg(orgId int64) bool {
	return crd.masterModel.DeleteById(orgId)
}

// GetItemByOrg 通过企业编号获取
func (crd *OrgExtend) GetItemByOrg(orgId int64) *organization_models.ExtendModel {
	model, _ := crd.masterModel.GetItemById(orgId)
	if model == nil {
		model = &organization_models.ExtendModel{}
	}
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

// FormatPlugLoginConfig 格式化：封堵-口令破解配置
func (crd *OrgExtend) FormatPlugLoginConfig(configJson string) (data enum_plug.PlugLoginConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatLoginSolveConfig 格式化：口令破解配置
func (crd *OrgExtend) FormatLoginSolveConfig(configJson string) (data enum_login_solve.LoginSolveConfig) {
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

// FormatPwJumpConfig 格式化： 密码跳变配置
func (crd *OrgExtend) FormatPwJumpConfig(configJson string) (data enum_pw_jump.PwJumpConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// FormatIpJumpConfig 格式化： IP跳变配置
func (crd *OrgExtend) FormatIpJumpConfig(configJson string) (data enum_ip_jump.IpJumpConfig) {
	err := json_helper.ToStruct(configJson, &data)
	if err != nil {
		data.Format()
	}
	return data
}

// SetPlugOfPort 设置封堵配置： 端口虚开
func (crd *OrgExtend) SetPlugOfPort(orgId int64, config enum_plug.PlugPortConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["plug_port_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetPlugOfLogin 设置封堵配置： 口令破解
func (crd *OrgExtend) SetPlugOfLogin(orgId int64, config enum_plug.PlugLoginConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["plug_login_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetLoginSolve 设置口令破解配置配置
func (crd *OrgExtend) SetLoginSolve(orgId int64, config enum_login_solve.LoginSolveConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["login_solve_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetAlarmNoticeConfig 设置告警通知配置
func (crd *OrgExtend) SetAlarmNoticeConfig(orgId int64, config enum_safety_report.AlarmNoticeConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["alarm_notice_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetHardwareNoticeConfig 硬件资源通知配置
func (crd *OrgExtend) SetHardwareNoticeConfig(orgId int64, config enum_safety_report.HardwareNoticeConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["hardware_notice_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetPwJumpConfig 设置密码跳变
func (crd *OrgExtend) SetPwJumpConfig(orgId int64, config enum_pw_jump.PwJumpConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["pw_jump_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}

// SetIpJumpConfig 设置IP跳变
func (crd *OrgExtend) SetIpJumpConfig(orgId int64, config enum_ip_jump.IpJumpConfig) error {
	config.Format()
	updates := make(map[string]interface{})
	updates["ip_jump_config"] = json_helper.Encode(config)

	if crd.masterModel.UpdateById(orgId, updates) {
		return nil
	} else {
		return errors.New("保存失败")
	}
}
