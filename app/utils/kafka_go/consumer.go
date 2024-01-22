package kafka_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"strings"
)

var readerCommon *kafka.Reader
var readerProcess *kafka.Reader
var readerThreaten *kafka.Reader
var readerAlarm *kafka.Reader
var readerSyslog *kafka.Reader

func CreateConsumerFactory() {
	config := getReaderConfig()

	readerCommon = newConsumer(TopicCommon, config)
}

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

func readMessage(ctx context.Context, r *kafka.Reader) (*MessageStruct, error) {
	result := &MessageStruct{}

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
