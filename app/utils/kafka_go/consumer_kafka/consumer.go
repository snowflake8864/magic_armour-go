package consumer_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/jobs/agent_job"
	"goskeleton/app/jobs/alarm_job"
	"goskeleton/app/jobs/email_job/send_alarm_cycle"
	"goskeleton/app/jobs/org_job"
	"goskeleton/app/jobs/plug_job"
	"goskeleton/app/jobs/proc_job"
	"goskeleton/app/jobs/syslog_job"
	"goskeleton/app/jobs/threaten_job"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"strings"
	"time"
)

func CreateConsumerFactory() *ConsumerClient {
	return &ConsumerClient{}
}

// 结构体
type ConsumerClient struct {
}

func (k *ConsumerClient) newConsumer(name string) *kafka.Reader {
	//主题
	topic := variable.ConfigYml.GetString("KafkaV2." + name + ".Name")
	//地址
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")
	//消费者信息
	groupId := variable.ConfigYml.GetString("KafkaV2.Consumer.GroupId")
	commitInterval := variable.ConfigYml.GetDuration("KafkaV2.Consumer.CommitInterval")
	startOffset := variable.ConfigYml.GetInt64("KafkaV2.Consumer.StartOffset")

	configReader := kafka.ReaderConfig{
		Topic:          topic,
		Brokers:        brokers,
		GroupID:        groupId + "_" + topic,
		CommitInterval: commitInterval,
		GroupBalancers: []kafka.GroupBalancer{
			&kafka.RangeGroupBalancer{},
		},
		StartOffset: startOffset, //kafka.FirstOffset
		MinBytes:    0,           // 最小读取字节数
		MaxBytes:    10e6,        // 最大读取字节数
	}
	variable.ZapLog.Sugar().Info(fmt.Sprintf("GroupID:%s", configReader.GroupID))
	return kafka.NewReader(configReader)
}

// 公司每次发货硬件（硬盘写入速度）都不一致, 是否全量并发需要看情况
func (k *ConsumerClient) newConsumers(name string) []*kafka.Reader {
	var readers []*kafka.Reader
	fullDose := variable.ConfigYml.GetInt("KafkaV2.Consumer.FullDose")
	partition := 1
	if fullDose == 1 {
		partition = variable.ConfigYml.GetInt("KafkaV2." + name + ".Partition")
	}
	//fmt.Println("newConsumers=====================================")
	//fmt.Println("fullDose:", fullDose)
	//fmt.Println("partition:", partition)

	for i := 0; i < partition; i++ { // 这里指定消费者的数量
		readers = append(readers, k.newConsumer(name))
	}
	return readers
}

func (k *ConsumerClient) readMessage(ctx context.Context, r *kafka.Reader) (*kafka_go.MessageStruct, error) {
	result := &kafka_go.MessageStruct{}

	msg, err := r.ReadMessage(ctx)
	if err != nil {
		variable.ZapLog.Sugar().Debug(fmt.Sprintf("topic:%s, 读kafka失败，err:%v", r.Config().Topic, err))
		return result, err
	} else if len(msg.Value) == 0 {
		return nil, nil
	} else {
		err = json.Unmarshal(msg.Value, result)
		if err != nil {
			variable.ZapLog.Sugar().Debug(fmt.Sprintf("topic:%s, json unmarshal msg value err，msg:%v,err:%v", r.Config().Topic, result, err))
		}
		//variable.ZapLog.Sugar().Debug(fmt.Sprintf("topic=%s,partition=%d,offset=%d,key=%s,result=%v", msg.Topic, msg.Partition, msg.Offset, msg.Key, result))
		return result, err

	}
}

