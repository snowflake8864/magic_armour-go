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
)

func CreateAlarmFactory() *AlarmPost {
	return &AlarmPost{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type AlarmPost struct {
	queryClient *producer_kafka.ProducerClient
}

// ToQueryBYComm 格式化参数后推送给队列
func (crd *AlarmPost) ToQueryBYComm(orgId int64, agentId, inputJson string) error {
	var params []alarmComm
	if err := json_helper.Json2Any(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	fmt.Println("params========================================")
	fmt.Println(json_helper.Encode(params))
	for _, item := range params {
		if item.Type == 0 {
			continue
		}

		//对应的模块
		module := enum_alarm.GetModuleByType(item.Type)
		switch module {
		case enum_alarm.ModuleProcess:
			crd.sendProcess(&item, orgId, agentId, now)
			break
		case enum_alarm.ModuleDll:
			crd.sendDll(&item, orgId, agentId, now)
			break
		case enum_alarm.ModuleTamper:
			crd.sendTamper(&item, orgId, agentId, now)
			break
		case enum_alarm.ModuleBlackmail:
			crd.sendBlackmail(&item, orgId, agentId, now)
			break
		case enum_alarm.ModuleUsb:
			crd.sendUsb(&item, orgId, agentId, now)
			break
		case enum_alarm.ModuleSelf:
			crd.sendSelf(&item, orgId, agentId, now)
			break
		}

	}

	return nil
}

func (crd *AlarmPost) getAgentAt(t int64, now string) string {
	if t == 0 {
		return now
	}
	return date_helper.TimestampToDatetime(t)
}

// 发送进程告警
func (crd *AlarmPost) sendProcess(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 || string_helper.IsEmpty(param.Hash) {
		return
	}

	model := alarm_models.AlarmProcessModel{
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
	fmt.Println("sendProcess ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutProcess,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// 发送模块告警
func (crd *AlarmPost) sendDll(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 || string_helper.IsEmpty(param.Hash) {
		return
	}

	model := alarm_models.AlarmDllModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		DlLHash:      param.Hash,
		DllDir:       param.Dir,
		ProcDir:      param.ExceptionProcess,
	}
	fmt.Println("sendDll ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutDll,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// 发送防篡改告警
func (crd *AlarmPost) sendTamper(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 {
		return
	}

	model := alarm_models.AlarmTamperModel{
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
	fmt.Println("sendTamper ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutTamper,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// 发送勒索告警
func (crd *AlarmPost) sendBlackmail(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 {
		return
	}

	model := alarm_models.AlarmBlackmailModel{
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
	fmt.Println("sendBlackmail ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutBlackmail,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// 发送外设告警
func (crd *AlarmPost) sendUsb(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 || string_helper.IsEmpty(param.PeripheralEid) {
		return
	}

	model := alarm_models.AlarmUsbModel{
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
	fmt.Println("sendUsb ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutUsb,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// 发送自保告警
func (crd *AlarmPost) sendSelf(param *alarmComm, orgId int64, agentId string, now string) {
	if param.Level == 0 {
		return
	}

	model := alarm_models.AlarmSelfModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         param.Type,
		Level:        param.Level,
		AgentAt:      crd.getAgentAt(param.AgentTime, now),
		HandleStatus: enum_alarm.HandleStatusNot,
		CreatedAt:    now,
		ProcHash:     param.ProcHash,
		ProcDir:      param.ProcHash,
		ProcParam:    param.ProcParam,
		FileDir:      param.FileDir,
		TargetDir:    param.TargetDir,
	}
	fmt.Println("sendSelf ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAlarmPutSelf,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceAlarm(context.Background(), msg)

}

// ToQueryBYPort 格式化参数后推送给队列
func (crd *AlarmPost) ToQueryBYPort(orgId int64, agentId, inputJson string) error {
	var params []alarmPort
	if err := json_helper.Json2Any(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	fmt.Println("params========================================")
	fmt.Println(json_helper.Encode(params))
	//统计攻击次数
	var attackTotal = make(map[string]int)
	var list []*alarm_models.AlarmPortModel
	for _, item := range params {
		if item.Weight == 0 || item.OpenPort == 0 || string_helper.IsEmpty(item.AttackIp) || string_helper.IsEmpty(item.DestinationIp) {
			break
		}
		unique := crd.uniqueBYPort(item.AttackIp, item.DestinationIp, item.OpenPort)
		if _, isOk := attackTotal[unique]; isOk {
			attackTotal[unique] += 1
		} else {
			attackTotal[unique] = 1

			list = append(list, crd.fillPort(&item, orgId, agentId, now))
		}
	}
	if len(list) > 0 {
		for _, item := range list {
			msg := &kafka_go.MessageStruct{
				Type: kafka_go.TypeAlarmPutPort,
				Data: json_helper.Encode(&item),
			}
			crd.queryClient.ProduceAlarm(context.Background(), msg)
		}

	}
	return nil
}

// 发送自保告警
func (crd *AlarmPost) fillPort(param *alarmPort, orgId int64, agentId string, now string) *alarm_models.AlarmPortModel {

	model := alarm_models.AlarmPortModel{
		OrgId:        orgId,
		AgentId:      agentId,
		Type:         enum_alarm.TypePortVirtual,
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

func (crd *AlarmPost) uniqueBYPort(attackIp string, targetIp string, targetPort int) string {
	return md5_encrypt.MD5(attackIp + targetIp + int_helper.IntToStr(targetPort))
}
