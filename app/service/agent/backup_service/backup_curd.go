package backup_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums/enum_backup"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/backup_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *BackupCurd {
	return &BackupCurd{masterModel: backup_models.CreateBackupFactory()}
}

type BackupCurd struct {
	masterModel *backup_models.BackupModel
}

// Add 新增
func (crd *BackupCurd) Add(param *backup_models.BackupModel) error {
	now := date_helper.CurDatetime()

	param.Id = fmt.Sprintf("%d", variable.SnowFlake.GetId())

	param.SourcePath = path_helper.ToStandard(param.SourcePath)
	param.StorageDir = path_helper.ToStandard(param.StorageDir)
	param.BackupStatus = enum_backup.BackupStatusDefault
	param.RollbackStatus = enum_backup.RollbackStatusDefault
	param.BackupStartAt = now
	param.CreatedAt = now

	//验证数据
	conditions := backup_models.SearchConditions{
		OrgId:            param.OrgId,
		AgentId:          param.AgentId,
		SourcePath:       param.SourcePath,
		BackupStatusList: []int{enum_backup.BackupStatusDefault, enum_backup.BackupStatusRunning},
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("有相同路径在备份中，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	///新增任务
	task_service.CreateCurdFactory().AddTaskForAgent(param.OrgId, param.AgentId, enum_task.TypeSystemBackup, "")

	return nil
}

// BackupRunning 备份执行中
func (crd *BackupCurd) BackupRunning(orgId int64, agentId string, ids []string) error {
	conditions := backup_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		Ids:     ids,
	}
	//修改
	updates := make(map[string]interface{})
	updates["backup_status"] = enum_backup.BackupStatusRunning
	updates["backup_start_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}

// BackupComplete 备份完成
func (crd *BackupCurd) BackupComplete(param *backup_models.BackupModel) error {
	conditions := backup_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Id:      param.Id,
	}

	//修改
	updates := make(map[string]interface{})
	updates["storage_dir"] = param.StorageDir
	updates["file_name"] = param.FileName
	updates["file_size"] = param.FileSize
	updates["backup_status"] = param.BackupStatus
	updates["backup_fail_reason"] = param.BackupFailReason
	updates["backup_param"] = param.BackupParam
	updates["backup_end_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}

// RollbackStart 回滚开始
func (crd *BackupCurd) RollbackStart(param *backup_models.BackupModel) error {
	existModel := crd.masterModel.GetItemByOrgAndAgentAndHash(param.OrgId, param.AgentId, param.Id)
	if existModel == nil || string_helper.IsEmpty(existModel.Id) {
		return errors.New("数据不存在")
	}
	if existModel.BackupStatus != enum_backup.BackupStatusSuccess {
		return errors.New("备份没有成功， 不能进行回滚")
	}
	if existModel.RollbackStatus == enum_backup.RollbackStatusWaiting {
		return errors.New("等待客户端执行，不能重复进行")
	}
	if existModel.RollbackStatus == enum_backup.RollbackStatusRunning {
		return errors.New("正在回滚中，不能重复进行")
	}

	//修改
	updates := make(map[string]interface{})
	updates["rollback_status"] = enum_backup.RollbackStatusWaiting
	updates["rollback_id"] = param.RollbackId
	updates["rollback_name"] = param.RollbackName
	updates["rollback_start_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(backup_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Id:      param.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	task_service.CreateCurdFactory().AddTaskForAgent(param.OrgId, param.AgentId, enum_task.TypeSystemRollback, "")
	return nil
}

// RollbackRunning 回滚执行中
func (crd *BackupCurd) RollbackRunning(orgId int64, agentId string, ids []string) error {
	conditions := backup_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		Ids:     ids,
	}
	//修改
	updates := make(map[string]interface{})
	updates["rollback_status"] = enum_backup.RollbackStatusRunning
	updates["rollback_start_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}

// RollbackComplete 回滚完成
func (crd *BackupCurd) RollbackComplete(param *backup_models.BackupModel) error {
	conditions := backup_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Id:      param.Id,
	}

	//修改
	updates := make(map[string]interface{})
	updates["rollback_status"] = param.RollbackStatus
	updates["rollback_fail_reason"] = param.RollbackFailReason
	updates["rollback_end_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(conditions, updates); !isOk {
		return errors.New("保存失败")
	}
	return nil
}
