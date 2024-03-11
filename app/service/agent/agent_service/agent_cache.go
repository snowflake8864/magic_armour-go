package agent_service

import (
	"encoding/json"
	"errors"
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
	redisClient *redis_factory.RedisClient
	authPrefix  string
	authExpire  int
}

type AgentCacheLogin struct {
	AgentId      string `json:"agent_id"`
	PlatformType int    `json:"platform_type"`
	OsType       int    `json:"os_type"`
	AgentIp      string `json:"agent_ip"`
	NetIp        string `json:"net_ip"`
	LoginAt      string `json:"login_at"`
}

func (crd *AgentCache) getKeyByLogin(val string) string {
	return crd.authPrefix + val
}

// EncryptByLogin 登录信息加密并缓存
func (crd *AgentCache) EncryptByLogin(param *model_agent.AgentModel) string {
	var loginInfo AgentCacheLogin
	loginInfo.AgentId = param.AgentId
	loginInfo.PlatformType = param.PlatformType
	loginInfo.OsType = param.OsType
	loginInfo.AgentIp = param.AgentIp
	loginInfo.NetIp = param.NetIp
	loginInfo.LoginAt = helpers.GetDatetime()

	jsonStr := helpers.JsonEncode(loginInfo)
	token := md5_encrypt.MD5(jsonStr)

	crd.redisClient.SetEx(crd.getKeyByLogin(token), jsonStr, crd.authExpire)
	return token
}

// DecryptByLogin 登录信息解密
func (crd *AgentCache) DecryptByLogin(token string, isDelay bool) (loginInfo AgentCacheLogin, err error) {
	jsonStr := crd.redisClient.Get(crd.getKeyByLogin(token))
	if helpers.StringIsEmpty(jsonStr) {
		return AgentCacheLogin{}, errors.New("token expired")
	}
	err = json.Unmarshal([]byte(jsonStr), &loginInfo)
	if err != nil {
		return AgentCacheLogin{}, errors.New("token json fail")
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.getKeyByLogin(token), crd.authExpire)
	}
	return loginInfo, nil
}
