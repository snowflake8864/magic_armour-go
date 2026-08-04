package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// AntivirusEngineController 杀毒引擎升级任务相关
type AntivirusEngineController struct {
	ControllerBase
}

// Upgrade 升级
func (u *BackupController) Upgrade(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeAntivirusEngineUpgrade}, ctx.ClientIP())

	result := []map[string]interface{}{} //需初始化
	/*
		"taskId":    0,
		"batchId":    0,
		"ver":    0,
		"size":     fmt.Sprintf("%d", downloadModel.FileSize),
		"hash":     downloadModel.FileHash,
		"download": addr + "/" + strings.TrimLeft(downloadModel.StoragePath, "/"),
	*/

	response.SuccessReturnData(ctx, result)
}

// UploadProgress 上报进度
func (u *BackupController) UploadProgress(ctx *gin.Context) {
	/*loginInfo := u.getLogin(ctx)
	param := engine_upload_data.DetectTaskProgress{
		TaskId:   int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "taskId")),
		BatchId:  int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "batchId")),
		Status:   int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "status")),
		Progress: ctx.GetFloat64(consts.ValidatorPrefixAgent + "progress"),
		Reason:   ctx.GetString(consts.ValidatorPrefixAgent + "reason"),
	}
	err := detect_resp_service.CreateRespFactory().HandleDetectTaskProgress(loginInfo.AgentId, &param)
	if err != nil {
		msg := fmt.Sprintf("agentId=%s, param=%s fail: %v", loginInfo.AgentId, json_helper.Encode(param), err)
		variable.ZapLog.Debug(msg)
		response.FailSimplify(ctx, err.Error())
		return
	}*/
	response.SuccessSimplify(ctx)
}
