package plug_calc_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/agent_service"
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
		masterModel: plug_models.CreatePlugDynamicFactory(),
		alarmModel:  alarm_models.CreateAlarmPlugFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "plug_",
		cacheExpire: 90,
	}
}

type PlugCalc struct {
	masterModel *plug_models.PlugDynamicModel
	alarmModel  *alarm_models.AlarmPlugModel
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

	//封堵配置
	config := agent_service.CreateExtendFactory().GetPlugOfPort(orgId, agentId)
	if config.Switch == enums.SwitchOff {
		return nil
	}
	//本次攻击不满足条件
	if sumWeight < config.Weight {
		//从端口虚开告警计算
		search := alarm_models.SearchConditions{
			OrgId:        orgId,
			AgentId:      agentId,
			AttackIp:     attackIp,
			HandleStatus: enum_alarm.HandleStatusNot,
			CreateAtGte:  date_helper.CalcDatetimeByMinute(attackAt, 0-int64(config.CountTime)),
		}
		fact := alarm_models.CreateAlarmPortFactory()
		calc := fact.CalcByConditions(search, fact.ColumnSumWeight())
		//累计权重不满足条件
		if calc.TotalNumber < config.Weight {
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
	if err := crd.add(orgId, agentId, attackIp, enum_plug.SourcePortVirtual, config.Type, config.PlugTime); err != nil {
		return err
	}

	//成功则缓存
	crd.cacheSuccess(orgId, agentId, attackIp)

	return nil
}

// CalcByLogin 计算: 来源登录日志
func (crd *PlugCalc) CalcByLogin(orgId int64, agentId string, attackIp string) error {

	return nil
}

// 新增
func (crd *PlugCalc) add(orgId int64, agentId string, ip string, source, plugType, duration int) error {
	now := date_helper.CurDatetime()
	var endAt string
	if plugType == enum_plug.PlugTypeTemporary {
		endAt = date_helper.CalcDatetimeByMinute(now, int64(duration))
	}

	//封堵
	modelPlug := plug_models.PlugDynamicModel{
		Id:         md5_encrypt.MD5(int_helper.Int64ToStr(orgId) + agentId + now + int_helper.IntRandom2Str(100, 999)),
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

	//告警
	alarmType := enum_alarm.TypePlugDynamicFromPort
	if source == enum_plug.SourcePortVirtual {
		alarmType = enum_alarm.TypePlugDynamicFromLogin
	}
	modelAlarm := alarm_models.AlarmPlugModel{
		Id:           variable.SnowFlake.GetId(),
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         alarmType,
		Level:        enum_alarm.LevelSerious,
		AgentAt:      now,
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		Ip:           modelPlug.Ip,
		PlugType:     modelPlug.PlugType,
		Direction:    modelPlug.Direction,
		Duration:     modelPlug.Duration,
		StartTime:    modelPlug.StartAt,
		EndTime:      modelPlug.EndAt,
	}
	if !crd.alarmModel.InsertItem(&modelAlarm) {
		return errors.New("保存告警失败")
	}

	//推送到幻势
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiAlarmSwitchPlug); hsCli != nil {
		hsData := huanshi_platform.AlarmPlug{
			AlarmBase: huanshi_platform.AlarmBase{
				Type:       huanshi_platform.AlarmTypePlug,
				OrgId:      modelAlarm.OrgId,
				AgentId:    modelAlarm.AgentId,
				Level:      modelAlarm.Level,
				NoticeType: modelAlarm.Type,
				NoticeTime: modelAlarm.CreatedAt,
			},

			Pip:        modelAlarm.Ip,
			PStartTime: modelAlarm.StartTime,
			PEndTime:   modelAlarm.EndTime,
			PMinute:    modelAlarm.Duration,
			PType:      modelAlarm.PlugType,
			PTypeTxt:   enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", modelAlarm.PlugType, enum_plug.EnumMap),
		}

		if err := hsCli.SyncAlarmPlug(hsData); err != nil {
			variable.ZapLog.Sugar().Debug(err.Error())
		}
	}

	return nil
}

// 动态封堵中已存在
func (crd *PlugCalc) hasDynamic(orgId int64, agentId string, attackIp string) bool {
	model := crd.masterModel.GetValidByOrgAndAgentAndIp(orgId, agentId, attackIp)
	if model == nil {
		return false
	} else if model.OrgId == 0 {
		return false
	} else {
		return true
	}
}

// 全局IP库中已存在
func (crd *PlugCalc) hasGlobalLib(orgId int64, attackIp string) bool {
	search := plug_models.SearchConditions{
		OrgId:      orgId,
		OpenSwitch: enums.SwitchOn,
	}
	models := plug_models.CreatePlugIpGlobalFactory().GetListByConditions(search)
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
	search := plug_models.SearchConditions{
		OrgId:      orgId,
		AgentId:    agentId,
		OpenSwitch: enums.SwitchOn,
	}
	models := plug_models.CreatePlugIpPolicyFactory().GetListByConditions(search)
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
