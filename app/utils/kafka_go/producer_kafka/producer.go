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
	"sync"
	"time"
)

func newProducer(name string) *kafka.Writer {
	//主题
	topic := variable.ConfigYml.GetString("KafkaV2." + name + ".Name")
	//partition := variable.ConfigYml.GetInt("KafkaV2." + name + ".Partition")
	//地址
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")

	//分区分配策略
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

	//写超时
	writeTimeout := variable.ConfigYml.GetDuration("KafkaV2.Producer.WriteTimeout")
	if writeTimeout == 0 {
		writeTimeout = 10
	}
	//每批次写的数量
	batchSize := variable.ConfigYml.GetInt("KafkaV2.Producer.BatchSize")
	if batchSize == 0 {
		batchSize = 20
	}
	//每批次写超时
	batchTimeout := variable.ConfigYml.GetDuration("KafkaV2.Producer.BatchTimeout")
	if batchTimeout == 0 {
		batchTimeout = 100
	}
	//确认方式
	requiredAcks := variable.ConfigYml.GetInt("KafkaV2.Producer.RequiredAcks")
	//是否异步
	async := variable.ConfigYml.GetBool("KafkaV2.Producer.Async")

	// 配置消费者
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     balancer,
		WriteTimeout: writeTimeout * time.Second,
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout * time.Millisecond,
		RequiredAcks: kafka.RequiredAcks(requiredAcks),
		Async:        async,
		Logger:       log.New(os.Stdout, "kafka-go: ", log.LstdFlags),
	}
}

// 包级别变量
var (
	instance *ProducerClient
	once     sync.Once // 保证初始化只执行一次
)

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

func CreateProducerFactory() *ProducerClient {
	once.Do(func() {
		// 初始化逻辑（可包含运行时条件）
		instance = &ProducerClient{
			Common:   newProducer(kafka_go.TopicCommon),
			Process:  newProducer(kafka_go.TopicProcess),
			Dll:      newProducer(kafka_go.TopicDll),
			Threaten: newProducer(kafka_go.TopicThreaten),
			Alarm:    newProducer(kafka_go.TopicAlarm),
			Syslog:   newProducer(kafka_go.TopicSyslog),
			Crontab:  newProducer(kafka_go.TopicCrontab),
		}

	})
	return instance
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
		Key:   nil,
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
