package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kakfa_v2/model"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
	"strings"
)

var configYml ymlconfig_interf.YmlConfigInterf

var consumerConfig kafka.ReaderConfig

var consumerCommon *kafka.Reader
var consumerProcess *kafka.Reader
var consumerThreaten *kafka.Reader
var consumerAlarm *kafka.Reader
var consumerSyslog *kafka.Reader

func init() {
	fmt.Println("kakfa_v2 consumer---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	getConsumerConfig()
	consumerCommon = newConsumer("TopicComm")
	consumerProcess = newConsumer("TopicSyslog")
	consumerThreaten = newConsumer("TopicThreaten")
	consumerAlarm = newConsumer("TopicAlarm")
	consumerSyslog = newConsumer("TopicSyslog")
}

func getConsumerConfig() {
	addr := configYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")
	//groupId := configYml.GetString("KafkaV2.Consumer.GroupId")
	commitInterval := configYml.GetDuration("KafkaV2.Consumer.CommitInterval")
	startOffset := configYml.GetInt64("KafkaV2.Consumer.StartOffset")

	consumerConfig = kafka.ReaderConfig{
		Brokers: brokers,
		//GroupID:        groupId,
		CommitInterval: commitInterval,
		StartOffset:    startOffset, //kafka.FirstOffset
		MinBytes:       0,           // 最小读取字节数
		MaxBytes:       10e6,        // 最大读取字节数
	}
}

func newConsumer(name string) *kafka.Reader {
	topic := configYml.GetString("KafkaV2.TopicList." + name + ".Name")
	partition := configYml.GetInt("KafkaV2.TopicList." + name + ".Partition")

	// 配置消费者
	consumerConfig.Topic = topic
	consumerConfig.Partition = partition

	return kafka.NewReader(consumerConfig)
}

// 消费消息
func readMessage(ctx context.Context, consumer *kafka.Reader) (*model.MessageStruct, error) {
	result := &model.MessageStruct{}

	msg, err := consumer.ReadMessage(ctx)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("topic:%s, 读kafka失败，err:%v", consumer.Config().Topic, err))
		return result, err
	} else {
		err = json.Unmarshal(msg.Value, result)
		if err != nil {
			fmt.Println(fmt.Sprintf("topic:%s, json unmarshal msg value err，msg:%v,err:%v", consumer.Config().Topic, result, err))
		}
		return result, err
		//fmt.Println(fmt.Sprintf("topic=%s,partition=%d,offset=%d,key=%s,result=%v", msg.Topic, msg.Partition, msg.Offset, msg.Key, result))
	}
}

// 消费主题： Common

func ReadCommon() {
	for {
		s, err := readMessage(context.Background(), consumerCommon)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}
		variable.ZapLog.Info(fmt.Sprintf("ReadCommon result:%v", s))
	}
}
