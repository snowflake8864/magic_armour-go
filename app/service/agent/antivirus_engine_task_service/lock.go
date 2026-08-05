package antivirus_engine_task_service

import (
	"fmt"
	"goskeleton/app/utils/redis_factory"
)

func CreateLockFactory() *Lock {
	return &Lock{
		redisClient:       redis_factory.GetOneRedisClient(),
		upgradeTaskPrefix: "lock_ant_eng_",
		upgradeTaskExpire: 86405,
	}
}

type Lock struct {
	redisClient       *redis_factory.RedisClient
	upgradeTaskPrefix string
	upgradeTaskExpire int
}

// UpgradeTask 锁：升级任务
func (crd *Lock) UpgradeTask(orgId int64, agentId string, taskId int64, batchId int) bool {
	key := crd.upgradeTaskPrefix + fmt.Sprintf("%d_%s_%d_%d", orgId, agentId, taskId, batchId)
	if isOk, _ := crd.redisClient.Lock(key, "1", crd.upgradeTaskExpire); isOk {
		return true
	} else {
		return false
	}
}
