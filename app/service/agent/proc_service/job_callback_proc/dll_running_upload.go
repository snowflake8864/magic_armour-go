package job_callback_proc

import (
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateDllRunningUploadFactory() *DllRunningUpload {
	return &DllRunningUpload{
		masterModel: proc_models.CreateDllRunningFactory(),
	}
}

type DllRunningUpload struct {
	ProcessBase
	masterModel *proc_models.DllRunningModel
}

// BatchAdd 新增
func (crd *DllRunningUpload) BatchAdd(list *[]proc_models.DllRunningModel) bool {
	if len(*list) == 0 {
		return true
	}

	createdAt := date_helper.CurDatetime()
	var inserts []proc_models.DllRunningModel
	for _, model := range *list {
		if string_helper.IsEmpty(model.DllHash) {
			continue
		}
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		model.DllDir = path_helper.ToStandard(model.DllDir)
		model.ProcDir = path_helper.ToStandard(model.ProcDir)

		inserts = append(inserts, model)
	}
	if len(inserts) == 0 {
		return true
	}
	return crd.masterModel.InsertBatch(&inserts)
}
