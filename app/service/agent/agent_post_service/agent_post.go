package agent_post_service

import (
	"context"
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/service/agent/dir_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
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
	monitorDay := date_helper.SpecialDatetime(variable.DateFormatToday)
	//fmt.Println("params========================================")
	//fmt.Println(json_helper.Encode(param))
	var cpuTops []agent_monitor_models.MonitorTop
	var memTops []agent_monitor_models.MonitorTop
	if len(param.CpuTops) > 0 {
		for _, item := range param.CpuTops {
			cpuTops = append(cpuTops, crd.transTop(item))
		}
	}
	if len(param.MemTops) > 0 {
		for _, item := range param.MemTops {
			memTops = append(memTops, crd.transTop(item))
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

	//fmt.Println("sendMonitor ============================")
	//fmt.Println(json_helper.Encode(&model))
	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeCommonAgentPutMonitor,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceSyslog(context.Background(), msg)

	return nil
}

// SaveDir 保存目录结构
func (crd *AgentPost) SaveDir(orgId int64, agentId, inputJson string) error {
	var params []agentDir
	//忽略错误, 上传数据有为空情况
	_ = json_helper.ToAny(inputJson, &params)

	//
	dirFact := dir_service.CreateDirCurdFactory()
	//获取查询路径
	searchDir := dirFact.GetBySearch(orgId, agentId)
	if string_helper.IsEmpty(searchDir) {
		return errors.New("search dir is empty")
	}
	//
	searchDir = path_helper.ToStandard(searchDir)
	now := date_helper.CurDatetime()

	var models []dir_models.DirModel
	if len(params) > 0 {
		for _, item := range params {
			model := dir_models.DirModel{
				OrgId:      orgId,
				AgentId:    agentId,
				SearchDir:  searchDir,
				FileName:   path_helper.ToStandard(item.Dir),
				FileType:   item.Type,
				FileRw:     item.Rw,
				FileGroup:  item.Group,
				FileUser:   item.User,
				FileSize:   item.Size,
				FileHash:   strings.ToLower(item.Hash),
				StartTime:  crd.getAgentAt(item.StartTime),
				UpdateTime: crd.getAgentAt(item.UpdateTime),
				CreatedAt:  now,
			}
			models = append(models, model)
		}
	}
	if err := dirFact.BatchAdd(orgId, agentId, searchDir, &models); err != nil {
		variable.ZapLog.Error(err.Error())
	}

	return nil
}

func (crd *AgentPost) getAgentAt(t string) string {
	t = strings.Trim(t, "-")
	if string_helper.IsEmpty(t) {
		return ""
	}

	return date_helper.TimestampToDatetime(string_helper.ToInt64(t))
}
