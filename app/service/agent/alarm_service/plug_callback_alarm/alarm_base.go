package plug_callback_alarm

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/string_helper"
)

type alarmBase struct {
}

// getHuanshiCli 获取幻势客户端
func (crd *alarmBase) getHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
	dictIds := []string{
		enum_dictionary.HuanshiCallbackUrl,
		enum_dictionary.HuanshiCallIp,
		enum_dictionary.HuanshiSwitch,
		enum_dictionary.HuanshiUrl,
		enum_dictionary.HuanshiAlarmSwitch,
		enum_dictionary.HuanshiAlarmIp,
		enum_dictionary.HuanshiAlarmPort,
		enum_dictionary.HuanshiLogSwitch,
		enum_dictionary.HuanshiLogIp,
		enum_dictionary.HuanshiLogPort,
		pushTypeFiled,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if len(vals) != 11 {
		return nil
	}

	if string_helper.ToInt(vals[enum_dictionary.HuanshiSwitch]) == enums.SwitchOff ||
		string_helper.ToInt(vals[enum_dictionary.HuanshiAlarmSwitch]) == enums.SwitchOff ||
		string_helper.ToInt(vals[pushTypeFiled]) == enums.SwitchOff {
		return nil
	}

	config := huanshi_platform.RequestConfig{}
	config.HuanshiCallbackUrl = vals[enum_dictionary.HuanshiCallbackUrl]
	config.HuanshiCallIp = vals[enum_dictionary.HuanshiCallIp]
	config.HuanshiUrl = vals[enum_dictionary.HuanshiUrl]
	config.HuanshiAlarmIp = vals[enum_dictionary.HuanshiAlarmIp]
	config.HuanshiAlarmPort = vals[enum_dictionary.HuanshiAlarmPort]
	config.HuanshiLogIp = vals[enum_dictionary.HuanshiLogIp]
	config.HuanshiLogPort = vals[enum_dictionary.HuanshiLogPort]

	return huanshi_platform.CreateRequestFactory(config)
}
