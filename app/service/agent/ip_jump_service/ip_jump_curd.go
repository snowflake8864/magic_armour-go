package ip_jump_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums/enum_ip_jump"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreatCurdFactory(orgId int64, jumpConfig *enum_ip_jump.IpJumpConfig) *IpJumpCurd {
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
	poolInitMax  int                        //ip次初始化最大长度
	ipPoolPrefix string                     // IP次前缀
	jumpConfig   *enum_ip_jump.IpJumpConfig //跳变配置
	orgId        int64
}

func (crd *IpJumpCurd) keyByIpPool() string {
	return crd.ipPoolPrefix + int_helper.Int64ToStr(crd.orgId)
}

// InitIpPool 初始化IP池
func (crd *IpJumpCurd) InitIpPool() error {
	//获取终端列表
	agentFact := agent_models.CreateAgentFactory()
	agentModels := agentFact.GetItemsByOrg(crd.orgId, agentFact.ColumnIp())
	agentSize := len(agentModels)
	if agentSize == 0 {
		crd.poolInitMax = 32
	} else if agentSize < crd.poolInitMax {
		crd.poolInitMax = agentSize
	}
	fmt.Println("poolInitMax: ", crd.poolInitMax)
	//存在ip集合
	existIps := map[string]struct{}{
		crd.jumpConfig.Gateway: {},
	}
	if agentSize > 0 {
		for _, agentModel := range agentModels {
			existIps[agentModel.NetIp] = struct{}{}
		}
	}

	//根据策略生成IP池
	var ipList []string
	iPv4Range := ip_helper.IPv4Range{IpWithMask: fmt.Sprintf("%s/%d", crd.jumpConfig.Ip, crd.jumpConfig.Mask)}
	if err, _, _, _, _, firstUsableNumber, lastUsableNumber := iPv4Range.CalcRange(); err != nil {
		return err
	} else {
		newNumber := firstUsableNumber
		for newNumber <= lastUsableNumber {
			//转成ip
			newIp := ip_helper.Uint32ToIPv4(newNumber)

			if _, isOK := existIps[newIp]; !isOK {
				ipList = append(ipList, newIp)
				//IP池到达预期数目
				if len(ipList) >= crd.poolInitMax {
					break
				}
			}
		}
	}

	//先删除
	cacheKey := crd.keyByIpPool()
	_, _ = crd.redisClient.Del(cacheKey)
	//在缓存
	if isOK, err := crd.redisClient.RPush(cacheKey, ipList); err != nil || !isOK {
		return errors.New("生成IP池失败")
	}

	return nil
}
