package dir_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreateDirCurdFactory() *DirCurd {
	return &DirCurd{
		masterModel:  dir_models.CreateDirFactory(),
		redisClient:  redis_factory.GetOneRedisClient(),
		searchPrefix: "agent_dir_search_",
		searchExpire: 86400,
	}
}

type DirCurd struct {
	masterModel  *dir_models.DirModel
	redisClient  *redis_factory.RedisClient
	searchPrefix string
	searchExpire int
}

func (crd *DirCurd) keyBySearch(orgId int64, agentId string) string {
	return crd.searchPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId
}

// SaveBySearch 保存当前查询路径
func (crd *DirCurd) SaveBySearch(orgId int64, agentId string, dir string, force bool) error {
	//为空默认根路径
	if string_helper.IsEmpty(dir) {
		dir = enum_tamper.RootDir
	}
	dir = path_helper.ToStandard(dir)

	//缓存路径
	cacheKey := crd.keyBySearch(orgId, agentId)
	cacheDir := crd.redisClient.Get(cacheKey)

	//任务
	task := enum_task.TypeDirStructureUpload
	existUnfinishedTask := task_service.CreateCacheFactory().ExistUnfinished(orgId, agentId, task)
	//终端存在未完成任务 && 查询路径一样
	if existUnfinishedTask && cacheDir == dir {
		return nil
	}

	//上次查询条件
	condition := dir_models.SearchConditions{
		OrgId:     orgId,
		AgentId:   agentId,
		SearchDir: dir,
		FileType:  enums.AllDefaultInt,
	}

	// 强制刷新
	if force {
		//更新查询路径
		_, _ = crd.redisClient.SetEx(cacheKey, dir, crd.searchExpire)
		//终端不存在未完成任务
		if !existUnfinishedTask {
			task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, task, "强刷")
		}

		//删除上传查询痕迹
		crd.masterModel.DeleteByConditions(condition)

	} else {
		//存在上次查询痕迹
		if crd.masterModel.Count(condition) > 0 {
			return nil
		}

		//更新查询路径
		_, _ = crd.redisClient.SetEx(cacheKey, dir, crd.searchExpire)
		//终端不存在未完成任务
		if !existUnfinishedTask {
			task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, task, "")
		}
	}

	return nil
}

// GetBySearch 获取当前查询路径
func (crd *DirCurd) GetBySearch(orgId int64, agentId string) string {
	return crd.redisClient.Get(crd.keyBySearch(orgId, agentId))
}

// BatchAdd 新增
func (crd *DirCurd) BatchAdd(orgId int64, agentId string, searchDir string, params *[]dir_models.DirModel) error {
	//删除缓存路径
	_, _ = crd.redisClient.Del(crd.keyBySearch(orgId, agentId))

	//先删除原来查询路径
	condition := dir_models.SearchConditions{
		OrgId:     orgId,
		AgentId:   agentId,
		SearchDir: searchDir,
		FileType:  enums.AllDefaultInt,
	}
	crd.masterModel.DeleteByConditions(condition)

	//批量参数
	crd.masterModel.InsertBatch(params)

	return nil

}
