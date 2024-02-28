package grant_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/dictionary_service"
	"strconv"
)

func CreateCurdFactory() *GrantCurd {
	return &GrantCurd{
		masterModel: model_platform.CreateGrantFactory(),
		moduleModel: model_platform.CreateGrantModuleFactory(),
		UsedModel:   model_platform.CreateGrantUsedFactory(),
	}
}

type GrantCurd struct {
	masterModel *model_platform.GrantModel
	moduleModel *model_platform.GrantModuleModel
	UsedModel   *model_platform.GrantUsedModel
}

// AutoGrant 自动授权给终端
func (crd *GrantCurd) AutoGrant(agentId string) {
	variable.ZapLog.Info("自动授权----------------")
	//自动授权开关
	autoSwitch, _ := strconv.Atoi(dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.SystemGrantAutoSwitch))
	if autoSwitch == enums.SwitchOff {
		return
	}

}
