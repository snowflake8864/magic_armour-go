package antivirus_engine_task_service

import (
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/utils/helpers/date_helper"
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

	return nil
}

// HandleOverdue 处理逾期的任务
func (crd *Resp) HandleOverdue(floorTs int64) error {

	return nil
}
