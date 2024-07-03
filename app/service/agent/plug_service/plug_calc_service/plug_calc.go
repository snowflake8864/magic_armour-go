package plug_calc_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/alarm_service/plug_callback_alarm"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/ip_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
)

func CreatePlugCalcFactory() *PlugCalc {
	return &PlugCalc{
		masterModel: plug_models.CreateDynamicFactory(),
		alarmModel:  alarm_models.CreatePlugFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "plug_",
		cacheExpire: 90,
	}
}

type PlugCalc struct {
	masterModel *plug_models.DynamicModel
	alarmModel  *alarm_models.PlugModel
	redisClient *redis_factory.RedisClient
	cachePrefix string
	cacheExpire int
}

func (crd *PlugCalc) keySuccess(orgId int64, agentId string, attackIp string) string {
	return md5_encrypt.MD5(crd.cachePrefix + int_helper.Int64ToStr(orgId) + "_" + agentId + "_" + attackIp)
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
	_, _ = crd.redisClient.SetEx(key, "0", crd.cacheExpire)
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

	//成功则缓存
	crd.cacheSuccess(orgId, agentId, attackIp)

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

	//成功则缓存
	crd.cacheSuccess(orgId, agentId, attackIp)

	return nil
}

// 新增
func (crd *PlugCalc) Add(orgId int64, agentId string, ip string, source, plugType, duration int) error {
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
		return errors.New("保存动态封堵失败")
	}

	//同步生成告警
	if err := plug_callback_alarm.CreatePlugSaveFactory().Add(&modelPlug); err != nil {
		variable.ZapLog.Sugar().Debug(err.Error())
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
	if crd.masterModel.Count(search) > 1 {
		return false
	} else {
		return true
	}
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
	for _, model := range models {
		if ip_helper.IsIPInRange(attackIp, model.Ip) {
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
	for _, model := range models {
		if ip_helper.IsIPInRange(attackIp, model.Ip) {
			return true
		}
	}
	return false
}

// getHuanshiCli 获取幻势客户端
func (crd *PlugCalc) getHuanshiCli(pushTypeFiled string) *huanshi_platform.RequestHuanshi {
	dictIds := []string{
		enum_dictionary.HuanshiCallbackUrl,
		enum_dictionary.HuanshiCallIp,
		enum_dictionary.HuanshiSwitch,
		enum_dictionary.HuanshiUrl,
		enum_dictionary.HuanshiAlarmSwitch,
		enum_dictionary.HuanshiAlarmIp,
		enum_dictionary.HuanshiAlarmPort,
		enum_dictionary.HuanshiLogSwitch,
		enum_dictionary.HuanshiLogIp,
		enum_dictionary.HuanshiLogPort,
		pushTypeFiled,
	}
	vals := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	if len(vals) != 11 {
		return nil
	}

	if string_helper.ToInt(vals[enum_dictionary.HuanshiSwitch]) == enums.SwitchOff ||
		string_helper.ToInt(vals[enum_dictionary.HuanshiAlarmSwitch]) == enums.SwitchOff ||
		string_helper.ToInt(vals[pushTypeFiled]) == enums.SwitchOff {
		return nil
	}

	config := huanshi_platform.RequestConfig{}
	config.HuanshiCallbackUrl = vals[enum_dictionary.HuanshiCallbackUrl]
	config.HuanshiCallIp = vals[enum_dictionary.HuanshiCallIp]
	config.HuanshiUrl = vals[enum_dictionary.HuanshiUrl]
	config.HuanshiAlarmIp = vals[enum_dictionary.HuanshiAlarmIp]
	config.HuanshiAlarmPort = vals[enum_dictionary.HuanshiAlarmPort]
	config.HuanshiLogIp = vals[enum_dictionary.HuanshiLogIp]

	return huanshi_platform.CreateRequestFactory(config)
}
