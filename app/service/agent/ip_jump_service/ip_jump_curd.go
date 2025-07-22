package ip_jump_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums/enum_ip_jump"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreatCurdFactory(orgId int64, jumpConfig *enum_ip_jump.IpJumpConfig) *IpJumpCurd {
	return &IpJumpCurd{
		redisClient:  redis_factory.GetOneRedisClient(),
		ipPoolPrefix: "ip_jump_",
		jumpConfig:   jumpConfig,
		orgId:        orgId,
	}
}

type IpJumpCurd struct {
	redisClient  *redis_factory.RedisClient
	ipPoolPrefix string                     // IP次前缀
	jumpConfig   *enum_ip_jump.IpJumpConfig //跳变配置
	orgId        int64
}

func (crd *IpJumpCurd) keyByIpPool() string {
	return crd.ipPoolPrefix + int_helper.Int64ToStr(crd.orgId)
}

func (crd *IpJumpCurd) getIpWithMask() string {
	return fmt.Sprintf("%s/%d", crd.jumpConfig.Ip, crd.jumpConfig.Mask)
}

// InitIpPool 初始化IP池
func (crd *IpJumpCurd) InitIpPool() error {
	//获取终端列表
	agentFact := agent_models.CreateAgentFactory()
	agentModels := agentFact.GetItemsByOrg(crd.orgId, agentFact.ColumnIp())
	agentSize := len(agentModels)
	//存在ip集合, 默认网关不能下发
	existIps := map[string]struct{}{
		crd.jumpConfig.Gateway: {},
	}
	if agentSize > 0 {
		//终端ip不下发
		for _, agentModel := range agentModels {
			existIps[agentModel.NetIp] = struct{}{}
		}
	}

	//根据策略生成IP池
	var ipList []string
	iPv4Range := ip_helper.IPv4Range{IpWithMask: crd.getIpWithMask()}
	if err, _, _, _, _, firstUsableNumber, lastUsableNumber := iPv4Range.CalcRange(); err != nil {
		return err
	} else {
		//IP池大小
		poolSize := agentSize + 16
		//从小到大结束值
		var ascEndNumber uint32
		for i := firstUsableNumber; i <= lastUsableNumber; i++ {
			ascEndNumber = i
			//转成ip
			newIp := ip_helper.Uint32ToIPv4(i)

			if _, isOK := existIps[newIp]; !isOK {
				ipList = append(ipList, newIp)
				//IP池到达预期数目
				if len(ipList) >= poolSize {
					break
				}
			}
		}

		//跳变方式：随机
		if crd.jumpConfig.Way == enum_ip_jump.WayRandom && ascEndNumber < lastUsableNumber {
			var n uint32
			for n = 0; n < 1024; n++ {
				//转成ip
				newIp := ip_helper.Uint32ToIPv4(int_helper.Unt32Random(ascEndNumber, lastUsableNumber))
				if _, isOK := existIps[newIp]; !isOK {
					ipList = append(ipList, newIp)
				}
				if n+ascEndNumber >= lastUsableNumber {
					break
				}
			}
		}
	}
	//跳变方式：随机
	if crd.jumpConfig.Way == enum_ip_jump.WayRandom {
		slice_helper.ShuffleStrings(ipList)
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

// NewestIp 新IP
func (crd *IpJumpCurd) NewestIp(agentIp string) (error, string) {
	cacheKey := crd.keyByIpPool()
	//从IP池取出
	ipStr, err := crd.redisClient.LPopToStr(cacheKey)
	if err != nil {
		return errors.New("get newest ip fail"), ""
	}

	//终端IP符合新IP池规则
	if ip_helper.IsIPInRange(agentIp, crd.getIpWithMask()) {
		//终端IP放入IP池末尾, 防止掩码比较大情况下ip池不够用
		_, _ = crd.redisClient.RPush(cacheKey, agentIp)
	}

	return nil, ipStr
}
