package job_callback_syslog

import (
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateProcessCurdFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: syslog_models.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	masterModel *syslog_models.ProcessModel
}
