package job_callback_proc

import (
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

// BatchAdd 新增
func (crd *ProcessRunningUpload) BatchAdd(list *[]proc_models.ProcessRunningModel) bool {
	if len(*list) == 0 {
		return true
	}

	createdAt := date_helper.CurDatetime()
	var inserts []proc_models.ProcessRunningModel
	for _, model := range *list {
		if string_helper.IsEmpty(model.ProcHash) {
			continue
		}
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		model.ProcDir = path_helper.ToStandard(model.ProcDir)

		inserts = append(inserts, model)
	}

	if len(inserts) == 0 {
		return true
	}
	return crd.masterModel.InsertBatch(&inserts)
}

// UploadDlls 上传模块信息
func (crd *ProcessRunningUpload) UploadDlls(orgId int64, agentId, procHash string, procId, number int, list *[]proc_models.ProcessRunningDll) error {
	updates := map[string]interface{}{
		"dll_status": enum_process.CollectDllStatusFinish,
		"dll_number": number,
		"DllList":    json_helper.Encode(&list),
	}

	crd.masterModel.UpdateByOrgAndAgentAndHashAndProcId(orgId, agentId, procHash, procId, updates)

	return nil
}
