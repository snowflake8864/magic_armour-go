package job_callback_syslog

import (
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/huanshi_helper"
)

type syslogBase struct {
}

// getHuanshiCli 获取幻势客户端
func (crd *syslogBase) getHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
	return huanshi_helper.GetHuanshiCli(pushTypeFiled)
}
