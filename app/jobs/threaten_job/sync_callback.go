package threaten_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func HandelUnknown(jsonStr string) {
	idArr := json_helper.ToSliceString(jsonStr)
	if len(idArr) == 0 {
		variable.ZapLog.Info(jsonStr)
		return
	}
	if !threaten_service.CreateCacheFactory().BatchUnknown(idArr) {
		variable.ZapLog.Info(jsonStr)
	}
	return
}

func HandelRewire(jsonStr string) {
	model := threaten_models.CacheModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if !threaten_service.CreateCacheFactory().Rewrite(&model) {
		variable.ZapLog.Info(jsonStr)
	}
	return
}
