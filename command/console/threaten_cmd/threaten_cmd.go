package threaten_cmd

import (
	"context"
	"errors"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
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
	"time"
)

func CreateCmdFactory() *ThreatenCmd {
	return &ThreatenCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "threaten_check_",
		cacheExpire: 300,
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

func (crd *ThreatenCmd) getSupportNumber() int {
	if number := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.ThreatenSupportNumber); string_helper.IsEmpty(number) {
		return enum_threaten.CenterSupportNumber
	} else {
		return string_helper.ToInt(number)
	}
}

// SyncForPushes  推送列表： 同步数据到威胁情报中心
func (crd *ThreatenCmd) SyncForPushes() {
	variable.ZapLog.Info("SyncForPushes start")
	//获取客户端
	err, threatenCli := crd.getThreatenCli(true, true)
	if err != nil {
		variable.ZapLog.Warn(fmt.Sprintf("SyncForPushes: %s ", err.Error()))
		return
	}

	conditions := threaten_models.SearchConditions{
		PushStatus: enum_threaten.PushStatusNot,
	}
	cacheFact := threaten_models.CreatePushFactory()
	page := 0
	limit := crd.getSupportNumber()
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
			variable.ZapLog.Sugar().Warn("SyncForPushes fail", err1.Error())
			break
		} else {
			threaten_service.CreatePushFactory().BatchPushed(idArr)
		}

	}

	variable.ZapLog.Sugar().Info("SyncForPushes end, page:=", page)
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
	variable.ZapLog.Debug("AddFromProcForCache start")
	//获取客户端

	if err, _ := crd.getThreatenCli(true, false); err != nil {
		variable.ZapLog.Sugar().Warn("AddFromProcForCache ", err.Error())
		return
	}

	orgModels := crd.getOrgModels()
	if orgModels == nil || len(orgModels) == 0 {
		variable.ZapLog.Info("AddFromProcForCache orgModels empty")
		return
	}

	//进程
	crd.addThreatenFromProcess(orgModels)
	//模块
	crd.addThreatenFromDll(orgModels)

	variable.ZapLog.Info("AddFromProcForCache end")
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
			var insList []threaten_models.CacheModel
			for _, procModel := range procModels {
				if _, isOk := existsMap[procModel.ProcHash]; !isOk {
					existsMap[procModel.ProcHash] = struct{}{}

					insList = append(insList, threaten_models.CacheModel{
						Id:        procModel.ProcHash,
						Source:    enum_threaten.SourceAgentUpload,
						Dir:       procModel.Dir,
						ParentDir: "",
						Os:        agentModel.OsName,
						Introduce: procModel.Introduce,
						Copyright: procModel.Copyright,
						Signature: procModel.Signature,
						Attribute: "",
						DataFrom:  enum_threaten.DataFromProcess,
					})
				}
			}
			crd.insertThreatens(insList)
		} //agentModels
		time.Sleep(3 * time.Millisecond)
	} //orgModels

	existsMap = nil
}

func (crd *ThreatenCmd) insertThreatens(insList []threaten_models.CacheModel) {
	total := len(insList)
	if total == 0 {
		return
	}
	// 循环切割
	chunkSize := 200
	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}

		chunks := insList[i:end]
		threaten_service.CreateCacheFactory().BatchAdd(&chunks, true)
	}
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
			var insList []threaten_models.CacheModel
			for _, procModel := range procModels {
				if _, isOk := existsMap[procModel.DllHash]; !isOk {
					existsMap[procModel.DllHash] = struct{}{}

					insList = append(insList, threaten_models.CacheModel{
						Id:        procModel.DllHash,
						Source:    enum_threaten.SourceAgentUpload,
						Dir:       procModel.Dir,
						ParentDir: procModel.ProcDir,
						Os:        agentModel.OsName,
						Introduce: procModel.Introduce,
						Copyright: procModel.Copyright,
						Signature: procModel.Signature,
						Attribute: procModel.Attribute,
						DataFrom:  enum_threaten.DataFromDll,
					})
				}
			}
			crd.insertThreatens(insList)
		} //agentModels
		time.Sleep(3 * time.Millisecond)
	} //orgModels

	existsMap = nil
}

// DeleteForCache  本地缓存：删除
func (crd *ThreatenCmd) DeleteForCache() {
	variable.ZapLog.Sugar().Info("DeleteForCache start")

	//获取客户端
	if err, _ := crd.getThreatenCli(false, false); err != nil {
		variable.ZapLog.Sugar().Warn("DeleteForCache ", err.Error())
		return
	}

	val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.ThreatenCacheExpire)
	expire := string_helper.ToInt(val)
	variable.ZapLog.Info(fmt.Sprintf("DeleteForCache expire: %d", expire))
	_, _, dayEnd, _ := date_helper.CalcByDay(-expire)

	conditions := threaten_models.SearchConditions{
		CreateAtLte: dayEnd,
	}
	threaten_models.CreateCacheFactory().DeleteByConditions(conditions)

	variable.ZapLog.Sugar().Info("DeleteForCache end")
}

