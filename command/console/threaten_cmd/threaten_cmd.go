package threaten_cmd

import (
	"context"
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/agent/threaten_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/threaten_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateCmdFactory() *ThreatenCmd {
	return &ThreatenCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "lock_threaten_check_",
		cacheExpire: 50,
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

// Cache2Pushes  推送列表： 本地缓存数据同步威胁情报推送列表
func (crd *ThreatenCmd) Cache2Pushes() {
	variable.ZapLog.Sugar().Debug("CacheToPushes start")

	variable.ZapLog.Sugar().Debug("CacheToPushes end")
}

// SyncForPushes  推送列表： 同步数据到威胁情报中心
func (crd *ThreatenCmd) SyncForPushes() {
	variable.ZapLog.Sugar().Debug("SyncForPushes start")
	//获取客户端
	err, threatenCli := crd.getThreatenCli(true, true)
	if err != nil {
		variable.ZapLog.Sugar().Debug("SyncForPushes ", err.Error())
		return
	}

	conditions := threaten_models.SearchConditions{
		PushStatus: enum_threaten.PushStatusNot,
	}
	cacheFact := threaten_models.CreatePushFactory()
	page := 0
	limit := enum_threaten.CenterSupportNumber
	for page = 1; page <= 10; page++ {
		//每次获取固定条数
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}
		//同步列表
		var syncList []threaten_platform.PushHashes
		var idArr []string
		for _, model := range models {
			idArr = append(idArr, model.Id)

			syncList = append(syncList, threaten_platform.PushHashes{
				Hash:      model.Id,
				Type:      model.Level - 1,
				Filename:  path_helper.GetFileName(model.Dir),
				Dir:       model.Dir,
				Os:        model.Os,
				Uname:     model.CreateName,
				ExecDir:   model.ParentDir,
				Attribute: model.Attribute,
				Remark:    model.Introduce,
				Copyright: model.Copyright,
			})
		}
		if err1 := threatenCli.PushHashes(&syncList); err1 != nil {
			variable.ZapLog.Sugar().Debug("SyncForPushes fail", err1.Error())
			break
		} else {
			threaten_service.CreatePushFactory().BatchPushed(idArr)
		}

	}

	variable.ZapLog.Sugar().Debug("SyncForPushes end, page:=", page)
}

func (crd *ThreatenCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

func (crd *ThreatenCmd) getAgentModels(orgId int64) []*agent_models.AgentModel {
	agentFact := agent_models.CreateAgentFactory()
	return agentFact.GetListByConditions(agent_models.SearchConditions{
		OrgId: orgId,
	}, agentFact.ColumnBasic(true, false))
}

// AddFromProcForCache  本地缓存：从进程或模块添加
func (crd *ThreatenCmd) AddFromProcForCache() {
	variable.ZapLog.Sugar().Debug("AddFromProcForCache start")
	//获取客户端

	if err, _ := crd.getThreatenCli(true, false); err != nil {
		variable.ZapLog.Sugar().Debug("AddFromProcForCache ", err.Error())
		return
	}

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Sugar().Debug("AddFromProcForCache orgModels empty")
		return
	}

	//进程
	crd.addThreatenFromProcess(orgModels)
	//模块
	crd.addThreatenFromDll(orgModels)

	variable.ZapLog.Sugar().Debug("AddFromProcForCache end")
}

func (crd *ThreatenCmd) addThreatenFromProcess(orgModels []*organization_models.OrganizationModel) {
	existsMap := make(map[string]struct{})
	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			condition := proc_models.SearchConditions{
				OrgId:   orgModel.Id,
				AgentId: agentModel.AgentId,
			}
			//获取全部进程
			procModels := proc_models.CreateProcessFactory().GetListByConditions(condition)
			if procModels == nil || len(procModels) == 0 {
				continue
			}
			for _, procModel := range procModels {

				if _, isOk := existsMap[procModel.ProcHash]; isOk {
					param := threaten_models.CacheModel{
						Id:        procModel.ProcHash,
						Source:    enum_threaten.SourceAgentUpload,
						Dir:       procModel.Dir,
						ParentDir: "",
						Os:        agentModel.OsName,
						Introduce: procModel.Introduce,
						Copyright: procModel.Copyright,
						Signature: procModel.Signature,
						Attribute: "",
						DataFrom:  enum_process.DataFromProcess,
					}
					threaten_service.CreateCacheFactory().Add(&param, true)
				} else {
					existsMap[procModel.ProcHash] = struct{}{}
				}
			}
		} //agentModels
	} //orgModels

	existsMap = nil
}

