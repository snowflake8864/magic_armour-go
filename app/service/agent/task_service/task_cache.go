package task_service

import "goskeleton/app/utils/redis_factory"

func CreateCacheFactory() *TaskCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &TaskCache{
		redisClient: redCli,
		cachePrefix: "task_",
	}
}

type TaskCache struct {
	redisClient *redis_factory.RedisClient
	cachePrefix string
}
