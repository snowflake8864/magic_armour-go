package threaten_job

import (
	"fmt"
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
		variable.ZapLog.Info("HandelUnknown fail")
	}
	return
}

func HandelRewire(jsonStr string) {
	model := threaten_models.CacheModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Warn(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}
	variable.ZapLog.Info(fmt.Sprintf("HandelRewire hash:%s, level:%d, dataFrom:%d", model.Id, model.Level, model.DataFrom))

	if !threaten_service.CreateCacheFactory().Rewrite(&model) {
		variable.ZapLog.Info("HandelRewire fail")
	}
	return
}

func HandelImport(filename string) {
	variable.ZapLog.Info(fmt.Sprintf("HandelImport filename:%s", filename))

	if err := threaten_service.CreateImportFactory().ParseExcel(filename); err != nil {
		variable.ZapLog.Warn(fmt.Sprintf("HandelImport fail: %s", err.Error()))
	}
	return
}
