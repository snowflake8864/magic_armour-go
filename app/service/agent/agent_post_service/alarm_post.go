package agent_post_service

import (
	"context"
	"fmt"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
)

func CreateAlarmFactory() *AlarmPost {
	return &AlarmPost{
		redisClient:   redis_factory.GetOneRedisClient(),
		cachePrefix:   "lock_alarm_",
		cacheExpire:   600,
		queryClient:   producer_kafka.CreateProducerFactory(),
		supportNumber: 100,
	}
}

type AlarmPost struct {
	redisClient   *redis_factory.RedisClient
	cachePrefix   string
	cacheExpire   int
	queryClient   *producer_kafka.ProducerClient
	supportNumber int
}

// 给hash加锁
func (crd *AlarmPost) lockHash(agentId, hash string) bool {
	key := crd.cachePrefix + md5_encrypt.MD5(agentId+"_"+hash)
	if isOk, _ := crd.redisClient.Lock(key, "1", crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// ToQueryBYComm 格式化参数后推送给队列
func (crd *AlarmPost) ToQueryBYComm(orgId int64, agentId, inputJson string) error {
	var params []alarmComm
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	//fmt.Println("params========================================")
	//fmt.Println(json_helper.Encode(params))

	var processModels []*alarm_models.ProcessModel
	var dllModels []*alarm_models.DllModel
	var tamperModels []*alarm_models.TamperModel
	var blackmailModels []*alarm_models.BlackmailModel
	var usbModels []*alarm_models.UsbModel
	var selfModels []*alarm_models.SelfModel
	var outreachModels []*alarm_models.OutreachModel
	for _, item := range params {
		if item.Type == 0 {
			continue
		}

		//对应的模块
		module := enum_alarm.GetModuleByType(item.Type)
		switch module {
		case enum_alarm.ModuleProcess:
			if model := crd.formatProcess(&item, orgId, agentId, now); model != nil {
				processModels = append(processModels, model)
			}
			break
		case enum_alarm.ModuleDll:
			if model := crd.formatDll(&item, orgId, agentId, now); model != nil {
				dllModels = append(dllModels, model)
			}
			break
		case enum_alarm.ModuleTamper:
			if model := crd.formatTamper(&item, orgId, agentId, now); model != nil {
				tamperModels = append(tamperModels, model)
			}
			break
		case enum_alarm.ModuleBlackmail:
			if model := crd.formatBlackmail(&item, orgId, agentId, now); model != nil {
				blackmailModels = append(blackmailModels, model)
			}
			break
		case enum_alarm.ModuleUsb:
			if model := crd.formatUsb(&item, orgId, agentId, now); model != nil {
				usbModels = append(usbModels, model)
			}
			break
		case enum_alarm.ModuleSelf:
			if model := crd.formatSelf(&item, orgId, agentId, now); model != nil {
				selfModels = append(selfModels, model)
			}
			break
		case enum_alarm.ModuleOutreach:
			if model := crd.formatOutreach(&item, orgId, agentId, now); model != nil {
				outreachModels = append(outreachModels, model)
			}
			break
		}

	}

	if len(processModels) > 0 {
		crd.sendProcess(processModels)
	}
	if len(dllModels) > 0 {
		crd.sendDll(dllModels)
	}
	if len(tamperModels) > 0 {
		crd.sendTamper(tamperModels)
	}
	if len(blackmailModels) > 0 {
		crd.sendBlackmail(blackmailModels)
	}
	if len(usbModels) > 0 {
		crd.sendUsb(usbModels)
	}
	if len(selfModels) > 0 {
		crd.sendSelf(selfModels)
	}
	if len(outreachModels) > 0 {
		crd.sendOutreach(outreachModels)
	}

	return nil
}

func (crd *AlarmPost) getAgentAt(t int64, now string) string {
	if t == 0 {
		return now
	}
	return now
}

func (crd *AlarmPost) formatProcess(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.ProcessModel {
	if param.Level == 0 || string_helper.IsEmpty(param.Hash) {
		return nil
	}

	if !crd.lockHash(agentId, fmt.Sprintf("%d_%s", param.Type, param.Hash)) {
		return nil
	}

	return &alarm_models.ProcessModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		ProcHash:     param.Hash,
		ProcDir:      param.Dir,
		ProcId:       param.Pid,
		ProcParam:    param.PParam,
	}
}

// 发送进程告警
func (crd *AlarmPost) sendProcess(models []*alarm_models.ProcessModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}
		//alarm_job.UploadProcess(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutProcess,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}
}

func (crd *AlarmPost) formatDll(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.DllModel {
	if param.Level == 0 || string_helper.IsEmpty(param.Hash) {
		return nil
	}

	if !crd.lockHash(agentId, fmt.Sprintf("%d_%s", param.Type, param.Hash)) {
		return nil
	}

	return &alarm_models.DllModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		DllHash:      param.Hash,
		DllDir:       param.Dir,
		ProcDir:      param.ExceptionProcess,
	}
}

// 发送模块告警
func (crd *AlarmPost) sendDll(models []*alarm_models.DllModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}
		//alarm_job.UploadDll(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutDll,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}
}

func (crd *AlarmPost) formatTamper(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.TamperModel {
	if param.Level == 0 || string_helper.IsEmpty(param.Dir) {
		return nil
	}

	if !crd.lockHash(agentId, fmt.Sprintf("%d_%s_%s", param.Type, param.Hash, param.Dir)) {
		return nil
	}

	return &alarm_models.TamperModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		ProcHash:     param.Hash,
		ProcDir:      param.ExceptionProcess,
		FileDir:      param.Dir,
	}
}

