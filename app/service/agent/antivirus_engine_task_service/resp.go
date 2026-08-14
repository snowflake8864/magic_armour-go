package antivirus_engine_task_service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreatRespFactory() *Resp {
	return &Resp{
		taskModel: antivirus_engine_task_models.CreateUpgradeTaskFactory(),
		failModel: antivirus_engine_task_models.CreateUpgradeFailFactory(),
	}
}

type Resp struct {
	taskModel *antivirus_engine_task_models.UpgradeTaskModel
	failModel *antivirus_engine_task_models.UpgradeFailModel
}

// HandleUnneeded 处理无需升级的
func (crd *Resp) HandleUnneeded(unneededModels []*antivirus_engine_task_models.UpgradeFailModel) error {
	if len(unneededModels) == 0 {
		return errors.New("len(unneededModels)==0")
	}
	currAt := date_helper.CurDatetime()
	var taskIds []int64
	for _, model := range unneededModels {
		taskIds = append(taskIds, model.TaskId)

		if model.Id == 0 {
			model.Id = variable.SnowFlake.GetId()
		}
		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = currAt
		}
	}

	//修改任务
	conditions := antivirus_engine_task_models.SearchConditions{
		Ids: taskIds,
	}
	//修改
	updates := make(map[string]interface{})
	updates["exec_unneeded_total"] = gorm.Expr("exec_unneeded_total+1")
	if !crd.taskModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	//插入错误日志
	crd.failModel.InsertBatch(unneededModels)

	return nil
}

// HandleSingleAgentProgress 处理单个终端的对应的任务进度
func (crd *Resp) HandleSingleAgentProgress(failModel *antivirus_engine_task_models.UpgradeFailModel) error {
	//是否存在任务
	existTaskModel := crd.taskModel.GetItemById(failModel.TaskId)
	if existTaskModel == nil || existTaskModel.Id == 0 {
		return errors.New("数据不存在")
	} else if existTaskModel.ExecStatus != enum_antivirus_engine_task.ExecStatusHandling {
		return fmt.Errorf("执行状态[%d]不是处理中", existTaskModel.ExecStatus)
	} else if existTaskModel.OrgId != failModel.OrgId {
		return fmt.Errorf("企业不一致[%d:%d]", existTaskModel.OrgId, failModel.OrgId)
	} else if existTaskModel.ExecBatch != failModel.ExecBatch {
		return fmt.Errorf("执行批次不一致[%d:%d]", existTaskModel.ExecBatch, failModel.ExecBatch)
	}

	//修改任务数据
	now := date_helper.CurDatetime()
	updates := map[string]interface{}{}
	//判断成功失败
	if failModel.Reason == enum_antivirus_engine_task.FailReasonEngine { // 失败数
		updates["exec_fail_total"] = gorm.Expr("exec_fail_total+1")
	} else {
		updates["exec_good_total"] = gorm.Expr("exec_good_total+1") //成功数
	}
	//（成功数+失败数+离线数+无需数+1） >= 终端数 : 完成
	completeTotal := existTaskModel.ExecGoodTotal + existTaskModel.ExecFailTotal + 1
	if completeTotal >= existTaskModel.ExecAgentTotal {
		updates["exec_status"] = enum_antivirus_engine_task.ExecStatusWaiting
		updates["exec_next_ts"] = enum_antivirus_engine_task.GetExecNextTs(existTaskModel.TaskType, existTaskModel.ExecDate,
			existTaskModel.ExecHis, existTaskModel.ExecWeekly, existTaskModel.ExecDaily)
		updates["exec_is_complete"] = enums.IsYes
		updates["exec_progress"] = 100.00
		updates["exec_complete_at"] = now
		updates["exec_completed_source"] = enums.OperateSourceAgent
		updates["exec_live_ts"] = 0
	} else {
		updates["exec_live_ts"] = date_helper.CurTimestamp()
		updates["exec_progress"] = float_helper.Div(completeTotal, existTaskModel.ExecAgentTotal, 2)
	}
	conditions := antivirus_engine_task_models.SearchConditions{
		Id: failModel.TaskId,
	}
	if !crd.taskModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	//保存错误
	if failModel.Reason == enum_antivirus_engine_task.FailReasonEngine {
		if failModel.Id == 0 {
			failModel.Id = variable.SnowFlake.GetId()
		}
		if string_helper.IsEmpty(failModel.CreatedAt) {
			failModel.CreatedAt = now
		}
		crd.failModel.InsertItem(failModel)
	}
	return nil
}

// HandleOverdue 处理逾期的任务
//
// 注意： 终端任务过期时间24小时，要保存升级任务过期时间一致
func (crd *Resp) HandleOverdue(floorTs int64) error {

	return nil
}
