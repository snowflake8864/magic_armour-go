package task_service

import (
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/md5_encrypt"
	"strconv"
)

func CreateCurdFactory() *TaskCurd {
	return &TaskCurd{model_agent.CreateTaskFactory()}
}

type TaskCurd struct {
	masterModel *model_agent.TaskModel
}

// AddTasksForAgents 多个终端添加任务
func (crd *TaskCurd) AddTasksForAgents(agentIds []string, types []int, remark string) {
	for _, agentId := range agentIds {
		crd.AddTasksForAgent(agentId, types, remark)
	}
}

// AddTasksForAgent 单个终端添加任务
func (crd *TaskCurd) AddTasksForAgent(agentId string, types []int, remark string) {
	//未完成的任务
	cacheFact := CreateCacheFactory()
	unfinishedMap := helpers.ConvertIntSlice2Map(cacheFact.GetUnfinished(agentId))

	//需要缓存列表
	var caches []int
	var inserts []model_agent.TaskModel
	nowDate := helpers.GetDatetime()

	//遍历新任务并去除未完成的
	for _, val := range types {
		if _, has := unfinishedMap[val]; !has {
			caches = append(caches, val)
			inserts = append(inserts, model_agent.TaskModel{
				Id:        md5_encrypt.MD5(agentId + strconv.Itoa(val) + nowDate),
				AgentId:   agentId,
				Type:      val,
				Remark:    remark,
				Status:    enum_task.StatusStart,
				NetIp:     "",
				CreatedAt: nowDate,
			})
		}
	}

	//没有
	if len(caches) == 0 {
		return
	}

	//缓存
	cacheFact.AddUnfinished(agentId, caches)
	//入库
	crd.masterModel.InsertBatch(&inserts)
}

// CloseTasks 关闭任务
func (crd *TaskCurd) CloseTasks(agentId string, types []int, netIp string) {
	CreateCacheFactory().RemUnfinished(agentId, types)

	conditions := map[string]interface{}{
		"agent_id": agentId,
		"type":     types,
		"status":   enum_task.StatusStart,
	}
	updates := map[string]interface{}{
		"status":       enum_task.StatusFinish,
		"net_ip":       netIp,
		"completed_at": helpers.GetDatetime(),
	}
	//入库
	crd.masterModel.UpdateByCondition(conditions, updates)
}
