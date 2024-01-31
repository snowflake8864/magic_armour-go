package dictionary_service

import (
	"goskeleton/app/utils/redis_factory"
)

func CreateCacheFactory(userId int64) *dictionaryCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &dictionaryCache{
		redisClient: redCli,
		cacheKey:    "dictionary",
	}
}

type dictionaryCache struct {
	redisClient *redis_factory.RedisClient
	cacheKey    string
}

// Rebuild 重构数据库缓存
func (cache *dictionaryCache) Rebuild() bool {
	curdFact := CreateDictionaryCurdFactory()
	items := curdFact.GetAll()

	var list []string
	for _, obj := range items {
		list = append(list, obj.Id, obj.Val)
	}

	if _, err := cache.redisClient.Execute("HMSET", u.userTokenKey); err == nil {
		return true
	}
	return false
}
