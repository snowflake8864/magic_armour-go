package task_service

import "goskeleton/app/utils/redis_factory"

func CreateCacheFactory() *TaskCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &TaskCache{
		redisClient:      redCli,
		unfinishedPrefix: "task_unfinished_",
		unfinishedExpire: 86400,
	}
}

type TaskCache struct {
	redisClient      *redis_factory.RedisClient
	unfinishedPrefix string
	unfinishedExpire int
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
	crd.redisClient.Expire(crd.getKeyByUnfinished(agentId), crd.unfinishedExpire)
	return isOk
}

func (crd *TaskCache) RemUnfinished(agentId string, tasks []int) bool {
	_, err := crd.redisClient.SRemByInt(crd.getKeyByUnfinished(agentId), tasks)
	if err == nil {
		return true
	}
	return false
}
