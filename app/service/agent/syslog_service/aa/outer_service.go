package job_callback_syslog

import (
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateOuterCurdFactory() *OuterCurd {
	return &OuterCurd{
		masterModel: syslog_models.CreateOuterFactory(),
	}
}

type OuterCurd struct {
	masterModel *syslog_models.OuterModel
}
