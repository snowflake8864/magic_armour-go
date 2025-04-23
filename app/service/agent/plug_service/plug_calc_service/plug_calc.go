package plug_calc_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/alarm_service/plug_callback_alarm"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/huanshi_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
)

func CreatePlugCalcFactory() *PlugCalc {
	return &PlugCalc{
		masterModel:   plug_models.CreateDynamicFactory(),
		alarmModel:    alarm_models.CreatePlugFactory(),
		redisClient:   redis_factory.GetOneRedisClient(),
		successPrefix: "plug_ok_",
		SuccessExpire: 60,
		lockPrefix:    "plug_lock_",
		lockExpire:    15,
	}
}

type PlugCalc struct {
	masterModel   *plug_models.DynamicModel
	alarmModel    *alarm_models.PlugModel
	redisClient   *redis_factory.RedisClient
	successPrefix string
	SuccessExpire int
	lockPrefix    string
	lockExpire    int
}

func (crd *PlugCalc) keySuccess(orgId int64, agentId string, attackIp string) string {
	return md5_encrypt.MD5(crd.successPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId + "_" + attackIp)
}

// 是否存在成功数据
func (crd *PlugCalc) existSuccess(orgId int64, agentId string, attackIp string) bool {
	key := crd.keySuccess(orgId, agentId, attackIp)
	if exist, err := crd.redisClient.Exists(key); err == nil {
		return exist
	} else {
		return false
	}
}

// 缓存成功数据
func (crd *PlugCalc) cacheSuccess(orgId int64, agentId string, attackIp string) {
	key := crd.keySuccess(orgId, agentId, attackIp)
	_, _ = crd.redisClient.SetEx(key, "0", crd.SuccessExpire)
}

func (crd *PlugCalc) keyLock(orgId int64, agentId string, attackIp string) string {
	return md5_encrypt.MD5(crd.lockPrefix + int_helper.Int64ToStr(orgId) + "_" + agentId + "_" + attackIp)
}

func (crd *PlugCalc) toLock(orgId int64, agentId string, attackIp string) bool {
	key := crd.keyLock(orgId, agentId, attackIp)
	isOk, _ := crd.redisClient.Lock(key, "1", crd.lockExpire)
	return isOk
}

func (crd *PlugCalc) unLock(orgId int64, agentId string, attackIp string) {
	key := crd.keyLock(orgId, agentId, attackIp)
	_, _ = crd.redisClient.Del(key)
}

// CalcByPort 计算: 来源端口虚开告警
func (crd *PlugCalc) CalcByPort(orgId int64, agentId string, attackIp string, sumWeight int, attackAt string) error {
	if crd.existSuccess(orgId, agentId, attackIp) {
		return nil
	}

	//终端开关
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(orgId, agentId, agentFact.ColumnConfig(false, false))
	if agentModel == nil || agentModel.PlugSwitch == enums.SwitchOff {
		return nil
	}

	//封堵配置 - 端口
	plugConfig := agent_service.CreateExtendFactory().GetPlugOfPort(orgId, agentId)
	if plugConfig.Switch == enums.SwitchOff {
		return nil
	}
	//本次攻击不满足条件
	if sumWeight < plugConfig.Weight {
		//从端口虚开告警计算
		search := alarm_models.SearchConditions{
			OrgId:        orgId,
			AgentId:      agentId,
			AttackIp:     attackIp,
			HandleStatus: enum_alarm.HandleStatusNot,
			CreateAtGte:  date_helper.CalcDatetimeByMinute(attackAt, 0-int64(plugConfig.CountTime)),
		}
		fact := alarm_models.CreatePortFactory()
		calc := fact.CalcByConditions(search, fact.ColumnSumWeight())
		//累计权重不满足条件
		if calc.TotalNumber < plugConfig.Weight {
			return nil
		}
	}

	//动态封堵中已存在
	if crd.hasDynamic(orgId, agentId, attackIp) {
		return nil
	}
	//全局IP库中已存在
	if crd.hasGlobalLib(orgId, attackIp) {
		return nil
	}
	//终端IP库中已存在
	if crd.hasAgentLib(orgId, agentId, attackIp) {
		return nil
	}

	//入库
	if err := crd.Add(orgId, agentId, attackIp, enum_plug.SourcePortVirtual, plugConfig.Type, plugConfig.PlugTime); err != nil {
		return err
	}

	return nil
}

// CalcByLogin 计算: 来源登录日志
func (crd *PlugCalc) CalcByLogin(orgId int64, agentId string, attackIp string) error {
	if crd.existSuccess(orgId, agentId, attackIp) {
		return nil
	}

	//终端开关
	agentFact := agent_models.CreateAgentFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(orgId, agentId, agentFact.ColumnConfig(false, false))
	if agentModel == nil || agentModel.PlugSwitch == enums.SwitchOff {
		return nil
	}

	//封堵配置 - 口令破解
	plugConfig := agent_service.CreateExtendFactory().GetPlugOfLogin(orgId, agentId)
	if plugConfig.Switch == enums.SwitchOff {
		return nil
	}

	//动态封堵中已存在
	if crd.hasDynamic(orgId, agentId, attackIp) {
		return nil
	}
	//全局IP库中已存在
	if crd.hasGlobalLib(orgId, attackIp) {
		return nil
	}
	//终端IP库中已存在
	if crd.hasAgentLib(orgId, agentId, attackIp) {
		return nil
	}

	//入库
	if err := crd.Add(orgId, agentId, attackIp, enum_plug.SourceLoginSolve, plugConfig.Type, plugConfig.PlugTime); err != nil {
		return err
	}

	return nil
}

