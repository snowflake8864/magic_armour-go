package kafka_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
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

var writerCommon *kafka.Writer
var writerProcess *kafka.Writer
var writerThreaten *kafka.Writer
var writerAlarm *kafka.Writer
var writerSyslog *kafka.Writer

func CreateProducerFactory() {
	config := getWriterConfig()

	writerCommon = newProducer(TopicCommon, config)
	writerProcess = newProducer(TopicProcess, config)
	writerThreaten = newProducer(TopicThreaten, config)
	writerAlarm = newProducer(TopicAlarm, config)
	writerSyslog = newProducer(TopicSyslog, config)
}

func getWriterConfig() writerConfig {
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")

	b := variable.ConfigYml.GetString("KafkaV2.Producer.Balancer")
	var balancer kafka.Balancer
	if b == "hash" {
		balancer = &kafka.Hash{}
	} else {
		balancer = &kafka.LeastBytes{}
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

// 生产消息
func writeMessage(ctx context.Context, w *kafka.Writer, msg kafka.Message) error {

	err := w.WriteMessages(ctx, msg)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("写入kafka失败，Tocic:%v,err:%v", w.Topic, err))
	}
	return err
}

func ProduceCommon(ctx context.Context, data *MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceCommon json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	msg := kafka.Message{
		Value: msgContent,
	}
	err = writeMessage(ctx, writerCommon, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceCommon 写入kafka失败，data:%v ,err:%v", data, err))
	}
}

func ProduceProcess(ctx context.Context, data *MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceProcess json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	msg := kafka.Message{
		Value: msgContent,
	}
	err = writeMessage(ctx, writerProcess, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceProcess 写入kafka失败，data:%v ,err:%v", data, err))
	}
}

func ProduceThreaten(ctx context.Context, data *MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceThreaten json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	msg := kafka.Message{
		Value: msgContent,
	}
	err = writeMessage(ctx, writerThreaten, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceThreaten 写入kafka失败，data:%v ,err:%v", data, err))
	}
}

func ProduceAlarm(ctx context.Context, data *MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceAlarm json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	msg := kafka.Message{
		Value: msgContent,
	}
	err = writeMessage(ctx, writerAlarm, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceAlarm 写入kafka失败，data:%v ,err:%v", data, err))
	}
}

func ProduceSyslog(ctx context.Context, data *MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceSyslog json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	msg := kafka.Message{
		Value: msgContent,
	}
	err = writeMessage(ctx, writerSyslog, msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceSyslog 写入kafka失败，data:%v ,err:%v", data, err))
	}
}
