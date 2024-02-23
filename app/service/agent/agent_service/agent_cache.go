package agent_service

import "goskeleton/app/utils/redis_factory"

func CreateCacheFactory() *AgentCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &AgentCache{
		redisClient: redCli,
		cachePrefix: "agent_",
	}
}

type AgentCache struct {
	redisClient *redis_factory.RedisClient
	cachePrefix string
}
