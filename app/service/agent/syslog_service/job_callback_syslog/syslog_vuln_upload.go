package job_callback_syslog

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateVulnUploadFactory() *SyslogVulnUpload {
	return &SyslogVulnUpload{
		masterModel: syslog_models.CreateVulnFactory(),
	}
}

type SyslogVulnUpload struct {
	syslogBase
	masterModel *syslog_models.VulnModel
}

// Add 添加
func (crd *SyslogVulnUpload) Add(model *syslog_models.VulnModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	return nil
}
