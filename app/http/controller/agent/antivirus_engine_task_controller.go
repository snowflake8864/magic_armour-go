package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/service/agent/antivirus_engine_task_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// AntivirusEngineController 杀毒引擎升级任务相关
type AntivirusEngineController struct {
	ControllerBase
}

// 任务类型优先级：立即 > 周期（默认最新）> 定时
//
// 注意：相同任务会被获取到多次，使用redis锁（保留24小时），24小时后批次会更新
//
// return: highestTaskModel 优先级最高的升级任务
//
// return:  engineModel 引擎信息
func (u *AntivirusEngineController) getTask(orgId int64, agentId string, deptId int64, usedVer string) (highestTaskModel *antivirus_engine_task_models.UpgradeTaskModel, engineModel *antivirus_engine_models.EngineModel) {

	//无需升级任务
	var unneededModels []*antivirus_engine_task_models.UpgradeTaskModel

	//获取任务
	conditions := antivirus_engine_task_models.SearchConditions{
		OrgId:          orgId,
		TargetAgentId:  agentId,
		TargetDeptId:   deptId,
		ExecStatusList: []int{enum_antivirus_engine_task.ExecStatusStarting, enum_antivirus_engine_task.ExecStatusHandling},
		ExecNextTsLte:  date_helper.CurTimestamp(),
	}
	//利用task_type字段, 按照立即、周期、定时排序
	if engineTaskModels := antivirus_engine_task_models.CreateUpgradeTaskFactory().GetItemsByConditions(
		conditions, 0, 0, "task_type ASC,id DESC"); len(engineTaskModels) > 0 {
		for key, taskModel := range engineTaskModels {
			if key == 0 {
				highestTaskModel = taskModel
			} else {
				unneededModels = append(unneededModels, taskModel)
			}
		}
	} else {
		return nil, nil
	}

	//判断是高优先级任务是否出处理过
	if antivirus_engine_task_service.CreateLockFactory().UpgradeTask(orgId, agentId, highestTaskModel.Id, highestTaskModel.ExecBatch) {
		//获取引擎信息
	} else {
		//
	}

	return highestTaskModel, engineModel
}

// Upgrade 升级
func (u *AntivirusEngineController) Upgrade(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeAntivirusEngineUpgrade}, ctx.ClientIP())

	//服务器地址
	err, serverAddr, _, _ := software_service.CreateDownloadFactory().GetServiceAddr()
	if err != nil {
		response.FailSimplify(ctx, fmt.Sprintf("serverAddr err: %s", err.Error()))
		return
	}

	//终端信息
	agentModel := agent_models.CreateAgentFactory().GetItemByOrgAndAgent(loginInfo.OrgId, loginInfo.AgentId,
		agent_models.CreateAgentFactory().ColumnBasic(true, false))
	if agentModel == nil || string_helper.IsEmpty(agentModel.AgentId) {
		response.FailSimplify(ctx, "agent not exists")
		return
	}

	result := map[string]interface{}{} //需初始化

	//获取到引擎任务
	highestTaskModel, engineModel := u.getTask(loginInfo.OrgId, loginInfo.AgentId, agentModel.DeptId, agentModel.AntivirusEngineVer)
	if highestTaskModel == nil || highestTaskModel.Id == 0 || engineModel == nil && engineModel.Id == 0 {
		response.FailSimplify(ctx, "No upgrade task")
		return
	}

	//返回信息
	result["taskId"] = highestTaskModel.Id
	result["batchId"] = highestTaskModel.ExecBatch
	result["ver"] = engineModel.Ver
	result["size"] = engineModel.FileSize
	result["hash"] = engineModel.FileHash
	result["download"] = serverAddr + "/" + strings.TrimLeft(engineModel.StoragePath, "/")

	response.SuccessReturnData(ctx, result)
}

// UploadProgress 上报进度
func (u *AntivirusEngineController) UploadProgress(ctx *gin.Context) {
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
