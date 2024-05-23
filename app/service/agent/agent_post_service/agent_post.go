package agent_post_service

import (
	"context"
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"strings"
)

func CreateAgentFactory() *AgentPost {
	return &AgentPost{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type AgentPost struct {
	queryClient *producer_kafka.ProducerClient
}

// transTop 转换top数据
func (crd *AgentPost) transTop(s agentMonitorTop) agent_monitor_models.MonitorTop {
	return agent_monitor_models.MonitorTop{
		Id:       json_helper.Number2Int(s.Id),
		Dir:      path_helper.ToStandard(s.Dir),
		Hash:     strings.ToLower(s.Hash),
		CpuUsage: float_helper.FloatScale(json_helper.Number2Float64(s.CpuUsage), 3),
		MemSize:  s.MemSize,
		User:     s.User,
	}
}

// MonitorToQuery 保存监控信息
func (crd *AgentPost) MonitorToQuery(orgId int64, agentId, inputJson string) error {
	var param agentMonitor
	if err := json_helper.ToStruct(inputJson, &param); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	monitorDay := date_helper.SpecialDatetime(variable.DayFormat)
	fmt.Println("params========================================")
	fmt.Println(json_helper.Encode(param))
	var cpuTops []agent_monitor_models.MonitorTop
	var memTops []agent_monitor_models.MonitorTop
	if len(param.CpuTops) > 0 {
		for _, item := range param.CpuTops {
			cpuTops = append(cpuTops, crd.transTop(item))
		}
	}
	if len(param.MemTops) > 0 {
		for _, item := range param.MemTops {
			cpuTops = append(cpuTops, crd.transTop(item))
		}
	}
	selfStat := agent_monitor_models.MonitorSelf{
		CpuUsage: float_helper.FloatScale(json_helper.Number2Float64(param.CpuUsage), 3),
		MemSize:  param.Self.MemSize,
	}

	model := agent_monitor_models.MonitorModel{
		OrgId:      orgId,
		AgentId:    agentId,
		MonitorDay: monitorDay,
		CpuNumber:  param.CpuNumber,
		CpuUsage:   float_helper.FloatScale(json_helper.Number2Float64(param.CpuUsage), 3),
		CpuTops:    json_helper.Encode(cpuTops),
		MemSize:    param.MemSize,
		MemUsage:   float_helper.FloatScale(json_helper.Number2Float64(param.MemUsage), 3),
		MemTops:    json_helper.Encode(memTops),
		HdSize:     param.HdSize,
		HdUsage:    float_helper.FloatScale(json_helper.Number2Float64(param.HdUsage), 3),
		SelfStat:   json_helper.Encode(selfStat),
		CreatedAt:  now,
	}

	fmt.Println("sendMonitor ============================")
	fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeAgentMonitor,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceCommon(context.Background(), msg)

	return nil
}
