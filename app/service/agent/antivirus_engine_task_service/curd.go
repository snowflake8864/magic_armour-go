package antivirus_engine_task_service

import (
	"goskeleton/app/model/agent/antivirus_engine_task_models"
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

// Add 添加
func (crd *Curd) Add(param *antivirus_engine_task_models.UpgradeTaskModel) error {

	return nil
}
