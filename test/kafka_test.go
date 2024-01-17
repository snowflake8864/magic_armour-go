package test

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/kakfa_v2/consumer"
	"goskeleton/app/utils/kakfa_v2/model"
	"goskeleton/app/utils/kakfa_v2/producer"
	"goskeleton/app/utils/kakfa_v2/topic"
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
	msg := model.MessageStruct{
		Type: "testt",
		Data: "fsdfs",
	}

	p := producer.GetOneKafkaProducer()
	p.PushCommon(msg)

}

func TestReadCommon(t *testing.T) {
	consumer.ReadCommon()
}
