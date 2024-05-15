package job_callcak_proc

import (
	"errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers"
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

// Add 添加
func (crd *DllRunningUpload) Add(dllRunningModel *proc_models.DllRunningModel) error {
	if helpers.StringIsEmpty(dllRunningModel.CreatedAt) {
		dllRunningModel.CreatedAt = helpers.DatetimeCurr()
	}

	//保存
	if !crd.masterModel.InsertItem(dllRunningModel) {
		return errors.New("保存失败")
	}

	return nil
}