// 发送防篡改告警
func (crd *AlarmPost) sendTamper(models []*alarm_models.TamperModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}

		//alarm_job.UploadTamper(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutTamper,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}

}

func (crd *AlarmPost) formatBlackmail(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.BlackmailModel {
	if param.Level == 0 || string_helper.IsEmpty(param.Dir) {
		return nil
	}

	if !crd.lockHash(agentId, fmt.Sprintf("%d_%s_%s", param.Type, param.Hash, param.Dir)) {
		return nil
	}

	return &alarm_models.BlackmailModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		ProcHash:     param.Hash,
		ProcDir:      param.ExceptionProcess,
		FileDir:      param.Dir,
	}
}

// 发送勒索告警
func (crd *AlarmPost) sendBlackmail(models []*alarm_models.BlackmailModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}

		//alarm_job.UploadBlackmail(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutBlackmail,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}
}

func (crd *AlarmPost) formatUsb(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.UsbModel {
	if param.Level == 0 || string_helper.IsEmpty(param.PeripheralEid) {
		return nil
	}

	return &alarm_models.UsbModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		UsbHash:      param.PeripheralEid,
		UsbName:      param.PeripheralName,
		UsbIntro:     param.PeripheralRemark,
	}
}

// 发送外设告警
func (crd *AlarmPost) sendUsb(models []*alarm_models.UsbModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}

		//alarm_job.UploadUsb(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutUsb,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}
}

func (crd *AlarmPost) formatSelf(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.SelfModel {
	if param.Level == 0 {
		return nil
	}

	if !crd.lockHash(agentId, fmt.Sprintf("%d_%s_%s_%s", param.Type, param.ProcHash, param.FileDir, param.TargetDir)) {
		return nil
	}

	return &alarm_models.SelfModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        enum_alarm.LevelSerious,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		ProcHash:     param.ProcHash,
		ProcDir:      param.ProcDir,
		ProcParam:    param.ProcParam,
		FileDir:      param.FileDir,
		TargetDir:    param.TargetDir,
	}
}

// 发送自保告警
func (crd *AlarmPost) sendSelf(models []*alarm_models.SelfModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}

		//alarm_job.UploadSelf(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutSelf,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}

}

func (crd *AlarmPost) formatOutreach(param *alarmComm, orgId int64, agentId string, now string) *alarm_models.OutreachModel {
	if param.Level == 0 || string_helper.IsEmpty(param.DomainName) || string_helper.IsEmpty(param.ResIp) {
		return nil
	}

	return &alarm_models.OutreachModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        enum_alarm.LevelImportant,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		DomainName:   param.DomainName,
		TargetIp:     param.ResIp,
	}
}

// 发送外联告警
func (crd *AlarmPost) sendOutreach(models []*alarm_models.OutreachModel) {
	modelLength := len(models)
	if modelLength == 0 {
		return
	}
	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}

		//alarm_job.UploadOutreach(json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type: kafka_go.TypeAlarmPutOutreach,
			Data: json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceAlarm(cxt, msg)
	}

}

// ToQueryBYPort 格式化参数后推送给队列
func (crd *AlarmPost) ToQueryBYPort(orgId int64, agentId, inputJson string) error {
	var params []alarmPort
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	//fmt.Println("params========================================")
	//fmt.Println(json_helper.Encode(params))
	//统计攻击次数
	var attackTotal = make(map[string]int)
	var list []*alarm_models.PortModel
	for _, item := range params {
		if item.Weight == 0 || item.OpenPort == 0 || string_helper.IsEmpty(item.AttackIp) || string_helper.IsEmpty(item.DestinationIp) {
			break
		}
		unique := crd.uniqueBYPort(item.AttackIp, item.DestinationIp, item.OpenPort, item.RedirectIp, item.RedirectPort, item.Weight)
		if _, isOk := attackTotal[unique]; isOk {
			attackTotal[unique] += 1
		} else {
			attackTotal[unique] = 1

			list = append(list, crd.fillPort(&item, orgId, agentId, now))
		}
	}
	if len(list) > 0 {
		cxt := context.Background()
		for _, item := range list {
			unique := crd.uniqueBYPort(item.AttackIp, item.TargetIp, item.TargetPort, item.RedirectIp, item.RedirectPort, item.Weight)
			item.AttackTotal = attackTotal[unique]

			msg := &kafka_go.MessageStruct{
				Type: kafka_go.TypeAlarmPutPort,
				Data: json_helper.Encode(&item),
			}
			crd.queryClient.ProduceAlarm(cxt, msg)
		}

	}
	return nil
}

// 发送自保告警
func (crd *AlarmPost) fillPort(param *alarmPort, orgId int64, agentId string, now string) *alarm_models.PortModel {

	model := alarm_models.PortModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         enum_alarm.GetTypeByWeight(param.Weight),
		Level:        enum_alarm.GetLevelByWeight(param.Weight),
		AgentAt:      crd.getAgentAt(param.Time, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		Weight:       param.Weight,
		AttackIp:     param.AttackIp,
		TargetIp:     param.DestinationIp,
		TargetPort:   param.OpenPort,
		RedirectIp:   param.RedirectIp,
		RedirectPort: param.RedirectPort,
		AttackTotal:  0,
	}

	return &model
}

func (crd *AlarmPost) uniqueBYPort(attackIp string, targetIp string, targetPort int, redirectIp string, redirectPort int, weight int) string {
	str := attackIp
	str += targetIp
	str += int_helper.IntToStr(targetPort)
	str += redirectIp
	str += int_helper.IntToStr(redirectPort)
	str += int_helper.IntToStr(weight)

	return md5_encrypt.MD5(str)
}