// SyncDefaultForCache  本地缓存：同步默认状态(未同步)数据到威胁情报中心
func (crd *ThreatenCmd) SyncDefaultForCache() {
	variable.ZapLog.Info("SyncDefaultForCache start")
	crd.syncCacheToThreaten(enum_threaten.LevelDefault, 600)
	variable.ZapLog.Info("SyncDefaultForCache end")
}

// SyncUnknownForCache  本地缓存：同步未知状态数据到威胁情报中心
func (crd *ThreatenCmd) SyncUnknownForCache() {
	variable.ZapLog.Info("SyncUnknownForCache start")
	crd.syncCacheToThreaten(enum_threaten.LevelUnknown, 2000)
	variable.ZapLog.Info("SyncUnknownForCache end")
}

// 同步缓存到威胁情报
func (crd *ThreatenCmd) syncCacheToThreaten(queryLevel int, pageSize int) {

	//获取客户端
	err, threatenCli := crd.getThreatenCli(true, false)
	if err != nil {
		variable.ZapLog.Sugar().Warn("SyncCacheToThreaten getThreatenCli:", err.Error())
		return
	}

	//获取要校验的列表
	conditions := threaten_models.SearchConditions{
		Level: &queryLevel,
	}
	page := 0
	cacheFact := threaten_models.CreateCacheFactory()
	var syncList []*threaten_platform.CheckHashes
	for {
		//每次获取固定条数
		page++
		offset := (page - 1) * pageSize
		variable.ZapLog.Info(fmt.Sprintf("SyncCacheToThreaten query queryLevel=%d, page=%d, offset=%d, pageSize=%d ", queryLevel, page, offset, pageSize))
		models := cacheFact.GetItemsByConditions(conditions, offset, pageSize, "")
		if models == nil || len(models) == 0 {
			variable.ZapLog.Info(fmt.Sprintf("SyncCacheToThreaten empty queryLevel=%d, page=%d, offset=%d, pageSize=%d ", queryLevel, page, offset, pageSize))
			break
		}
		for _, model := range models {
			//短时间不允许校验
			if !crd.allowForCheck(model.Id) {
				continue
			}

			syncList = append(syncList, &threaten_platform.CheckHashes{
				ReqType:    threaten_platform.ReqType,
				Key:        model.Id,
				Filename:   path_helper.GetFileName(model.Dir),
				Dir:        model.Dir,
				Os:         model.Os,
				ParentName: model.ParentDir,
				DataFrom:   model.DataFrom,
			})
		}

		//未同步一次最多验证1回
		if queryLevel == enum_threaten.LevelDefault && page > 1 {
			variable.ZapLog.Info(fmt.Sprintf("SyncCacheToThreaten break queryLevel=%d, page=%d, offset=%d, pageSize=%d ", queryLevel, page, offset, pageSize))
			break
		}
	}

	if tmpLen := len(syncList); tmpLen > 0 {
		variable.ZapLog.Info(fmt.Sprintf("SyncCacheToThreaten queryLevel=%d,len=%d start", queryLevel, tmpLen))
		//认证最大数
		supportNumber := crd.getSupportNumber()
		for i := 0; i < tmpLen; i += supportNumber {
			end := i + supportNumber
			if end > tmpLen {
				end = tmpLen
			}

			if err1 := crd.checkCacheToThreaten(threatenCli, queryLevel, syncList[i:end]); err1 != nil {
				variable.ZapLog.Sugar().Warnln("SyncCacheToThreaten check: ", err1.Error())
			}
		}

		variable.ZapLog.Info(fmt.Sprintf("SyncCacheToThreaten queryLevel=%d,len=%d end", queryLevel, tmpLen))
	}
}

func (crd *ThreatenCmd) keyForCheck(id string) string {
	return crd.cachePrefix + id
}

// 允许校验
func (crd *ThreatenCmd) allowForCheck(id string) bool {
	if isExist, _ := crd.redisClient.Exists(crd.keyForCheck(id)); isExist {
		return false
	}
	return true
}

// 添加最新校验校验
func (crd *ThreatenCmd) addForCheck(id string, level int) {
	_, _ = crd.redisClient.SetEx(crd.keyForCheck(id), int_helper.IntToStr(level), crd.cacheExpire)
}

