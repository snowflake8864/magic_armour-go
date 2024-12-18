package huanshi_helper

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/string_helper"
)

// GetHuanshiCli 获取幻势客户端
func GetHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
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
		enum_dictionary.HuanshiThirdSwitch,
		enum_dictionary.HuanshiThirdIp,
		enum_dictionary.HuanshiThirdPort,
	}
	if !string_helper.IsEmpty(pushTypeFiled) {
		dictIds = append(dictIds, pushTypeFiled)
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if len(vals) != len(dictIds) {
		return nil
	}
	// 要推送信息对应已关闭
	if !string_helper.IsEmpty(pushTypeFiled) && string_helper.ToInt(vals[pushTypeFiled]) == enums.SwitchOff {
		return nil
	}

	config := huanshi_platform.RequestConfig{}
	config.HuanshiSwitch = string_helper.ToInt(vals[enum_dictionary.HuanshiSwitch])
	config.HuanshiCallbackUrl = vals[enum_dictionary.HuanshiCallbackUrl]
	config.HuanshiCallIp = vals[enum_dictionary.HuanshiCallIp]
	config.HuanshiUrl = vals[enum_dictionary.HuanshiUrl]

	config.HuanshiAlarmSwitch = string_helper.ToInt(vals[enum_dictionary.HuanshiAlarmSwitch])
	config.HuanshiAlarmIp = vals[enum_dictionary.HuanshiAlarmIp]
	config.HuanshiAlarmPort = vals[enum_dictionary.HuanshiAlarmPort]

	config.HuanshiLogSwitch = string_helper.ToInt(vals[enum_dictionary.HuanshiLogSwitch])
	config.HuanshiLogIp = vals[enum_dictionary.HuanshiLogIp]
	config.HuanshiLogPort = vals[enum_dictionary.HuanshiLogPort]

	config.HuanshiThirdSwitch = string_helper.ToInt(vals[enum_dictionary.HuanshiThirdSwitch])
	config.HuanshiThirdIp = vals[enum_dictionary.HuanshiThirdIp]
	config.HuanshiThirdPort = vals[enum_dictionary.HuanshiThirdPort]

	return huanshi_platform.CreateRequestFactory(config)
}
