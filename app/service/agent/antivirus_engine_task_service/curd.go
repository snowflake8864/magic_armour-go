package antivirus_engine_task_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreatCurdFactory() *Curd {
	return &Curd{
		taskModel: antivirus_engine_task_models.CreateUpgradeTaskFactory(),
		failModel: antivirus_engine_task_models.CreateUpgradeFailFactory(),
	}
}

type Curd struct {
	taskModel *antivirus_engine_task_models.UpgradeTaskModel
	failModel *antivirus_engine_task_models.UpgradeFailModel
}

func (crd *Curd) getDeptIds(targetIds string) []int64 {
	var deptIds []int64
	for _, val := range strings.Split(targetIds, ",") {
		deptIds = append(deptIds, string_helper.ToInt64(val))
	}
	return deptIds
}

func (crd *Curd) getAgentTotal(param *antivirus_engine_task_models.UpgradeTaskModel) int {

	if param.TargetType == enum_antivirus_engine_task.TargetTypeAllAgent {
		return int(agent_models.CreateAgentFactory().Count(agent_models.SearchConditions{OrgId: param.OrgId}))
	} else if param.TargetType == enum_antivirus_engine_task.TargetTypeSingleAgent {
		return 1
	} else if param.TargetType == enum_antivirus_engine_task.TargetTypeDepartments {
		return int(agent_models.CreateAgentFactory().Count(agent_models.SearchConditions{
			OrgId:   param.OrgId,
			DeptIds: crd.getDeptIds(param.TargetIds),
		}))
	} else {
		return 0
	}
}

func (crd *Curd) addTask(param *antivirus_engine_task_models.UpgradeTaskModel, taskType int, remark string) {

	var agentIdList []string
	if param.TargetType == enum_antivirus_engine_task.TargetTypeAllAgent {
		agentIdList = agent_models.CreateAgentFactory().GetIdsByByConditions(agent_models.SearchConditions{OrgId: param.OrgId})
	} else if param.TargetType == enum_antivirus_engine_task.TargetTypeDepartments {
		agentIdList = agent_models.CreateAgentFactory().GetIdsByByConditions(agent_models.SearchConditions{
			OrgId:   param.OrgId,
			DeptIds: crd.getDeptIds(param.TargetIds),
		})
	} else if param.TargetType == enum_antivirus_engine_task.TargetTypeSingleAgent {
		agentIdList = append(agentIdList, param.TargetIds)
	}

	if len(agentIdList) > 0 {
		//终端任务指令
		task_service.CreateCurdFactory().AddTaskForAgents(param.OrgId, agentIdList, taskType, remark)
	}
}

// Add 添加
func (crd *Curd) Add(param *antivirus_engine_task_models.UpgradeTaskModel) error {
	currAt := date_helper.CurDatetime()

	param.ExecStatus = enum_antivirus_engine_task.ExecStatusWaiting
	param.ExecNextTs = enum_antivirus_engine_task.GetExecNextTs(param.TaskType, param.ExecDate,
		param.ExecHis, param.ExecWeekly, param.ExecDaily)
	param.ExecIsComplete = enums.IsNot
	param.ExecProgress = 0.00
	param.CreatedAt = currAt
	param.UpdatedAt = currAt

	if param.TaskType == enum_antivirus_engine_task.TaskTypeManualInstantly {
		param.ExecAgentTotal = crd.getAgentTotal(param)
		param.ExecStatus = enum_antivirus_engine_task.ExecStatusStarting
		param.ExecBatch = 1
		param.ExecLiveTs = date_helper.CurTimestamp()
		param.ExecStartAt = currAt
	}
	if !crd.taskModel.InsertItem(param) {
		return errors.New("添加失败")
	}

	if param.TaskType == enum_antivirus_engine_task.TaskTypeManualInstantly {
		crd.addTask(param, enum_task.TypeAntivirusEngineUpgrade, "添加")
	}
	return nil
}