// 同步缓存到威胁情报并进行处理
func (crd *ThreatenCmd) checkCacheToThreaten(cli *threaten_platform.RequestThreaten, preLevel int,
	syncList []*threaten_platform.CheckHashes) error {

	//没有需要校验数据
	if syncList == nil || len(syncList) == 0 {
		variable.ZapLog.Info("syncList == nil")
		return nil
	}
	//哈希和来源的映射
	var idAndDataFrom = make(map[string]int)
	for _, item := range syncList {
		idAndDataFrom[item.Key] = item.DataFrom
	}

	//去校验
	err, checkReps := cli.CheckHashes(syncList)
	if err != nil {
		variable.ZapLog.Warn(err.Error())
		return nil
	}
	//没有返回数据
	if checkReps == nil || len(checkReps) == 0 {
		//
		crd.handleUnknown(preLevel, idAndDataFrom)
		return nil
	}
	//有返回数据
	ctx := context.Background()
	now := date_helper.CurDatetime()

	//需要快速处理的数据, 减少时间差
	var processWhites []string
	var processBlacks []string
	var dllWhites []string
	var dllBlacks []string
	hasLevelOnly := make(map[int][]string)

	for _, item := range checkReps {
		hash := strings.ToLower(item.Hash)
		dataFrom, isOk := idAndDataFrom[hash]
		if !isOk {
			continue
		} else {
			delete(idAndDataFrom, hash)
		}

		//数据
		model := threaten_models.CacheModel{
			Id:        strings.ToLower(item.Hash),
			Level:     item.Level,
			DataFrom:  dataFrom,
			CheckAt:   now,
			CheckInfo: item.Attr,
		}

		//添加最新校验校验, 防止重复校验
		crd.addForCheck(model.Id, model.Level)

		//进入队列
		crd.queryClient.ProduceThreaten(ctx, &kafka_go.MessageStruct{
			Type: kafka_go.TypeThreatenSyncRewrite,
			Data: json_helper.Encode(&model),
		})

		//需要快速处理
		if enum_threaten.IsWhite(model.Level) { //白名单
			if dataFrom == enum_threaten.DataFromProcess {
				processWhites = append(processWhites, model.Id)
			} else if dataFrom == enum_threaten.DataFromDll {
				dllWhites = append(dllWhites, model.Id)
			} else {
				processWhites = append(processWhites, model.Id)
				dllWhites = append(dllWhites, model.Id)
			}

		} else if enum_threaten.IsBlack(model.Level) { //黑名单
			if dataFrom == enum_threaten.DataFromProcess {
				processBlacks = append(processBlacks, model.Id)
			} else if dataFrom == enum_threaten.DataFromDll {
				dllBlacks = append(dllBlacks, model.Id)
			} else {
				processBlacks = append(processBlacks, model.Id)
				dllBlacks = append(dllBlacks, model.Id)
			}
		} else { //仅有等级的
			hasLevelOnly[model.Level] = append(hasLevelOnly[model.Level], model.Id)
		}
	}

	//快速处理
	crd.fastHandleHasLevelAndInfluenceProc(processWhites, enum_threaten.LevelWhite, enum_threaten.DataFromProcess)
	crd.fastHandleHasLevelAndInfluenceProc(processBlacks, enum_threaten.LevelBlack, enum_threaten.DataFromProcess)
	crd.fastHandleHasLevelAndInfluenceProc(dllWhites, enum_threaten.LevelWhite, enum_threaten.DataFromDll)
	crd.fastHandleHasLevelAndInfluenceProc(dllBlacks, enum_threaten.LevelBlack, enum_threaten.DataFromDll)
	crd.handleUnknown(preLevel, idAndDataFrom)
	crd.fastHandleHasLevelOnly(hasLevelOnly)

	return nil
}

// 快速处理快速有等级和影响进程的
func (crd *ThreatenCmd) fastHandleHasLevelAndInfluenceProc(idArr []string, level, dataFrom int) {
	if len(idArr) > 0 {
		variable.ZapLog.Info(fmt.Sprintf("level:%d, dataFrom:%d, idArr:%s", level, dataFrom, json_helper.Encode(idArr)))
		threaten_service.CreateCacheFactory().FastBatchSetHasLevelAndInfluenceProc(idArr, level, dataFrom)
	}
}

// 快速处理快速有等级和影响进程的
func (crd *ThreatenCmd) fastHandleHasLevelOnly(param map[int][]string) {
	if len(param) > 0 {
		fact := threaten_service.CreateCacheFactory()
		for level, ids := range param {
			fact.FastBatchSetHasLevelOnly(ids, level)
		}
	}
}

// 处理未知
func (crd *ThreatenCmd) handleUnknown(level int, idAndDataFrom map[string]int) {
	if len(idAndDataFrom) > 0 && level == enum_threaten.LevelDefault {
		var idArr []string
		for id := range idAndDataFrom {
			idArr = append(idArr, id)
		}
		//要直接修改，防止下次验证时获取到重复数据
		threaten_service.CreateCacheFactory().BatchUnknown(idArr)
	}
}
