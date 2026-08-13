package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/model/platform/antivirus_engine_models"
	"goskeleton/app/service/agent/agent_service"
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

func (u *AntivirusEngineController) formatFailModel(taskModel *antivirus_engine_task_models.UpgradeTaskModel, agentId string, deptId int64,
	remark, currAt string) *antivirus_engine_task_models.UpgradeFailModel {

	return &antivirus_engine_task_models.UpgradeFailModel{
		OrgId:     taskModel.OrgId,
		TaskId:    taskModel.Id,
		ExecBatch: taskModel.ExecBatch,
		AgentId:   agentId,
		DeptId:    deptId,
		CreatedAt: currAt,
		Id:        variable.SnowFlake.GetId(),
		Reason:    enum_antivirus_engine_task.FailReasonUnneeded,
		Remark:    remark,
	}
}

// 终端获取到优先级最高的升级任务，并处理其它任务
//
// 任务类型优先级：立即 > 周期（默认最新）> 定时
//
// 注意：相同任务会被获取到多次，使用redis锁（保留24小时），24小时后批次会更新
//
// return: err 错误描述
//
// return: highestTaskModel 优先级最高的升级任务
//
// return:  engineModel 引擎信息
func (u *AntivirusEngineController) getTask(loginInfo *agent_service.AgentCacheLogin, deptId int64, usedVer string) (err error,
	highestTaskModel *antivirus_engine_task_models.UpgradeTaskModel, engineModel *antivirus_engine_models.EngineModel) {

	//无需升级任务
	var unneededModels []*antivirus_engine_task_models.UpgradeFailModel

	//获取任务
	conditions := antivirus_engine_task_models.SearchConditions{
		OrgId:          loginInfo.OrgId,
		TargetAgentId:  loginInfo.AgentId,
		TargetDeptId:   deptId,
		ExecStatusList: []int{enum_antivirus_engine_task.ExecStatusStarting, enum_antivirus_engine_task.ExecStatusHandling},
		ExecNextTsLte:  date_helper.CurTimestamp(),
	}
	currAt := date_helper.CurDatetime()
	//利用task_type字段, 按照立即、周期、定时排序
	if engineTaskModels := antivirus_engine_task_models.CreateUpgradeTaskFactory().GetItemsByConditions(
		conditions, 0, 0, "task_type ASC,id DESC"); len(engineTaskModels) > 0 {
		for _, taskModel := range engineTaskModels {
			//升级任务批次短时间没有重复执行过
			if antivirus_engine_task_service.CreateLockFactory().UpgradeTask(loginInfo.OrgId, loginInfo.AgentId, highestTaskModel.Id, highestTaskModel.ExecBatch) {
				//没有优先级最高任务
				if highestTaskModel == nil {
					//默认第一个优先级最高
					highestTaskModel = taskModel
				} else {
					//追加无需升级任务列表
					unneededModels = append(unneededModels, u.formatFailModel(taskModel, loginInfo.AgentId, deptId,
						fmt.Sprintf("优先级小于任务[%d]、批次[%d]", highestTaskModel.Id, highestTaskModel.ExecBatch), currAt))
				}
			}
		}
		if highestTaskModel == nil {
			return fmt.Errorf("tasks repeatedly executed"), nil, nil
		}
	} else {
		return fmt.Errorf("no eligible upgrade task"), nil, nil
	}

	//获取引擎信息
	isPublish := enums.IsYes
	engineConditions := antivirus_engine_models.SearchConditions{
		PlatformType: enum_software.PlatformTypeWindows,
		IsPublish:    &isPublish,
	}
	if loginInfo.PlatformType == enum_software.PlatformTypeWindows && highestTaskModel.EngineWindowsId > 0 {
		engineConditions.Id = highestTaskModel.EngineWindowsId
	} else if loginInfo.PlatformType == enum_software.PlatformTypeLinux && highestTaskModel.EngineLinuxId > 0 {
		engineConditions.Id = highestTaskModel.EngineLinuxId
	} else {
		isNewest := enums.IsYes
		engineConditions.IsNewest = &isNewest
	}
	engineModel = antivirus_engine_models.CreateEngineFactory().GetItemByConditions(engineConditions)
	if engineModel == nil || engineModel.Id == 0 {
		//记录错误原因
		err = fmt.Errorf("task[%d], batch[%d] no eligible upgrade", highestTaskModel.Id, highestTaskModel.ExecBatch)

		//追加无需升级任务列表
		unneededModels = append(unneededModels, u.formatFailModel(highestTaskModel, loginInfo.AgentId, deptId,
			"no engine that can be upgraded", currAt))

	} else if engineModel.Ver == usedVer {
		//记录错误原因
		err = fmt.Errorf("task[%d], batch[%d] same version[%s]", highestTaskModel.Id, highestTaskModel.ExecBatch, usedVer)

		//追加无需升级任务列表
		unneededModels = append(unneededModels, u.formatFailModel(highestTaskModel, loginInfo.AgentId, deptId,
			fmt.Sprintf("same version[%s]", usedVer), currAt))
	}

	//处理无需任务列表
	if len(unneededModels) > 0 {

	}

	return err, highestTaskModel, engineModel
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
	if err, highestTaskModel, engineModel := u.getTask(loginInfo, agentModel.DeptId, agentModel.AntivirusEngineVer); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else {
		//返回信息
		result["taskId"] = highestTaskModel.Id
		result["batchId"] = highestTaskModel.ExecBatch
		result["ver"] = engineModel.Ver
		result["size"] = engineModel.FileSize
		result["hash"] = engineModel.FileHash
		result["download"] = serverAddr + "/" + strings.TrimLeft(engineModel.StoragePath, "/")
	}
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
