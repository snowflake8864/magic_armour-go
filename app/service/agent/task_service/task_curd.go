package task_service

import (
	"goskeleton/app/model/model_agent"
)

func CreateCurdFactory() *TaskCurd {
	return &TaskCurd{model_agent.CreateTaskFactory()}
}

type TaskCurd struct {
	masterModel *model_agent.TaskModel
}
