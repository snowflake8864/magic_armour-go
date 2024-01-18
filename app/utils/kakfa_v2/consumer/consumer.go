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
var configReader kafka.ReaderConfig

var readerCommon *kafka.Reader
var readerProcess *kafka.Reader
var readerThreaten *kafka.Reader
var readerAlarm *kafka.Reader
var readerSyslog *kafka.Reader

func init() {
	fmt.Println("kakfa_v2 consumer---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	getReaderConfig()
	readerCommon = newConsumer("TopicComm")
	readerProcess = newConsumer("TopicSyslog")
	readerThreaten = newConsumer("TopicThreaten")
	readerAlarm = newConsumer("TopicAlarm")
	readerSyslog = newConsumer("TopicSyslog")
}

func getReaderConfig() {
	addr := configYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")
	groupId := configYml.GetString("KafkaV2.Consumer.GroupId")
	commitInterval := configYml.GetDuration("KafkaV2.Consumer.CommitInterval")
	startOffset := configYml.GetInt64("KafkaV2.Consumer.StartOffset")

	configReader = kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupId,
		CommitInterval: commitInterval,
		StartOffset:    startOffset, //kafka.FirstOffset
		MinBytes:       0,           // 最小读取字节数
		MaxBytes:       10e6,        // 最大读取字节数
	}
}

func newConsumer(name string) *kafka.Reader {
	topic := configYml.GetString("KafkaV2.TopicList." + name + ".Name")
	//partition := configYml.GetInt("KafkaV2.TopicList." + name + ".Partition")

	// 配置消费者
	configReader.Topic = topic
	//configReader.Partition = partition

	return kafka.NewReader(configReader)
}

// 消费消息
func readMessage(ctx context.Context, r *kafka.Reader) (*model.MessageStruct, error) {
	result := &model.MessageStruct{}

	msg, err := r.ReadMessage(ctx)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("topic:%s, 读kafka失败，err:%v", r.Config().Topic, err))
		return result, err
	} else {
		err = json.Unmarshal(msg.Value, result)
		if err != nil {
			variable.ZapLog.Error(fmt.Sprintf("topic:%s, json unmarshal msg value err，msg:%v,err:%v", r.Config().Topic, result, err))
		}
		//fmt.Println(fmt.Sprintf("topic=%s,partition=%d,offset=%d,key=%s,result=%v", msg.Topic, msg.Partition, msg.Offset, msg.Key, result))
		return result, err

	}
}

// 消费主题： Common

func ConsumeCommon() {
	for {
		s, err := readMessage(context.Background(), readerCommon)
		if err != nil {
			break // 当前消息读取失败时，并不退出for终止所有后续消费，而是跳过该消息即可
		}
		variable.ZapLog.Info(fmt.Sprintf("ReadCommon result:%v", s))
	}
}
