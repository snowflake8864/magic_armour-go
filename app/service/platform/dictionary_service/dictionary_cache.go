package dictionary_service

import (
	"goskeleton/app/utils/redis_factory"
)

func CreateCacheFactory() *dictionaryCache {
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
	total := len(items)
	if total == 0 {
		return true
	} else {
		fieldsAndValues := make(map[string]string)
		for _, obj := range items {
			fieldsAndValues[obj.Id] = obj.Val
		}
		isOK, _ := cache.redisClient.HmSet(cache.cacheKey, fieldsAndValues)
		return isOK
	}
}

func (cache *dictionaryCache) GetMapByIds(ids []string) map[string]string {
	return cache.redisClient.HmGet(cache.cacheKey, ids)
}

func (cache *dictionaryCache) Lock() bool {
	isOK, _ := cache.redisClient.Lock("ggg_test", "fff", 300)
	return isOK
}
