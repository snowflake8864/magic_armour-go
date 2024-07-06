package consumer_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/jobs/agent_job"
	"goskeleton/app/jobs/alarm_job"
	"goskeleton/app/jobs/org_job"
	"goskeleton/app/jobs/proc_job"
	"goskeleton/app/jobs/syslog_job"
	"goskeleton/app/jobs/threaten_job"
	"goskeleton/app/utils/kafka_go"
	"strings"
)

func getReaderConfig() kafka.ReaderConfig {
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")
	groupId := variable.ConfigYml.GetString("KafkaV2.Consumer.GroupId")
	commitInterval := variable.ConfigYml.GetDuration("KafkaV2.Consumer.CommitInterval")
	startOffset := variable.ConfigYml.GetInt64("KafkaV2.Consumer.StartOffset")

	return kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupId,
		CommitInterval: commitInterval,
		GroupBalancers: []kafka.GroupBalancer{
			&kafka.RangeGroupBalancer{},
		},
		StartOffset: startOffset, //kafka.FirstOffset
		MinBytes:    0,           // 最小读取字节数
		MaxBytes:    10e6,        // 最大读取字节数
	}

}

func newConsumer(name string, configReader kafka.ReaderConfig) *kafka.Reader {
	topic := variable.ConfigYml.GetString("KafkaV2." + name + ".Name")
	// 配置消费者
	configReader.Topic = topic
	return kafka.NewReader(configReader)
}

func CreateConsumerFactory() *ConsumerClient {
	config := getReaderConfig()
	return &ConsumerClient{
		Common:   newConsumer(kafka_go.TopicCommon, config),
		Process:  newConsumer(kafka_go.TopicProcess, config),
		Threaten: newConsumer(kafka_go.TopicThreaten, config),
		Alarm:    newConsumer(kafka_go.TopicAlarm, config),
		Syslog:   newConsumer(kafka_go.TopicSyslog, config),
		Crontab:  newConsumer(kafka_go.TopicCrontab, config),
	}

}

// 结构体
type ConsumerClient struct {
	Common   *kafka.Reader
	Process  *kafka.Reader
	Threaten *kafka.Reader
	Alarm    *kafka.Reader
	Syslog   *kafka.Reader
	Crontab  *kafka.Reader
}

func (k *ConsumerClient) readMessage(ctx context.Context, r *kafka.Reader) (*kafka_go.MessageStruct, error) {
	result := &kafka_go.MessageStruct{}

	msg, err := r.ReadMessage(ctx)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("topic:%s, 读kafka失败，err:%v", r.Config().Topic, err))
		return result, err
	} else {
		err = json.Unmarshal(msg.Value, result)
		if err != nil {
			variable.ZapLog.Error(fmt.Sprintf("topic:%s, json unmarshal msg value err，msg:%v,err:%v", r.Config().Topic, result, err))
		}
		fmt.Println(fmt.Sprintf("topic=%s,partition=%d,offset=%d,key=%s,result=%v", msg.Topic, msg.Partition, msg.Offset, msg.Key, result))
		return result, err

	}
}

// ConsumeCommon 消费主题： Common
func (k *ConsumerClient) ConsumeCommon() {
	for {
		data, err := k.readMessage(context.Background(), k.Common)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
		case kafka_go.TypeAgentDelete:
			agent_job.Delete(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeAgentOpenAutoWhite:
			agent_job.OpenAutoWhite(data.OrgId, data.AgentId)
			break
		case kafka_go.TypeOrgClean:
			org_job.Clean(data.OrgId)
			break
		case kafka_go.TypeAgentPutMonitor:
			agent_job.UploadMonitor(data.Data)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("ReadCommon result:%T", data))
		}
	}
}

// ConsumeProcess 消费主题： Process
func (k *ConsumerClient) ConsumeProcess() {
	for {
		data, err := k.readMessage(context.Background(), k.Process)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
		case kafka_go.TypeProcessPutMaster:
			proc_job.UploadProcess(data.Data)
			break
		case kafka_go.TypeProcessPutRunning:
			proc_job.UploadProcessRunning(data.Data)
			break
		case kafka_go.TypeDllPutMaster:
			proc_job.UploadDll(data.Data)
			break
		case kafka_go.TypeDllPutRunning:
			proc_job.UploadDllRunning(data.Data)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("readerProcess result:%T", data))
		}
	}
}

// ConsumeAlarm 消费主题： Alarm
func (k *ConsumerClient) ConsumeAlarm() {
	for {
		data, err := k.readMessage(context.Background(), k.Alarm)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
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
		case kafka_go.TypeAlarmPutPort:
			alarm_job.UploadPort(data.Data)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("readerAlarm result:%T", data))
		}
	}
}

// ConsumeSyslog 消费主题： Syslog
func (k *ConsumerClient) ConsumeSyslog() {
	for {
		data, err := k.readMessage(context.Background(), k.Syslog)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
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
		case kafka_go.TypeSyslogPutLogin:
			syslog_job.UploadLogin(data.Data)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("readerSyslog result:%T", data))
		}
	}
}

// ConsumeThreaten 消费主题： Threaten
func (k *ConsumerClient) ConsumeThreaten() {
	for {
		data, err := k.readMessage(context.Background(), k.Threaten)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
		case kafka_go.TypeThreatenSyncRewrite:
			threaten_job.HandelRewire(data.Data)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("readerThreaten result:%T", data))
		}
	}
}

// ConsumeCrontab 消费主题： Crontab
func (k *ConsumerClient) ConsumeCrontab() {
	for {
		data, err := k.readMessage(context.Background(), k.Crontab)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}

		switch data.Type {
		case kafka_go.TypeCrontabStatAlarm:
			alarm_job.StatInStartAtAndEndAt(data.Data, data.AgentId, data.AgentOS)
			break
		default:
			variable.ZapLog.Info(fmt.Sprintf("readerCrontab result:%T", data))
		}
	}
}
