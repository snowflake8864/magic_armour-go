package test

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/kakfa_v2/producer"
	"goskeleton/app/utils/kakfa_v2/topic"
	_ "goskeleton/bootstrap"
	"net"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
}

func TestCreateTopic(t *testing.T) {
	topic.CreateTopic()
}

func TestCreateTopic2(t *testing.T) {

	topic1 := "my-topic"

	conn, err := kafka.Dial("tcp", "192.168.17.251:9092")
	if err != nil {
		fmt.Println("1111111111111111111111111111111111111")
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	fmt.Println(controller.Host)
	fmt.Println(controller.Port)
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic1,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
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
	msg := producer.ProducerMessage{
		Type: "testt",
		Data: "fsdfs",
	}

	p := producer.GetOneKafkaProducer()
	p.PushCommon(msg)

}
