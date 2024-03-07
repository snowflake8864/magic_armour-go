package task_service

import "goskeleton/app/utils/redis_factory"

func CreateCacheFactory() *TaskCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &TaskCache{
		redisClient:       redCli,
		cacheExpire:       86400,
		unfinishedPrefix:  "task_unfinished_",
		emptyNumberPrefix: "task_empty_number_",
	}
}

type TaskCache struct {
	redisClient       *redis_factory.RedisClient
	cacheExpire       int
	unfinishedPrefix  string
	emptyNumberPrefix string
}

func (crd *TaskCache) getKeyByEmptyTask(val string) string {
	return crd.emptyNumberPrefix + val
}

func (crd *TaskCache) getKeyByUnfinished(val string) string {
	return crd.unfinishedPrefix + val
}

func (crd *TaskCache) GetUnfinished(agentId string) []int {
	data, _ := crd.redisClient.SMembersByInt(crd.getKeyByUnfinished(agentId))
	return data
}

func (crd *TaskCache) AddUnfinished(agentId string, tasks []int) bool {
	isOk, _ := crd.redisClient.SAddByInt(crd.getKeyByUnfinished(agentId), tasks)
	crd.redisClient.Expire(crd.getKeyByUnfinished(agentId), crd.cacheExpire)
	return isOk
}

func (crd *TaskCache) RemUnfinished(agentId string, tasks []int) bool {
	_, err := crd.redisClient.SRemByInt(crd.getKeyByUnfinished(agentId), tasks)
	if err == nil {
		return true
	}
	return false
}

// SetHeartbeat 设置心跳,并判断是否需要追加修改配置的任务
// 【新需求 6705】1分钟没有任务下发，主动调整回连频率为2分钟一次
func (crd *TaskCache) SetHeartbeat(agentId string, hasUnfinishedTask bool) (appendConf bool) {

	return appendConf
}

// GetHeartbeatFlag 获取心跳标志
// @return bool true:普通, false:快速
func (crd *TaskCache) GetHeartbeatFlag(agentId string) (flag bool) {

	return flag
}
