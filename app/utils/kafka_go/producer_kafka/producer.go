package producer_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kafka_go"
	"log"
	"os"
	"strings"
	"time"
)

type writerConfig struct {
	Addr         []string
	Balancer     kafka.Balancer
	WriteTimeout time.Duration
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks int
	Async        bool
}

func getWriterConfig() writerConfig {
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")

	b := variable.ConfigYml.GetString("KafkaV2.Producer.Balancer")
	var balancer kafka.Balancer
	if b == "roundRobin" {
		balancer = &kafka.RoundRobin{}
	} else if b == "hash" {
		balancer = &kafka.Hash{}
	} else if b == "leastBytes" {
		balancer = &kafka.LeastBytes{}
	} else {
		balancer = &kafka.RoundRobin{}
	}
	return writerConfig{
		Addr:         brokers,
		Balancer:     balancer,
		WriteTimeout: variable.ConfigYml.GetDuration("KafkaV2.Producer.WriteTimeout"),
		BatchSize:    variable.ConfigYml.GetInt("KafkaV2.Producer.BatchSize"),
		BatchTimeout: variable.ConfigYml.GetDuration("KafkaV2.Producer.BatchTimeout"),
		RequiredAcks: variable.ConfigYml.GetInt("KafkaV2.Producer.RequiredAcks"),
		Async:        variable.ConfigYml.GetBool("KafkaV2.Producer.Async"),
	}
}

func newProducer(name string, c writerConfig) *kafka.Writer {
	topic := variable.ConfigYml.GetString("KafkaV2." + name + ".Name")
	//partition := variable.ConfigYml.GetInt("KafkaV2." + name + ".Partition")

	// 配置消费者
	return &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		Topic:        topic,
		Balancer:     c.Balancer,
		WriteTimeout: c.WriteTimeout * time.Second,
		BatchSize:    c.BatchSize,
		BatchTimeout: c.BatchTimeout * time.Millisecond,
		RequiredAcks: kafka.RequiredAcks(c.RequiredAcks),
		Async:        c.Async,
		Logger:       log.New(os.Stdout, "kafka-go: ", log.LstdFlags),
	}
}

func CreateProducerFactory() *ProducerClient {
	config := getWriterConfig()
	return &ProducerClient{
		Common:   newProducer(kafka_go.TopicCommon, config),
		Process:  newProducer(kafka_go.TopicProcess, config),
		Dll:      newProducer(kafka_go.TopicDll, config),
		Threaten: newProducer(kafka_go.TopicThreaten, config),
		Alarm:    newProducer(kafka_go.TopicAlarm, config),
		Syslog:   newProducer(kafka_go.TopicSyslog, config),
		Crontab:  newProducer(kafka_go.TopicCrontab, config),
	}
}

// 结构体
type ProducerClient struct {
	Common   *kafka.Writer
	Process  *kafka.Writer
	Dll      *kafka.Writer
	Threaten *kafka.Writer
	Alarm    *kafka.Writer
	Syslog   *kafka.Writer
	Crontab  *kafka.Writer
}

// 生产消息
func (k *ProducerClient) writeMessage(ctx context.Context, w *kafka.Writer, msg kafka.Message) error {
	err := w.WriteMessages(ctx, msg)
	if err != nil {
		variable.ZapLog.Warn(fmt.Sprintf("写入kafka失败，Tocic:%v,err:%v", w.Topic, err))
	}
	return err
}

// 生成信息
func (k *ProducerClient) makeMessage(data *kafka_go.MessageStruct) kafka.Message {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Info(fmt.Sprintf("json marshal err data:%v,err:%v", data, err))
	}

	return kafka.Message{
		Value: msgContent,
	}
}

func (k *ProducerClient) ProduceCommon(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Common, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceCommon data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceProcess(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Process, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceProcess data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceDll(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Dll, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceDll data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceThreaten(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Threaten, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceThreaten data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceAlarm(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Alarm, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceAlarm data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceSyslog(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Syslog, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceSyslog data:%v ,err:%v", data, err))
	}
}

func (k *ProducerClient) ProduceCrontab(ctx context.Context, data *kafka_go.MessageStruct) {
	if err := k.writeMessage(ctx, k.Crontab, k.makeMessage(data)); err != nil {
		variable.ZapLog.Info(fmt.Sprintf("ProduceCrontab data:%v ,err:%v", data, err))
	}
}
