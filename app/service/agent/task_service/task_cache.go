package task_service

import (
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
	"strconv"
)

func CreateCacheFactory() *TaskCache {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &TaskCache{
		redisClient:          redCli,
		orgUninstallPrefix:   "task_org_uninstall_",
		orgUninstallExpire:   86400 * 3,
		unfinishedPrefix:     "task_unfinished_",
		unfinishedExpire:     86400,
		emptyNumberPrefix:    "task_empty_number_",
		emptyNumberThreshold: 3,
	}
}

type TaskCache struct {
	redisClient          *redis_factory.RedisClient
	orgUninstallPrefix   string
	orgUninstallExpire   int
	unfinishedPrefix     string
	unfinishedExpire     int
	emptyNumberPrefix    string
	emptyNumberThreshold int
}

func (crd *TaskCache) keyByEmptyTask(orgId int64, agentId string) string {
	return crd.emptyNumberPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId
}

func (crd *TaskCache) keyByUnfinished(orgId int64, agentId string) string {
	return crd.unfinishedPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId
}

func (crd *TaskCache) keyByOrgUninstall(orgId int64) string {
	return crd.orgUninstallPrefix + int_helper.Int64ToStr(orgId)
}

// AddOrgUninstall 企业级卸载任务: 新增
func (crd *TaskCache) AddOrgUninstall(orgId int64) bool {
	isOk, _ := crd.redisClient.SetEx(crd.keyByOrgUninstall(orgId), date_helper.CurDatetime(), crd.orgUninstallExpire)
	return isOk
}

// ExistsOrgUninstall 企业级卸载任务: 判断
func (crd *TaskCache) ExistsOrgUninstall(orgId int64) bool {
	isOk, _ := crd.redisClient.Exists(crd.keyByOrgUninstall(orgId))
	return isOk
}

// DeleteOrgUninstall 企业级卸载任务: 删除
func (crd *TaskCache) DeleteOrgUninstall(orgId int64) bool {
	_, err := crd.redisClient.Del(crd.keyByOrgUninstall(orgId))
	if err == nil {
		return true
	}
	return false
}

// GetUnfinished 终端未完成任务: 获取
func (crd *TaskCache) GetUnfinished(orgId int64, agentId string) []int {
	data, _ := crd.redisClient.SMembersByInt(crd.keyByUnfinished(orgId, agentId))
	return data
}

// AddUnfinished 终端未完成任务:新增
func (crd *TaskCache) AddUnfinished(orgId int64, agentId string, tasks []int) bool {
	isOk, _ := crd.redisClient.SAddByInt(crd.keyByUnfinished(orgId, agentId), tasks)
	crd.redisClient.Expire(crd.keyByUnfinished(orgId, agentId), crd.unfinishedExpire)
	return isOk
}

// RemUnfinished 终端未完成任务:移除
func (crd *TaskCache) RemUnfinished(orgId int64, agentId string, tasks []int) bool {
	_, err := crd.redisClient.SRemByInt(crd.keyByUnfinished(orgId, agentId), tasks)
	if err == nil {
		return true
	}
	return false
}

// ExistUnfinished 终端未完成任务:存在某个任务
func (crd *TaskCache) ExistUnfinished(orgId int64, agentId string, task int) bool {
	exist, err := crd.redisClient.SisMember(crd.keyByUnfinished(orgId, agentId), task)
	if err != nil {
		return false
	}
	return exist
}

// DeleteUnfinished 终端未完成任务： 删除
func (crd *TaskCache) DeleteUnfinished(orgId int64, agentId string) bool {
	_, err := crd.redisClient.Del(crd.keyByUnfinished(orgId, agentId))
	if err != nil {
		return false
	}
	return true
}

// SetHeartbeat 设置心跳,并判断是否需要追加修改配置的任务
// 【新需求 6705】1分钟没有任务下发，主动调整回连频率为2分钟一次
func (crd *TaskCache) SetHeartbeat(orgId int64, agentId string, hasUnfinishedTask bool) (appendConf bool) {
	cacheKey := crd.keyByEmptyTask(orgId, agentId)
	//原空任务次数

	oldNumber := string_helper.ToInt(crd.redisClient.Get(cacheKey))
	//有未完成的任务
	if hasUnfinishedTask {
		//删除空任务次数
		crd.redisClient.Del(cacheKey)

		//原空任务次数大于等于阀值
		if oldNumber >= crd.emptyNumberThreshold {
			//agent心跳时间改为:快速(agent.fasttime)
			appendConf = true
		}
	} else {
		//新空任务次数
		newNumber := oldNumber + 1
		//设置空任务次数
		crd.redisClient.SetEx(cacheKey, strconv.Itoa(newNumber), crd.unfinishedExpire)
		//原空任务次数小于阀值 & 新空任务次数大于等于阀值
		if oldNumber < crd.emptyNumberThreshold && newNumber >= crd.emptyNumberThreshold {
			//agent心跳时间改为:普通(agent.comtime)
			appendConf = true
		}
	}

	return appendConf
}

// GetHeartbeatFlag 获取心跳标志
// @return bool true:普通, false:快速
func (crd *TaskCache) GetHeartbeatFlag(orgId int64, agentId string) bool {
	cacheKey := crd.keyByEmptyTask(orgId, agentId)
	//原空任务次数
	oldNumber := string_helper.ToInt(crd.redisClient.Get(cacheKey))
	//空任务次数是否大于等于阀值
	return oldNumber >= crd.emptyNumberThreshold
}
