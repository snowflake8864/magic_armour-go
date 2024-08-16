package agent_post_service

import (
	"context"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/syslog_service/job_callback_syslog"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateSyslogFactory() *SyslogPost {
	return &SyslogPost{
		queryClient:   producer_kafka.CreateProducerFactory(),
		supportNumber: 100,
	}
}

type SyslogPost struct {
	queryClient   *producer_kafka.ProducerClient
	supportNumber int
}

// ToQuery 格式化参数后推送给队列
func (crd *SyslogPost) ToQuery(orgId int64, agentId, inputJson string) error {
	var params []syslogComm
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	//fmt.Println("params========================================")
	//fmt.Println(json_helper.Encode(params))

	//没有其他业务的
	var innerModels []syslog_models.InnerModel
	var outerModels []syslog_models.OuterModel
	var processModels []syslog_models.ProcessModel
	for _, item := range params {
		logType := json_helper.Number2Int64(item.LogType)
		//对应的模块
		switch logType {
		case 1:
			crd.sendDns(&item, orgId, agentId, now)
			break
		case 2:
			outerModels = append(outerModels, crd.formatOuter(&item, orgId, agentId, now))
			break
		case 3:
			innerModels = append(innerModels, crd.formatInner(&item, orgId, agentId, now))
			break
		case 4:
			processModels = append(processModels, crd.formatProcess(&item, orgId, agentId, now))
			break
		case 5:
			crd.sendLogin(&item, orgId, agentId, now)
			break
		}

	}

	if err := crd.sendInner(&innerModels); err != nil {
		variable.ZapLog.Sugar().Debug("sendInner: " + err.Error())
	}

	if err := crd.sendOuter(&outerModels); err != nil {
		variable.ZapLog.Sugar().Debug("sendOuter: " + err.Error())
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

	model := syslog_models.DnsModel{
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
	//fmt.Println("sendDns ============================")
	//fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutDns,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

func (crd *SyslogPost) formatInner(param *syslogComm, orgId int64, agentId string, now string) syslog_models.InnerModel {
	return syslog_models.InnerModel{
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
}

// 发送对内通讯日志
func (crd *SyslogPost) sendInner(models *[]syslog_models.InnerModel) error {
	modelLength := len(*models)
	if modelLength == 0 {
		return nil
	}
	/*
		data := *models
		cxt := context.Background()
		for i := 0; i < modelLength; i += crd.supportNumber {
			end := i + crd.supportNumber
			if end > modelLength {
				end = modelLength
			}
			//syslog_job.UploadInner(json_helper.Encode(data[i:end]))

			msg := &kafka_go.MessageStruct{
				Type: kafka_go.TypeSyslogPutInner,
				Data: json_helper.Encode(data[i:end]),
			}
			crd.queryClient.ProduceDll(cxt, msg)
		}
	*/

	if err := job_callback_syslog.CreateInnerUploadFactory().BatchAdd(models); err != nil {
		return err
	}

	return nil
}

func (crd *SyslogPost) formatOuter(param *syslogComm, orgId int64, agentId string, now string) syslog_models.OuterModel {

	return syslog_models.OuterModel{
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
}

// 发送对外通讯日志
func (crd *SyslogPost) sendOuter(models *[]syslog_models.OuterModel) error {

	modelLength := len(*models)
	if modelLength == 0 {
		return nil
	}
	/*
		data := *models
		cxt := context.Background()
		for i := 0; i < modelLength; i += crd.supportNumber {
			end := i + crd.supportNumber
			if end > modelLength {
				end = modelLength
			}
			syslog_job.UploadOuter(json_helper.Encode(data[i:end]))
			msg := &kafka_go.MessageStruct{
				Type: kafka_go.TypeSyslogPutOuter,
				Data: json_helper.Encode(data[i:end]),
			}
			crd.queryClient.ProduceDll(cxt, msg)
		}*/

	if err := job_callback_syslog.CreateOuterUploadFactory().BatchAdd(models); err != nil {
		return err
	}

	return nil
}

func (crd *SyslogPost) formatProcess(param *syslogComm, orgId int64, agentId string, now string) syslog_models.ProcessModel {

	return syslog_models.ProcessModel{
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
}

// 发送进程启动日志
func (crd *SyslogPost) sendProcess(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.ProcessModel{
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
	//fmt.Println("sendProcess ============================")
	//fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutProcess,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}

// 发送登录志
func (crd *SyslogPost) sendLogin(param *syslogComm, orgId int64, agentId string, now string) {

	model := syslog_models.LoginModel{
		OrgId:     orgId,
		AgentId:   agentId,
		AgentAt:   crd.getAgentAt(json_helper.Number2Int64(param.Time), now),
		CreatedAt: now,
		SourceIp:  param.Ip,
		Username:  param.Username,
		Type:      param.Type,
		Status:    json_helper.Number2Int(param.Status),
	}
	//fmt.Println("sendLogin ============================")
	//fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeSyslogPutLogin,
		Data: json_helper.Encode(&model),
	}

	crd.queryClient.ProduceSyslog(context.Background(), msg)
}
