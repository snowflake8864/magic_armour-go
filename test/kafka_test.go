package test

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/consumer_kafka"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/kafka_go/topic_kakfa"
	_ "goskeleton/bootstrap"
	"testing"
)

func TestConfig(t *testing.T) {
}

func TestCreateTopic(t *testing.T) {
	topic_kakfa.CreateTopicFactory().Init()
}

func TestShowTopic(t *testing.T) {
	conn, err := kafka.Dial("tcp", "192.168.16.251:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
}

func TestDeleteTopic(t *testing.T) {
	conn, err := kafka.Dial("tcp", "192.168.16.251:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	err = conn.DeleteTopics("my-topic", "alarm")
	if err != nil {
		panic(err.Error())
	}
}

func TestWriteCommon(t *testing.T) {
	ctx := context.Background()

	msg := &kafka_go.MessageStruct{
		Type:    kafka_go.TypeCommonAgentDelete,
		AgentId: "sssss1",
		OrgId:   2,
	}
	producerFact := producer_kafka.CreateProducerFactory()
	for i := 0; i <= 2; i++ {
		producerFact.ProduceCommon(ctx, msg)
	}

}

func TestReadCommon(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeCommon()
}

func TestReadProcess(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeProcess()
}

func TestReadDll(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeDll()
}

func TestReadAlarm(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeAlarm()

}

func TestReadSyslog(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeSyslog()

}

func TestReadThreaten(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeThreaten()

}

func TestConsumeGroup(t *testing.T) {
	consumer_kafka.CreateConsumerFactory().ConsumeProcess()

}
