package job_callback_syslog

import (
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateInnerCurdFactory() *InnerCurd {
	return &InnerCurd{
		masterModel: syslog_models.CreateInnerFactory(),
	}
}

type InnerCurd struct {
	masterModel *syslog_models.InnerModel
}