// 新增
func (crd *PlugCalc) Add(orgId int64, agentId string, ip string, source, plugType, duration int) error {
	//防止并发
	if !crd.toLock(orgId, agentId, ip) {
		return errors.New("toLock fail")
	}
	defer crd.unLock(orgId, agentId, ip)

	now := date_helper.CurDatetime()
	var endAt string
	if plugType == enum_plug.PlugTypeTemporary {
		endAt = date_helper.CalcDatetimeByMinute(now, int64(duration))
	}

	//封堵
	modelPlug := plug_models.DynamicModel{
		Id:         variable.SnowFlake.GetId(),
		OrgId:      orgId,
		AgentId:    agentId,
		PlugType:   plugType,
		Source:     source,
		Ip:         ip,
		Direction:  enum_plug.DirectionOut + enum_plug.DirectionIn,
		Duration:   duration,
		StartAt:    now,
		EndAt:      endAt,
		OpenSwitch: enums.SwitchOn,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if !crd.masterModel.InsertItem(&modelPlug) {
		return errors.New("save fail")
	}

	//成功则缓存
	crd.cacheSuccess(orgId, agentId, ip)

	//新增任务
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypePluggingDownload, "")

	//同步生成告警
	if err := plug_callback_alarm.CreatePlugSaveFactory().Add(&modelPlug); err != nil {
		variable.ZapLog.Info(err.Error())
	}

	return nil
}

// 动态封堵中已存在
func (crd *PlugCalc) hasDynamic(orgId int64, agentId string, attackIp string) bool {
	SwitchOn := enums.SwitchOn
	search := plug_models.SearchConditions{
		OrgId:      orgId,
		AgentId:    agentId,
		Ip:         attackIp,
		OpenSwitch: &SwitchOn,
	}
	if crd.masterModel.Count(search) >= 1 {
		return true
	} else {
		return false
	}
}

func (crd *PlugCalc) InIpEffect(nowTimestamp int64, effectMode int, effectStartAt, effectEndAt string) bool {
	if effectMode == enum_plug.EffectModeRange {
		//正确： startTimestamp <= nowTimestamp <= endTimestamp
		if !string_helper.IsEmpty(effectStartAt) {
			//没有开始: 开始时间大于现在时间
			if startTimestamp := date_helper.DatetimeToTimestamp(effectStartAt+" 00:00:00", variable.DateFormatNow); startTimestamp > nowTimestamp {
				return false
			}
		}
		if !string_helper.IsEmpty(effectEndAt) {
			//已经结束: 现在时间大于结束时间
			if endTimestamp := date_helper.DatetimeToTimestamp(effectEndAt+" 23:59:59", variable.DateFormatNow); nowTimestamp > endTimestamp {
				return false
			}

		}
	}

	return true
}

// 全局IP库中已存在
func (crd *PlugCalc) hasGlobalLib(orgId int64, attackIp string) bool {
	SwitchOn := enums.SwitchOn
	search := plug_models.SearchConditions{
		OrgId:      orgId,
		OpenSwitch: &SwitchOn,
	}
	models := plug_models.CreateIpGlobalFactory().GetListByConditions(search)
	if models == nil || len(models) == 0 {
		return false
	}
	nowTimestamp := date_helper.CurTimestamp()
	for _, model := range models {
		if crd.InIpEffect(nowTimestamp, model.EffectMode, model.EffectStartAt, model.EffectEndAt) && ip_helper.IsIPInRange(attackIp, model.Ip) {
			return true
		}
	}
	return false
}

// 终端IP库中已存在
func (crd *PlugCalc) hasAgentLib(orgId int64, agentId string, attackIp string) bool {
	SwitchOn := enums.SwitchOn
	search := plug_models.SearchConditions{
		OrgId:      orgId,
		AgentId:    agentId,
		OpenSwitch: &SwitchOn,
	}
	models := plug_models.CreateIpPolicyFactory().GetListByConditions(search)
	if models == nil || len(models) == 0 {
		return false
	}
	nowTimestamp := date_helper.CurTimestamp()
	for _, model := range models {
		if crd.InIpEffect(nowTimestamp, model.EffectMode, model.EffectStartAt, model.EffectEndAt) && ip_helper.IsIPInRange(attackIp, model.Ip) {
			return true
		}
	}
	return false
}

// getHuanshiCli 获取幻势客户端
func (crd *PlugCalc) getHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
	return huanshi_helper.GetHuanshiCli(pushTypeFiled)
}
