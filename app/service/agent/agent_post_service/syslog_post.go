package agent_post_service

import (
	"context"
	"fmt"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateSyslogFactory() *SyslogPost {
	return &SyslogPost{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type SyslogPost struct {
	queryClient *producer_kafka.ProducerClient
}

// ToQuery 格式化参数后推送给队列
func (crd *SyslogPost) ToQuery(orgId int64, agentId, inputJson string) error {
	var params []syslogComm
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	fmt.Println("params========================================")
	fmt.Println(json_helper.Encode(params))
	for _, item := range params {
		logType := json_helper.Number2Int64(item.LogType)
		//对应的模块
		switch logType {
		case 1:
			crd.sendDns(&item, orgId, agentId, now)
			break
		case 2:
			crd.sendOuter(&item, orgId, agentId, now)
			break
		case 3:
			crd.sendInner(&item, orgId, agentId, now)
			break
		case 4:
			crd.sendProcess(&item, orgId, agentId, now)
			break
		case 5:
			crd.sendLogin(&item, orgId, agentId, now)
			break
		}

	}

	return nil
}

func (crd *SyslogPost) getAgentAt(t int64, now string) string {
	if t == 0 {
		return now
	}
	return date_helper.TimestampToDatetime(t)
}

// 发送DNS日志
func (crd *SyslogPost) sendDns(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.SyslogDnsModel{
		OrgId:      orgId,
		AgentId:    agentId,
		AgentAt:    crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt:  now,
		ProcHash:   param.Hash,
		ProcId:     json_helper.Number2Int(param.PId),
		ProcDir:    param.PDir,
		DomainName: param.DomainName,
		TargetIp:   param.ResIp,
	}
	fmt.Println("sendDns ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutDns,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

// 发送对内通讯日志
func (crd *SyslogPost) sendInner(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.SyslogInnerModel{
		OrgId:            orgId,
		AgentId:          agentId,
		AgentAt:          crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt:        now,
		ProcHash:         param.Hash,
		ProcId:           json_helper.Number2Int(param.PId),
		ProcDir:          param.PDir,
		TargetIp:         param.ResIp,
		TargetPort:       json_helper.Number2Int(param.RsPort),
		TargetPortStatus: json_helper.Number2Int(param.ResPortStatus),
		SourceIp:         param.SourceIp,
		SourcePort:       json_helper.Number2Int(param.SourcePort),
		Proto:            json_helper.Number2Int(param.Proto),
	}
	fmt.Println("sendInner ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutInner,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

// 发送对外通讯日志
func (crd *SyslogPost) sendOuter(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.SyslogOuterModel{
		OrgId:            orgId,
		AgentId:          agentId,
		AgentAt:          crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt:        now,
		ProcHash:         param.Hash,
		ProcId:           json_helper.Number2Int(param.PId),
		ProcDir:          param.PDir,
		TargetIp:         param.ResIp,
		TargetPort:       json_helper.Number2Int(param.RsPort),
		TargetPortStatus: json_helper.Number2Int(param.ResPortStatus),
		SourceIp:         param.SourceIp,
		SourcePort:       json_helper.Number2Int(param.SourcePort),
		Proto:            json_helper.Number2Int(param.Proto),
	}
	fmt.Println("sendOuter ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutOuter,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

// 发送进程启动日志
func (crd *SyslogPost) sendProcess(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.SyslogProcessModel{
		OrgId:       orgId,
		AgentId:     agentId,
		AgentAt:     crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt:   now,
		ProcHash:    param.Hash,
		ProcId:      json_helper.Number2Int(param.PId),
		ProcDir:     param.PDir,
		ProcParam:   param.PParam,
		ParentHash:  param.PpHash,
		ParentId:    json_helper.Number2Int(param.PpId),
		ParentDir:   param.PpDir,
		ParentParam: param.PpParam,
	}
	fmt.Println("sendProcess ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutProcess,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

// 发送登录志
func (crd *SyslogPost) sendLogin(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.SyslogLoginModel{
		OrgId:     orgId,
		AgentId:   agentId,
		AgentAt:   crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt: now,
		SourceIp:  param.Ip,
		Username:  param.Username,
		Type:      param.Type,
		Status:    json_helper.Number2Int(param.Status),
	}
	fmt.Println("sendLogin ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutLogin,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}
