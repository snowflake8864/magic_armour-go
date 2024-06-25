package job_callback_syslog

import (
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateLoginCurdFactory() *LoginCurd {
	return &LoginCurd{
		masterModel: syslog_models.CreateLoginFactory(),
	}
}

type LoginCurd struct {
	masterModel *syslog_models.LoginModel
}
