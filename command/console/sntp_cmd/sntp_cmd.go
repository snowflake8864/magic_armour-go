package sntp_cmd

import (
	"fmt"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/sntp_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCmdFactory() *SntpCmd {
	return &SntpCmd{}
}

type SntpCmd struct {
}

// Sync 同步时间
func (crd *SntpCmd) Sync() {
	variable.ZapLog.Info("Sync start")
	addr := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformSntpAddr)
	if !string_helper.IsEmpty(addr) {
		//获取命令
		command := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.PlatformSntpCmd)

		if err, output := sntp_helper.Exec(command, addr); err != nil {
			variable.ZapLog.Warn(fmt.Sprintf("Sync fail: %s", err.Error()))
		} else {
			variable.ZapLog.Info(fmt.Sprintf("Sync ok: %s", output))
		}
	}

	variable.ZapLog.Info("Sync end")
}