func (k *ConsumerClient) consumeMessage(ctx context.Context, r *kafka.Reader) {
	for {
		data, err := k.readMessage(ctx, r)
		if err != nil {
			continue // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		} else if data == nil {
			// 当前消息读取为空时，并不退出for终止所有后续消费，休眠后继续处理
			time.Sleep(3 * time.Second)
			continue
		}
		switch data.Type {
		case kafka_go.TypeCommonAgentDelete:
			agent_job.Delete(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeCommonAgentOpenAutoWhite:
			agent_job.OpenAutoWhite(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeCommonOrgClean:
			org_job.Clean(data.OrgId)
			break
		case kafka_go.TypeCommonEmailAlarmCycle:
			send_alarm_cycle.Send(data.OrgId)
			break

		case kafka_go.TypeProcessBatchPutMaster:
			proc_job.BatchUploadProcess(data.OrgId, data.AgentId, data.Data)
			break
		case kafka_go.TypeDllBatchPutMaster:
			proc_job.BatchUploadDll(data.OrgId, data.AgentId, data.Data)
			break

		case kafka_go.TypeAlarmPutProcess:
			alarm_job.UploadProcess(data.Data)
			break
		case kafka_go.TypeAlarmPutDll:
			alarm_job.UploadDll(data.Data)
			break
		case kafka_go.TypeAlarmPutTamper:
			alarm_job.UploadTamper(data.Data)
			break
		case kafka_go.TypeAlarmPutBlackmail:
			alarm_job.UploadBlackmail(data.Data)
			break
		case kafka_go.TypeAlarmPutUsb:
			alarm_job.UploadUsb(data.Data)
			break
		case kafka_go.TypeAlarmPutSelf:
			alarm_job.UploadSelf(data.Data)
			break
		case kafka_go.TypeAlarmPutOutreach:
			alarm_job.UploadOutreach(data.Data)
			break
		case kafka_go.TypeAlarmPutPort:
			alarm_job.UploadPort(data.Data)
			break

		case kafka_go.TypeSyslogPutDns:
			syslog_job.UploadDns(data.Data)
			break
		case kafka_go.TypeSyslogPutOuter:
			syslog_job.UploadOuter(data.Data)
			break
		case kafka_go.TypeSyslogPutInner:
			syslog_job.UploadInner(data.Data)
			break
		case kafka_go.TypeSyslogPutProcess:
			syslog_job.UploadProcess(data.Data)
			break
		case kafka_go.TypeSyslogPutLoginSuccess:
			syslog_job.UploadLoginSuccess(data.Data)
			break
		case kafka_go.TypeSyslogPutLoginFail:
			syslog_job.UploadLoginFail(data.Data)
			break
		case kafka_go.TypeCommonAgentPutMonitor:
			agent_job.UploadMonitor(data.Data)
			break

		case kafka_go.TypeThreatenSyncRewrite:
			threaten_job.HandelRewire(data.Data)
			break
		case kafka_go.TypeThreatenImport:
			threaten_job.HandelImport(data.Data)
			break

		case kafka_go.TypeCrontabAlarmStat:
			alarm_job.StatInStartAtAndEndAt(data.Data, data.AgentId, data.AgentOS)
			break
		case kafka_go.TypeCrontabAgentSafe:
			agent_job.HandleSafe(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeCrontabPlugDisable:
			plug_job.DynamicDisable(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeCrontabTaskFormProcess:
			proc_job.TaskFromProcess(data.OrgId, data.AgentId, data.Data)
			break
		case kafka_go.TypeCrontabTaskFormDll:
			proc_job.TaskFromDll(data.OrgId, data.AgentId, data.Data)
			break
		case kafka_go.TypeCrontabThreatenRepairProc:
			threaten_job.RepairProc(string_helper.ToInt(data.AgentId), string_helper.ToInt(data.AgentOS), data.Data)
			break
		default:
			variable.ZapLog.Sugar().Info(fmt.Sprintf("consumeMessage data.Type: %s", data.Type))
		}
	}
}

// ConsumeCommon 消费主题： Common
func (k *ConsumerClient) ConsumeCommon() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicCommon)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeProcess 消费主题： Process
func (k *ConsumerClient) ConsumeProcess() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicProcess)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeDll 消费主题： Dll
func (k *ConsumerClient) ConsumeDll() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicDll)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeAlarm 消费主题： Alarm
func (k *ConsumerClient) ConsumeAlarm() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicAlarm)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeSyslog 消费主题： Syslog
func (k *ConsumerClient) ConsumeSyslog() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicSyslog)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeThreaten 消费主题： Threaten
func (k *ConsumerClient) ConsumeThreaten() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicThreaten)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}

// ConsumeCrontab 消费主题： Crontab
func (k *ConsumerClient) ConsumeCrontab() {
	ctx := context.Background()
	consumers := k.newConsumers(kafka_go.TopicCrontab)

	for _, consumer := range consumers {
		go k.consumeMessage(ctx, consumer)
	}

	// 主 goroutine 可以选择等待或者处理其他逻辑
	select {}
}
