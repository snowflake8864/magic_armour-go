package test

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kakfa_v2"
	_ "goskeleton/bootstrap"
	"net"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(variable.ConfigKafkaV2Yml.GetString("KafkaV2.Producer.Addr"))

}

func TestCreateTopic(t *testing.T) {
	kakfa_v2.CreateTopicFactory()
}

func TestCreateTopic2(t *testing.T) {

	topic := "my-topic"

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
			Topic:             topic,
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
