package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kakfa_v2/model"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
	"log"
	"os"
	"strings"
	"time"
)

var configYml ymlconfig_interf.YmlConfigInterf

type writerConfigWrite struct {
	Addr         []string
	Balancer     kafka.Balancer
	WriteTimeout time.Duration
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks int
	Async        bool
}

var configWriter *writerConfigWrite

var writerCommon *kafka.Writer
var writerProcess *kafka.Writer
var writerThreaten *kafka.Writer
var writerAlarm *kafka.Writer
var writerSyslog *kafka.Writer

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	fmt.Println("kakfa_v2 producer---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	getWriterConfig()

	writerCommon = newProducer("TopicComm")
}

func newProducer(name string) *kafka.Writer {
	topic := configYml.GetString("KafkaV2.TopicList." + name + ".Name")
	//partition := configYml.GetInt("KafkaV2.TopicList." + name + ".Partition")

	// 配置消费者
	return &kafka.Writer{
		Addr:         kafka.TCP(configWriter.Addr...),
		Topic:        topic,
		Balancer:     configWriter.Balancer,
		WriteTimeout: configWriter.WriteTimeout * time.Second,
		BatchSize:    configWriter.BatchSize,
		BatchTimeout: configWriter.BatchTimeout * time.Millisecond,
		RequiredAcks: kafka.RequiredAcks(configWriter.RequiredAcks),
		Async:        configWriter.Async,
		Logger:       log.New(os.Stdout, "kafka-go: ", log.LstdFlags),
	}
}

func getWriterConfig() {
	addr := configYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")

	b := configYml.GetString("KafkaV2.Producer.Balancer")
	var balancer kafka.Balancer
	if b == "hash" {
		balancer = &kafka.Hash{}
	} else {
		balancer = &kafka.LeastBytes{}
	}
	configWriter = &writerConfigWrite{
		Addr:         brokers,
		Balancer:     balancer,
		WriteTimeout: configYml.GetDuration("KafkaV2.Producer.WriteTimeout"),
		BatchSize:    configYml.GetInt("KafkaV2.Producer.BatchSize"),
		BatchTimeout: configYml.GetDuration("KafkaV2.Producer.BatchTimeout"),
		RequiredAcks: configYml.GetInt("KafkaV2.Producer.RequiredAcks"),
		Async:        configYml.GetBool("KafkaV2.Producer.Async"),
	}
}

// 生产消息,发送user信息
func writeMessage(ctx context.Context, w *kafka.Writer, msg kafka.Message) error {

	err := w.WriteMessages(ctx, msg)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("写入kafka失败，Tocic:%v,err:%v", w.Topic, err))
	}
	return err
}

func ProduceCommon(ctx context.Context, data *model.MessageStruct) {
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

func ProduceProcess(ctx context.Context, data *model.MessageStruct) {
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

func ProduceThreaten(ctx context.Context, data *model.MessageStruct) {
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

func ProduceAlarm(ctx context.Context, data *model.MessageStruct) {
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

func ProduceSyslog(ctx context.Context, data *model.MessageStruct) {
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
