package job_callback_alarm

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateVirusUploadFactory() *VirusUpload {
	return &VirusUpload{
		masterModel: alarm_models.CreateVirusFactory(),
	}
}

type VirusUpload struct {
	alarmBase
	masterModel *alarm_models.VirusModel
}

// BatchAdd 新增
func (crd *VirusUpload) BatchAdd(list *[]alarm_models.VirusModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []alarm_models.VirusModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		model.FileHash = strings.ToLower(model.FileHash)
		model.FileDir = path_helper.ToStandard(model.FileDir)
		model.VirusType = path_helper.ToStandard(model.VirusType)
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	return nil
}
