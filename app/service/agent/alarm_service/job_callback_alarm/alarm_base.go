package job_callback_alarm

import (
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/huanshi_helper"
)

type alarmBase struct {
}

// getHuanshiCli 获取幻势客户端
func (crd *alarmBase) getHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
	return huanshi_helper.GetHuanshiCli(pushTypeFiled)
}
