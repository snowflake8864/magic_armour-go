package proc_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreateSampleCurdFactory() *SampleCurd {
	return &SampleCurd{
		masterModel: proc_models.CreateProcSampleFactory(),
	}
}

type SampleCurd struct {
	masterModel *proc_models.ProcSampleModel
}

// Add 新增
func (crd *SampleCurd) Add(model *proc_models.ProcSampleModel) error {

	return nil
}

// UploadSuccess 上传成功
func (crd *SampleCurd) UploadSuccess(orgId int64, agentId string, id string, uploadModel *upload_models.UploadModel) error {
	conditions := proc_models.SearchConditionsForSample{
		OrgId:   orgId,
		AgentId: agentId,
		Id:      id,
	}

	updates := map[string]interface{}{
		"status":        enum_process.SampleStatusFinished,
		"original_name": uploadModel.OriginalName,
		"storage_path":  path_helper.ToStandard(uploadModel.StoragePath),
		"full_path":     path_helper.ToStandard(uploadModel.FullPath),
	}
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("修改失败")
	}
	return nil
}

// UploadFail 上传失败
func (crd *SampleCurd) UploadFail(orgId int64, agentId string, id string) error {
	conditions := proc_models.SearchConditionsForSample{
		OrgId:   orgId,
		AgentId: agentId,
		Id:      id,
	}

	updates := map[string]interface{}{
		"status": enum_process.SampleStatusFail,
	}
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("修改失败")
	}
	return nil
}

// DelByOrg 删除：企业
func (crd *SampleCurd) DelByOrg(orgId int64) bool {
	search := proc_models.SearchConditionsForSample{
		OrgId: orgId,
	}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *SampleCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := proc_models.SearchConditionsForSample{
		OrgId:   orgId,
		AgentId: agentId,
	}
	return crd.masterModel.DeleteByConditions(search)
}
