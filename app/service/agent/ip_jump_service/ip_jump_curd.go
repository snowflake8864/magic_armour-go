package ip_jump_service

import (
	"goskeleton/app/http/validator/organization/ip_jump"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreatCurdFactory(orgId int64, jumpConfig *ip_jump.GlobalConfig) *IpJumpCurd {
	return &IpJumpCurd{
		redisClient:  redis_factory.GetOneRedisClient(),
		poolInitMax:  256,
		ipPoolPrefix: "ip_jump_",
		jumpConfig:   jumpConfig,
		orgId:        orgId,
	}
}

type IpJumpCurd struct {
	redisClient  *redis_factory.RedisClient
	poolInitMax  int32                 //ip次初始化最大长度
	ipPoolPrefix string                // IP次前缀
	jumpConfig   *ip_jump.GlobalConfig //跳变配置
	orgId        int64
}

func (crd *IpJumpCurd) keyByIpPool() string {
	return crd.ipPoolPrefix + int_helper.Int64ToStr(crd.orgId)
}

// InitIpPool 初始化IP池
func (crd *IpJumpCurd) InitIpPool() error {
	//获取
	return nil
}