func (crd *ThreatenCmd) addThreatenFromDll(orgModels []*organization_models.OrganizationModel) {
	existsMap := make(map[string]struct{})
	for _, orgModel := range orgModels {
		//获取全部终端
		agentModels := crd.getAgentModels(orgModel.Id)
		if agentModels == nil || len(agentModels) == 0 {
			continue
		}

		for _, agentModel := range agentModels {
			condition := proc_models.SearchConditions{
				OrgId:   orgModel.Id,
				AgentId: agentModel.AgentId,
			}
			//获取全部模块
			procModels := proc_models.CreateDllFactory().GetListByConditions(condition)
			if procModels == nil || len(procModels) == 0 {
				continue
			}
			for _, procModel := range procModels {

				if _, isOk := existsMap[procModel.DllHash]; isOk {
					param := threaten_models.CacheModel{
						Id:        procModel.DllHash,
						Source:    enum_threaten.SourceAgentUpload,
						Dir:       procModel.Dir,
						ParentDir: procModel.ProcDir,
						Os:        agentModel.OsName,
						Introduce: procModel.Introduce,
						Copyright: procModel.Copyright,
						Signature: procModel.Signature,
						Attribute: procModel.Attribute,
						DataFrom:  enum_process.DataFromDll,
					}
					threaten_service.CreateCacheFactory().Add(&param, true)
				} else {
					existsMap[procModel.DllHash] = struct{}{}
				}
			}
		} //agentModels
	} //orgModels

	existsMap = nil
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

// RepairProcForCache  本地缓存：修复进程模块（进程和模块上传时入队列, 威胁情报验证回写入队列, 有时间差）
func (crd *ThreatenCmd) RepairProcForCache() {
	variable.ZapLog.Sugar().Debug("RepairProcForCache start")
	//66f00330d9014b59cd3d100a886ddd2d
	floorAt := date_helper.TimestampToDatetime(date_helper.AddTimestampByHour(-3))
	variable.ZapLog.Sugar().Debug("RepairProcForCache floorAt: ", floorAt)

	cacheFact := threaten_models.CreateCacheFactory()
	page := 0
	limit := enum_threaten.CenterSupportNumber
	conditions := threaten_models.SearchConditions{
		LevelList:  []int{enum_threaten.LevelWhite, enum_threaten.LevelBlack},
		CheckAtGte: floorAt,
	}
	ctx := context.Background()
	for {
		page++
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}
		var processWhites, processBlacks, dllWhites, dllBlacks []string
		for _, model := range models {
			if model.DataFrom == enum_threaten.DataFromProcess {
				if model.Level == enum_threaten.LevelWhite {
					processWhites = append(processWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					processBlacks = append(processBlacks, model.Id)
				}
			} else if model.DataFrom == enum_threaten.DataFromDll {
				if model.Level == enum_threaten.LevelWhite {
					dllWhites = append(dllWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					dllBlacks = append(dllBlacks, model.Id)
				}
			} else {
				if model.Level == enum_threaten.LevelWhite {
					processWhites = append(processWhites, model.Id)
					dllWhites = append(dllWhites, model.Id)
				} else if model.Level == enum_threaten.LevelBlack {
					processBlacks = append(processBlacks, model.Id)
					dllBlacks = append(dllBlacks, model.Id)
				}
			}
		}

		//进入队列
		if len(processWhites) > 0 {
			crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabThreatenRepairProc,
				AgentId: int_helper.IntToStr(enum_threaten.DataFromProcess),
				AgentOS: int_helper.IntToStr(enum_threaten.LevelWhite),
				Data:    json_helper.Encode(processWhites),
			})
		}
		if len(processBlacks) > 0 {
			crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabThreatenRepairProc,
				AgentId: int_helper.IntToStr(enum_threaten.DataFromProcess),
				AgentOS: int_helper.IntToStr(enum_threaten.LevelBlack),
				Data:    json_helper.Encode(processBlacks),
			})
		}
		if len(dllWhites) > 0 {
			crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabThreatenRepairProc,
				AgentId: int_helper.IntToStr(enum_threaten.DataFromDll),
				AgentOS: int_helper.IntToStr(enum_threaten.LevelWhite),
				Data:    json_helper.Encode(dllWhites),
			})
		}
		if len(dllBlacks) > 0 {
			crd.queryClient.ProduceCrontab(ctx, &kafka_go.MessageStruct{
				Type:    kafka_go.TypeCrontabThreatenRepairProc,
				AgentId: int_helper.IntToStr(enum_threaten.DataFromDll),
				AgentOS: int_helper.IntToStr(enum_threaten.LevelBlack),
				Data:    json_helper.Encode(dllBlacks),
			})
		}

	}
	variable.ZapLog.Sugar().Debug("RepairProcForCache page:", page)
}

