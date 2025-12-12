package consumer_kafka_v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/core/event_manage"
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
	"sync"
	"time"
)

func CreateConsumerFactory() *ConsumerClient {
	return &ConsumerClient{}
}

type ConsumerClient struct {
}

// 创建分区消费者
func (k *ConsumerClient) createPartitionConsumers(topicNames []string) (readers []*kafka.Reader) {
	//地址
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")
	//消费者信息
	groupId := variable.ConfigYml.GetString("KafkaV2.Consumer.GroupId")
	commitInterval := variable.ConfigYml.GetDuration("KafkaV2.Consumer.CommitInterval")
	startOffset := variable.ConfigYml.GetInt64("KafkaV2.Consumer.StartOffset")
	//每个主题消费者是否与主题分区数量保持一致(0|1)
	fullDose := variable.ConfigYml.GetInt("KafkaV2.Consumer.FullDose")

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, name := range topicNames {
		//主题是否存在
		topicName := variable.ConfigYml.GetString("KafkaV2." + name + ".Name")
		if string_helper.IsEmpty(topicName) {
			continue
		}
		// 主题暂停
		if variable.ConfigYml.GetInt("KafkaV2."+name+".Suspend") == 1 {
			continue
		}

		//默认消费者数量
		partitions := 1
		if fullDose == 1 {
			//指定消费者数量
			if numbers := variable.ConfigYml.GetInt("KafkaV2." + name + ".Numbers"); numbers > 0 {
				partitions = numbers
			} else {
				//跟分区一致
				partitions = variable.ConfigYml.GetInt("KafkaV2." + name + ".Partition")
			}
		}
		variable.ZapLog.Info(fmt.Sprintf("Consumers topic=%s, partitions=%d", topicName, partitions))
		for p := 0; p < partitions; p++ {
			wg.Add(1)
			go func(topic string, partition int) {
				defer wg.Done()
				//不能同时指定Partition 和 GroupID
				reader := kafka.NewReader(kafka.ReaderConfig{
					Brokers: brokers,
					GroupID: groupId,
					Topic:   topic,
					//Partition:      partition,
					GroupBalancers: []kafka.GroupBalancer{&kafka.RangeGroupBalancer{}},
					MinBytes:       10e3,
					MaxBytes:       10e6,
					StartOffset:    startOffset, //kafka.FirstOffset
					CommitInterval: commitInterval,
				})
				mu.Lock()
				readers = append(readers, reader)
				mu.Unlock()
			}(topicName, p)
		}

	}

	wg.Wait()
	return readers
}

// Consume 消费全部主题
func (k *ConsumerClient) Consume() {
	topicNames := []string{
		kafka_go.TopicCommon,
		kafka_go.TopicProcess,
		kafka_go.TopicDll,
		kafka_go.TopicThreaten,
		kafka_go.TopicAlarm,
		kafka_go.TopicSyslog,
		kafka_go.TopicCrontab,
	}

	// 创建分区消费者
	readers := k.createPartitionConsumers(topicNames)

	// 信号处理
	ctx, cancel := context.WithCancel(context.Background())

	// 启动所有消费者
	var consumerWg sync.WaitGroup
	for i, reader := range readers {
		consumerWg.Add(1)
		go func(idx int, r *kafka.Reader) {
			_ = k.consumeTopic(ctx, idx, r)
		}(i, reader)

	}

	// 注册关闭事件
	eventManageFactory := event_manage.CreateEventManageFactory()
	if _, exists := eventManageFactory.Get(variable.EventDestroyPrefix + "KafkaConsumer"); !exists {
		eventManageFactory.Set(variable.EventDestroyPrefix+"KafkaConsumer", func(args ...interface{}) {
			variable.ZapLog.Info("Trigger Kafka shutdown")
			cancel() // 通知所有消费者协程退出
			variable.ZapLog.Info("All consumers stopped")

			// 显式关闭readers
			variable.ZapLog.Info(fmt.Sprintf("readers len=%d", len(readers)))
			for index, reader := range readers {
				variable.ZapLog.Info(fmt.Sprintf("event Closing topic=%s, index=%d", reader.Config().Topic, index))
				if err := reader.Close(); err != nil {
					variable.ZapLog.Warn(fmt.Sprintf("Close error: %v", err))
				}
			}
			variable.ZapLog.Info("KafkaConsumer Quit====================================")
		})
	}

	// 等待所有消费者完成（正常情况下会一直阻塞在这里）
	consumerWg.Wait()
	variable.ZapLog.Info("Consume function exiting")
}

func (k *ConsumerClient) consumeTopic(ctx context.Context, idx int, reader *kafka.Reader) error {
	defer func() {
		if err := recover(); err != nil {
			variable.ZapLog.Error(fmt.Sprintf("Consumer [%d] panic: %v", idx, err))
		}
		variable.ZapLog.Info(fmt.Sprintf("Consumer [%d] exited", idx))
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				variable.ZapLog.Info(fmt.Sprintf("Consumer [%d] FetchMessage err: %v", idx, err))
				time.Sleep(100 * time.Millisecond) //短暂休眠避免CPU忙等待
				continue
			}

			if err := k.handleData(m.Value); err != nil {
				variable.ZapLog.Info(fmt.Sprintf("Consumer [%d] handleData err: %v", idx, err))
			}

			if err := reader.CommitMessages(ctx, m); err != nil && !errors.Is(err, context.Canceled) {
				variable.ZapLog.Info(fmt.Sprintf("Consumer [%d] CommitMessages err: %v", idx, err))
			}
		}
	}
}

func (k *ConsumerClient) handleData(msg []byte) error {
	data := &kafka_go.MessageStruct{}
	if err := json.Unmarshal(msg, data); err != nil {
		return fmt.Errorf("data.Unmarshal: %v", err)
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
		variable.ZapLog.Info(fmt.Sprintf("handleData data.Type: %s", data.Type))
	}

	return nil
}
