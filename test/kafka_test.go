package test

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/kafka_v2/consumer"
	"goskeleton/app/utils/kafka_v2/model"
	"goskeleton/app/utils/kafka_v2/producer"
	"goskeleton/app/utils/kafka_v2/topic"
	_ "goskeleton/bootstrap"
	"testing"
)

func TestConfig(t *testing.T) {
}

func TestCreateTopic(t *testing.T) {
	topic.CreateTopic()
}

func TestShowTopic(t *testing.T) {
	conn, err := kafka.Dial("tcp", "192.168.17.251:9092")
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
	conn, err := kafka.Dial("tcp", "192.168.17.251:9092")
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

	msg := &model.MessageStruct{
		Sign: "agent",
		Type: "safe",
		Data: "data",
	}
	for i := 0; i <= 100; i++ {
		producer.ProduceCommon(ctx, msg)
	}

}

func TestReadCommon(t *testing.T) {
	consumer.ConsumeCommon()
}
