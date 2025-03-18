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
)

func CreateCurdFactory() *BackupCurd {
	return &BackupCurd{masterModel: backup_models.CreateBackupFactory()}
}

type BackupCurd struct {
	masterModel *backup_models.BackupModel
}

// 新增任务
func (crd *BackupCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeBlackmailPolicyDownload, remark)
}

// Add 新增
func (crd *BackupCurd) Add(param *backup_models.BackupModel) error {
	now := date_helper.CurDatetime()

	param.Id = fmt.Sprintf("%d", variable.SnowFlake.GetId())

	param.SourcePath = path_helper.ToStandard(param.SourcePath)
	param.StorageDir = path_helper.ToStandard(param.StorageDir)
	param.BackupStatus = enum_backup.BackupStatusDefault
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

	return nil
}
