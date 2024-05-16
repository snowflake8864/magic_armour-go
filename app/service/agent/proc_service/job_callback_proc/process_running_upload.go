package job_callback_proc

import (
	"errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateProcessRunningUploadFactory() *ProcessRunningUpload {
	return &ProcessRunningUpload{
		masterModel: proc_models.CreateProcessRunningFactory(),
	}
}

type ProcessRunningUpload struct {
	ProcessBase
	masterModel *proc_models.ProcessRunningModel
}

// Add 添加
func (crd *ProcessRunningUpload) Add(processRunningModel *proc_models.ProcessRunningModel) error {
	if string_helper.IsEmpty(processRunningModel.CreatedAt) {
		processRunningModel.CreatedAt = date_helper.CurDatetime()
	}
	processRunningModel.ProcDir = path_helper.ToStandard(processRunningModel.ProcDir)

	//保存
	if !crd.masterModel.InsertItem(processRunningModel) {
		return errors.New("保存失败")
	}

	return nil
}
