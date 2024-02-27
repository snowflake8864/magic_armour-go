package agent_service

import (
	"encoding/json"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
)

func CreateCacheFactory() *AgentCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &AgentCache{
		redisClient: redCli,
		authPrefix:  "agent_auth_",
		authExpire:  7200,
	}
}

type AgentCache struct {
	redisClient  *redis_factory.RedisClient
	authPrefix   string
	authExpire   int
	AgentId      string `json:"agent_id"`
	PlatformType int    `json:"platform_type"`
	OsType       int    `json:"os_type"`
	CreatedAt    string `json:"created_at"`
}

func (crd *AgentCache) getKeyByLogin(val string) string {
	return crd.authPrefix + val
}

// EncryptByLogin 登录信息加密并缓存
func (crd *AgentCache) EncryptByLogin(param *model_agent.AgentModel) string {
	crd.AgentId = param.AgentId
	crd.PlatformType = param.PlatformType
	crd.OsType = param.OsType
	crd.CreatedAt = helpers.GetDatetime()

	jsonByte, _ := json.Marshal(crd)
	jsonStr := string(jsonByte)
	token := md5_encrypt.MD5(jsonStr)

	crd.redisClient.SetEx(crd.getKeyByLogin(token), jsonStr, crd.authExpire)
	return token
}

// DecryptByLogin 登录信息解密
func (crd *AgentCache) DecryptByLogin(token string, isDelay bool) *AgentCache {
	jsonStr := crd.redisClient.Get(crd.getKeyByLogin(token))
	if helpers.IsEmptyString(jsonStr) {
		return nil
	}
	err := json.Unmarshal([]byte(jsonStr), crd)
	if err != nil {
		return nil
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.getKeyByLogin(token), crd.authExpire)
	}

	return crd
}
