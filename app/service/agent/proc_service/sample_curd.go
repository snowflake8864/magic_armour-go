package proc_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
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
	hasModel := crd.masterModel.GetItemByOrgAndAgentAndHash(model.OrgId, model.AgentId, model.Id)
	if hasModel != nil {
		if hasModel.Status == enum_process.SampleStatusFinished {
			return errors.New("数据已存在且提取样本成功")
		}
		if hasModel.Status == enum_process.SampleStatusGetting {
			return errors.New("数据已存在且正在提取中")
		}
	}

	model.Status = enum_process.SampleStatusDefault
	model.CreatedAt = date_helper.CurDatetime()
	model.UpdatedAt = model.CreatedAt

	if !crd.masterModel.InsertItem(model) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//
	task_service.CreateCurdFactory().AddTaskForAgent(model.OrgId, model.AgentId, enum_task.TypeProcessSampleUpload, "")

	return nil
}

// UploadSuccess 上传成功
func (crd *SampleCurd) UploadSuccess(orgId int64, agentId string, id string, uploadModel *upload_models.UploadModel) error {
	conditions := proc_models.SearchConditions{
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
	conditions := proc_models.SearchConditions{
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
	search := proc_models.SearchConditions{
		OrgId: orgId,
	}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *SampleCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	return crd.masterModel.DeleteByConditions(search)
}

// DeleteById 通过hash删除
func (crd *SampleCurd) DeleteById(orgId int64, agentId string, id string) bool {
	search := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		Id:      id,
	}
	return crd.masterModel.DeleteByConditions(search)
}
