package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_backup"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/backup_models"
	"goskeleton/app/service/agent/backup_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/response"
)

// BackupController 备份回滚相关
type BackupController struct {
	ControllerBase
}

// GetBackups 获取备份策略
func (u *BackupController) GetBackups(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeSystemBackup}, ctx.ClientIP())

	result := []map[string]interface{}{} //需初始化
	var ids []string                     //需修改状态

	condition := backup_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		AgentId:      loginInfo.AgentId,
		BackupStatus: enum_backup.BackupStatusDefault,
	}
	models := backup_models.CreateBackupFactory().GetListByConditions(condition)
	if models != nil && len(models) > 0 {
		for _, model := range models {
			var sourcePath, storageDir string

			if loginInfo.PlatformType == enum_software.PlatformTypeWindows {
				sourcePath = path_helper.ToWindows(model.SourcePath)
				storageDir = path_helper.ToWindows(model.StorageDir)
			} else {
				sourcePath = path_helper.ToStandard(model.SourcePath)
				storageDir = path_helper.ToStandard(model.StorageDir)
			}
			result = append(result, map[string]interface{}{
				"id":          model.Id,
				"source_path": sourcePath,
				"storage_dir": storageDir,
			})

			ids = append(ids, model.Id)
		}
	}

	if len(ids) > 0 {
		_ = backup_service.CreateCurdFactory().BackupRunning(loginInfo.OrgId, loginInfo.AgentId, ids)
	}

	response.SuccessReturnData(ctx, result)
}

// UploadBackup 上传备份结果
func (u *BackupController) getResp(ctx *gin.Context, isBackup bool) *backup_models.BackupModel {
	loginInfo := u.getLogin(ctx)

	param := backup_models.BackupModel{}
	param.OrgId = loginInfo.OrgId
	param.AgentId = loginInfo.AgentId
	param.Id = ctx.GetString(consts.ValidatorPrefixAgent + "id")
	state := int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "state"))
	if isBackup {
		param.FileName = ctx.GetString(consts.ValidatorPrefixAgent + "file_name")
		param.FileSize = ctx.GetString(consts.ValidatorPrefixAgent + "file_size")
		param.BackupFailReason = ctx.GetString(consts.ValidatorPrefixAgent + "fail_reason")
		if state == 0 {
			param.BackupStatus = enum_backup.BackupStatusFail
		} else {
			param.BackupStatus = enum_backup.BackupStatusSuccess
		}
	} else {
		param.RollbackFailReason = ctx.GetString(consts.ValidatorPrefixAgent + "fail_reason")
		if state == 0 {
			param.RollbackStatus = enum_backup.RollbackStatusFail
		} else {
			param.RollbackStatus = enum_backup.RollbackStatusSuccess
		}
	}

	return &param
}

// UploadBackup 上传备份结果
func (u *BackupController) UploadBackup(ctx *gin.Context) {
	param := u.getResp(ctx, true)

	if err := backup_service.CreateCurdFactory().BackupComplete(param); err != nil {
		response.FailData(ctx, "BackupComplete fail", err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// GetRollbacks 获取回滚策略
func (u *BackupController) GetRollbacks(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeSystemRollback}, ctx.ClientIP())

	result := []map[string]interface{}{} //需初始化

	condition := backup_models.SearchConditions{
		OrgId:          loginInfo.OrgId,
		AgentId:        loginInfo.AgentId,
		BackupStatus:   enum_backup.BackupStatusSuccess,
		RollbackStatus: enum_backup.RollbackStatusWaiting,
	}
	models := backup_models.CreateBackupFactory().GetListByConditions(condition)
	if models != nil && len(models) > 0 {
		for _, model := range models {
			var sourcePath, storageDir string

			if loginInfo.PlatformType == enum_software.PlatformTypeWindows {
				sourcePath = path_helper.ToWindows(model.SourcePath)
				storageDir = path_helper.ToWindows(model.StorageDir)
			} else {
				sourcePath = path_helper.ToStandard(model.SourcePath)
				storageDir = path_helper.ToStandard(model.StorageDir)
			}
			result = append(result, map[string]interface{}{
				"id":          model.Id,
				"source_path": sourcePath,
				"storage_dir": storageDir,
				"file_name":   model.FileName,
				"file_size":   model.FileSize,
			})
		}
	}

	response.SuccessReturnData(ctx, result)
}

// UploadRollback 上传回滚结果
func (u *BackupController) UploadRollback(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}
