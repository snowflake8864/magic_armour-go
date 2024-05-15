package task_service

import (
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/md5_encrypt"
	"strconv"
)

func CreateCurdFactory() *TaskCurd {
	return &TaskCurd{task_models.CreateTaskFactory()}
}

type TaskCurd struct {
	masterModel *task_models.TaskModel
}

// AddTasksForAgents 多个终端添加任务
func (crd *TaskCurd) AddTasksForAgents(orgId int64, agentIds []string, types []int, remark string) {
	for _, agentId := range agentIds {
		crd.AddTasksForAgent(orgId, agentId, types, remark)
	}
}

// AddTasksForAgent 单个终端添加任务
func (crd *TaskCurd) AddTasksForAgent(orgId int64, agentId string, types []int, remark string) {
	//未完成的任务
	cacheFact := CreateCacheFactory()
	unfinishedMap := slice_helper.SliceInt2Map(cacheFact.GetUnfinished(orgId, agentId))

	//需要缓存列表
	var caches []int
	var inserts []task_models.TaskModel
	nowDate := date_helper.CurDatetime()

	//遍历新任务并去除未完成的
	for _, val := range types {
		if _, has := unfinishedMap[val]; !has {
			caches = append(caches, val)
			inserts = append(inserts, task_models.TaskModel{
				Id:        md5_encrypt.MD5(agentId + strconv.Itoa(val) + nowDate),
				OrgId:     orgId,
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
	cacheFact.AddUnfinished(orgId, agentId, caches)
	//入库
	crd.masterModel.InsertBatch(&inserts)
}

// CloseTasks 关闭任务
func (crd *TaskCurd) CloseTasks(orgId int64, agentId string, types []int, netIp string) {
	CreateCacheFactory().RemUnfinished(orgId, agentId, types)

	conditions := task_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		Types:   types,
	}
	updates := map[string]interface{}{
		"status":       enum_task.StatusFinish,
		"net_ip":       netIp,
		"completed_at": date_helper.CurDatetime(),
	}
	//入库
	crd.masterModel.UpdateByConditions(conditions, updates)
}
