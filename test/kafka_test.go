package test

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/kafka_go"
	_ "goskeleton/bootstrap"
	"testing"
)

func TestConfig(t *testing.T) {
}

func TestCreateTopic(t *testing.T) {
	kafka_go.CreateTopic()
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
		Sign: kafka_go.SignAgent,
		Type: kafka_go.TypeAgentDelete,
		Data: "123333",
	}
	kafka_go.CreateProducerFactory()
	for i := 0; i <= 2; i++ {
		kafka_go.ProduceCommon(ctx, msg)
	}

}

func TestReadCommon(t *testing.T) {
	kafka_go.CreateConsumerFactory()
	kafka_go.ConsumeCommon()
}