// SyncDefaultForCache  本地缓存：同步默认状态(未同步)数据到威胁情报中心
func (crd *ThreatenCmd) SyncDefaultForCache() {
	variable.ZapLog.Sugar().Debug("SyncDefaultForCache start")
	crd.syncCacheToThreaten(enum_threaten.LevelDefault)
	variable.ZapLog.Sugar().Debug("SyncDefaultForCache end")
}

// SyncUnknownForCache  本地缓存：同步未知状态数据到威胁情报中心
func (crd *ThreatenCmd) SyncUnknownForCache() {
	variable.ZapLog.Sugar().Debug("SyncUnknownForCache start")
	crd.syncCacheToThreaten(enum_threaten.LevelUnknown)
	variable.ZapLog.Sugar().Debug("SyncUnknownForCache end")
}

// 同步缓存到威胁情报
func (crd *ThreatenCmd) syncCacheToThreaten(level int) {
	variable.ZapLog.Sugar().Debug("SyncCacheToThreaten level:", level)

	//获取客户端
	err, threatenCli := crd.getThreatenCli(true, false)
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
		//未同步一次最多验证3回
		if level == enum_threaten.LevelDefault && page >= 10 {
			break
		}
		//每次获取固定条数
		page++
		offset := (page - 1) * limit
		models := cacheFact.GetItemsByConditions(conditions, offset, limit, "")
		if models == nil || len(models) == 0 {
			break
		}

		//同步列表
		var syncList []threaten_platform.CheckHashes
		var idAndDataFrom = make(map[string]int)
		for _, model := range models {
			idAndDataFrom[model.Id] = model.DataFrom

			syncList = append(syncList, threaten_platform.CheckHashes{
				ReqType:    threaten_platform.ReqType,
				Key:        model.Id,
				Filename:   path_helper.GetFileName(model.Dir),
				Dir:        model.Dir,
				Os:         model.Os,
				ParentName: model.ParentDir,
			})
		}
		//variable.ZapLog.Sugar().Debug("idAndDataFrom============================")
		//variable.ZapLog.Sugar().Debug(json_helper.Encode(idAndDataFrom))
		//fmt.Println(json_helper.Encode(syncList))
		if err1 := crd.checkCacheToThreaten(threatenCli, level, &syncList, idAndDataFrom); err1 != nil {
			variable.ZapLog.Sugar().Debug("SyncCacheToThreaten check: ", err.Error())
		}
	}

	variable.ZapLog.Sugar().Debug("SyncCacheToThreaten level:", level, ", page:", page)
}

// 同步缓存到威胁情报并进行处理
func (crd *ThreatenCmd) checkCacheToThreaten(cli *threaten_platform.RequestThreaten, level int,
	list *[]threaten_platform.CheckHashes, idAndDataFrom map[string]int) error {
	variable.ZapLog.Sugar().Debug("checkCacheToThreaten=============================", level)

	//没有需要校验数据
	if list == nil || len(*list) == 0 {
		variable.ZapLog.Sugar().Debug("list == nil")
		return nil
	}

	//去校验
	err, checkReps := cli.CheckHashes(list)
	if err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
		return nil
	}
	//没有返回数据
	if checkReps == nil || len(*checkReps) == 0 {
		//
		crd.handleUnUnknown(level, idAndDataFrom)
		return nil
	}
	//有返回数据
	ctx := context.Background()
	now := date_helper.CurDatetime()

	for _, item := range *checkReps {
		hash := strings.ToLower(item.Hash)
		dataFrom, isOk := idAndDataFrom[hash]
		if !isOk {
			continue
		} else {
			delete(idAndDataFrom, hash)
		}

		model := threaten_models.CacheModel{
			Id:        strings.ToLower(item.Hash),
			Level:     item.Level,
			DataFrom:  dataFrom,
			CheckAt:   now,
			CheckInfo: item.Attr,
		}
		//进入队列
		crd.queryClient.ProduceThreaten(ctx, &kafka_go.MessageStruct{
			Type: kafka_go.TypeThreatenSyncRewrite,
			Data: json_helper.Encode(&model),
		})
	}

	//
	crd.handleUnUnknown(level, idAndDataFrom)

	return nil
}

// 威胁情报不存在的hash & 原来是未同步
func (crd *ThreatenCmd) handleUnUnknown(level int, idAndDataFrom map[string]int) {
	if len(idAndDataFrom) > 0 && level == enum_threaten.LevelDefault {
		var idArr []string
		for id := range idAndDataFrom {
			idArr = append(idArr, id)
		}
		//要直接修改，防止下次验证时获取到重复数据
		threaten_service.CreateCacheFactory().BatchUnknown(idArr)
	}
}
