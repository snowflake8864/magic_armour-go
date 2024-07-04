package threaten_cmd

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/redis_factory"
)

func CreateCmdFactory() *ThreatenCmd {
	return &ThreatenCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "lock_threaten_check_",
		cacheExpire: 120,
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type ThreatenCmd struct {
	redisClient *redis_factory.RedisClient
	cachePrefix string
	cacheExpire int
	queryClient *producer_kafka.ProducerClient
}

// getThreatenCli  获取威胁情报中心客户端
func (crd *ThreatenCmd) getThreatenCli(isOnline, isToken bool) (err error, cli *threaten_platform.RequestThreaten) {
	dictIds := []string{
		enum_dictionary.ThreatenSwitch,
		enum_dictionary.ThreatenUrl,
		enum_dictionary.ThreatenAppId,
		enum_dictionary.ThreatenAppKey,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if string_helper.ToInt(vals[enum_dictionary.ThreatenSwitch]) == enums.SwitchOff {
		return errors.New("威胁情报功能已关闭"), nil
	}

	//客户端
	cli = threaten_platform.CreateRequestFactory(threaten_platform.RequestConfig{
		Domain:    vals[enum_dictionary.ThreatenUrl],
		AppId:     vals[enum_dictionary.ThreatenAppId],
		AppSecret: vals[enum_dictionary.ThreatenAppKey],
	})

	//验证在线
	if isOnline {
		if err = cli.Online(); err != nil {
			return errors.New("网络连接失败： " + err.Error()), nil
		}
	}

	//生成token
	if isToken {
		if err = cli.GenerateToken(); err != nil {
			return errors.New("生成口令： " + err.Error()), nil
		}
	}

	return nil, cli
}

// SyncForPushes  推送列表： 同步数据到威胁情报中心
func (crd *ThreatenCmd) SyncForPushes() {
	variable.ZapLog.Sugar().Debug("SyncForPushes")
}

// AddFromProcessForCache  本地缓存：从进程添加
func (crd *ThreatenCmd) AddFromProcessForCache() {
	variable.ZapLog.Sugar().Debug("AddFromProcessForCache")
}

// AddFromDllForCache  本地缓存：从模块添加
func (crd *ThreatenCmd) AddFromDllForCache() {
	variable.ZapLog.Sugar().Debug("AddFromDllForCache")
}

// DeleteForCache  本地缓存：删除
func (crd *ThreatenCmd) DeleteForCache() {
	variable.ZapLog.Sugar().Debug("DeleteForCache start")

	//获取客户端
	if err, _ := crd.getThreatenCli(false, false); err != nil {
		variable.ZapLog.Sugar().Debug("DeleteForCache ", err.Error())
		return
	}

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.ThreatenCacheExpire)
	expire := string_helper.ToInt(val)
	variable.ZapLog.Sugar().Debug("DeleteForCache expire:", expire)
	_, _, dayEnd, _ := date_helper.CalcByDay(-expire)

	conditions := threaten_models.SearchConditions{
		CreateAtLte: dayEnd,
	}
	cacheFact := threaten_models.CreateCacheFactory()
	page := 0
	limit := enum_threaten.CenterSupportNumber
	for {
		page++
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}
		var ids []string
		for _, model := range models {
			ids = append(ids, model.Id)
		}

		threaten_service.CreateCacheFactory().BatchDelete(ids)
	}

	variable.ZapLog.Sugar().Debug("DeleteForCache end page:", page)
}

// SyncDefaultForCache  本地缓存：同步默认状态(未同步)数据到威胁情报中心
func (crd *ThreatenCmd) SyncDefaultForCache() {
	variable.ZapLog.Sugar().Debug("SyncDefaultForCache")
}

// SyncUnknownForCache  本地缓存：同步未知状态数据到威胁情报中心
func (crd *ThreatenCmd) SyncUnknownForCache() {
	variable.ZapLog.Sugar().Debug("SyncUnknownForCache start")

	variable.ZapLog.Sugar().Debug("SyncUnknownForCache end:")
}

// 给hash加锁
func (crd *ThreatenCmd) lockHash(hash string) bool {
	key := crd.cachePrefix + hash
	if isOk, _ := crd.redisClient.Lock(key, "1", crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// 同步缓存到威胁情报
func (crd *ThreatenCmd) syncCacheToThreaten(level int) {
	variable.ZapLog.Sugar().Debug("SyncCacheToThreaten level", level)

	//获取客户端
	err, threatenCli := crd.getThreatenCli(true, true)
	if err != nil {
		variable.ZapLog.Sugar().Debug("SyncCacheToThreaten ", err.Error())
		return
	}

	conditions := threaten_models.SearchConditions{
		Level: &level,
	}
	cacheFact := threaten_models.CreateCacheFactory()
	page := 0
	limit := enum_threaten.CenterSupportNumber
	for {
		//每次获取固定条数
		page++
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}

		//同步列表
		var syncList []threaten_platform.CheckHashes
		var syncIdMap = make(map[string]struct{})
		for _, model := range models {
			//加锁失败
			if !crd.lockHash(model.Id) {
				continue
			}
			syncIdMap[model.Id] = struct{}{}

			syncList = append(syncList, threaten_platform.CheckHashes{
				ReqType:    threaten_platform.ReqType,
				Key:        model.Id,
				Filename:   path_helper.GetFileName(model.Dir),
				Dir:        model.Dir,
				Os:         model.Os,
				ParentName: model.ParentDir,
			})
		}

		if err1 := crd.checkCacheToThreaten(threatenCli, level, &syncList, &syncIdMap); err1 != nil {
			variable.ZapLog.Sugar().Debug("SyncCacheToThreaten check: ", err.Error())
		}
	}

	variable.ZapLog.Sugar().Debug("SyncCacheToThreaten level:", level, ", page:", page)
}

// 同步缓存到威胁情报并进行处理
func (crd *ThreatenCmd) checkCacheToThreaten(cli *threaten_platform.RequestThreaten, level int,
	list *[]threaten_platform.CheckHashes, idMap *map[string]struct{}) error {

	//没有需要校验数据
	if list == nil || len(*list) == 0 {
		return nil
	}

	//去校验
	err, checkReps := cli.CheckHashes(list)
	if err != nil {
		return err
	}
	//有返回数据
	if checkReps != nil && len(checkReps) > 0 {
		for _, item := range checkReps {

		}
	}

	return nil
}
