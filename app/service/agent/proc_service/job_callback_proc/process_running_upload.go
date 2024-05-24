package job_callback_proc

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
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

// UploadDlls 上传模块信息
func (crd *ProcessRunningUpload) UploadDlls(orgId int64, agentId, procHash string, procId, number int, list *[]proc_models.ProcessRunningDll) error {
	updates := map[string]interface{}{
		"dll_status": enum_process.ProcGetDllStatusFinish,
		"dll_number": number,
		"DllList":    json_helper.Encode(&list),
	}

	crd.masterModel.UpdateByOrgAndAgentAndHashAndProcId(orgId, agentId, procHash, procId, updates)

	return nil
}
