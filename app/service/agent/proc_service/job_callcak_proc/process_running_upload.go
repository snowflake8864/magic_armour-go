package job_callcak_proc

import (
	"errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers"
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
	if helpers.StringIsEmpty(processRunningModel.CreatedAt) {
		processRunningModel.CreatedAt = helpers.DatetimeCurr()
	}

	//保存
	if !crd.masterModel.InsertItem(processRunningModel) {
		return errors.New("保存失败")
	}

	return nil
}
