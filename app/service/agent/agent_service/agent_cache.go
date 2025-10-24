package agent_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
		authExpire:  cur_app.GetAgentTokenExpire(),
	}
}

type AgentCache struct {
	redisClient *redis_factory.RedisClient
	authPrefix  string
	authExpire  int
}

type AgentCacheLogin struct {
	OrgId        int64  `json:"org_id"`
	AgentId      string `json:"agent_id"`
	PlatformType int    `json:"platform_type"`
	OsType       int    `json:"os_type"`
	AgentIp      string `json:"agent_ip"`
	NetIp        string `json:"net_ip"`
	Os           string `json:"os"`
	LoginAt      string `json:"login_at"`
	Ver          string `json:"ver"`
}

func (crd *AgentCache) getKeyByLogin(val string) string {
	return crd.authPrefix + val
}

func (crd *AgentCache) getToken(loginInfo *AgentCacheLogin) string {
	return md5_encrypt.MD5(fmt.Sprintf("%d_%s_%s_%s", loginInfo.OrgId, loginInfo.AgentId, loginInfo.Ver, loginInfo.LoginAt))
}

// EncryptByLogin 登录信息加密并缓存
func (crd *AgentCache) EncryptByLogin(param *agent_models.AgentModel) string {
	var loginInfo AgentCacheLogin
	loginInfo.OrgId = param.OrgId
	loginInfo.AgentId = param.AgentId
	loginInfo.PlatformType = param.PlatformType
	loginInfo.OsType = param.OsType
	loginInfo.AgentIp = param.AgentIp
	loginInfo.NetIp = param.NetIp
	loginInfo.Os = param.OsName
	loginInfo.LoginAt = date_helper.CurDatetime()
	loginInfo.Ver = param.Ver

	token := crd.getToken(&loginInfo)

	_, _ = crd.redisClient.SetEx(crd.getKeyByLogin(token), json_helper.Encode(loginInfo), crd.authExpire)
	return token
}

// DecryptByLogin 登录信息解密
func (crd *AgentCache) DecryptByLogin(token string, isDelay bool) (loginInfo AgentCacheLogin, err error) {
	jsonStr := crd.redisClient.Get(crd.getKeyByLogin(token))
	if string_helper.IsEmpty(jsonStr) {
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

// ModifyByLogin 修改登录信息
func (crd *AgentCache) ModifyByLogin(cacheLogin *AgentCacheLogin, agentIp, netIp string) (err error) {
	token := crd.getToken(cacheLogin)
	jsonStr := crd.redisClient.Get(crd.getKeyByLogin(token))
	if string_helper.IsEmpty(jsonStr) {
		return errors.New("token expired")
	}

	var loginInfo AgentCacheLogin
	if err = json.Unmarshal([]byte(jsonStr), &loginInfo); err != nil {
		return errors.New("token json fail")
	}

	//agent服务去重启情况下只能修改下面信息
	loginInfo.AgentIp = agentIp
	loginInfo.NetIp = netIp

	_, _ = crd.redisClient.SetEx(crd.getKeyByLogin(token), json_helper.Encode(loginInfo), crd.authExpire)

	return nil
}
