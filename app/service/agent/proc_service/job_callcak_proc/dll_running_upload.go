package job_callcak_proc

import (
	"errors"
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

// Add 添加
func (crd *DllRunningUpload) Add(dllRunningModel *proc_models.DllRunningModel) error {
	if string_helper.IsEmpty(dllRunningModel.CreatedAt) {
		dllRunningModel.CreatedAt = date_helper.CurDatetime()
	}
	dllRunningModel.DllDir = path_helper.ToStandard(dllRunningModel.DllDir)
	dllRunningModel.ProcDir = path_helper.ToStandard(dllRunningModel.ProcDir)

	//保存
	if !crd.masterModel.InsertItem(dllRunningModel) {
		return errors.New("保存失败")
	}

	return nil
}
