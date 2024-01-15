package kakfa_v2

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"net"
	"strconv"
)

func CreateTopicFactory() {
	addr := variable.ConfigKafkaV2Yml.GetString("KafkaV2.Producer.Addr")
	fmt.Println(addr)
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
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

	var topicConfigs []kafka.TopicConfig

	topics := []string{"TopicComm", "TopicProcess", "TopicThreaten", "TopicAlarm", "TopicSyslog"}
	for _, val := range topics {
		name := variable.ConfigKafkaV2Yml.GetString("KafkaV2.TopicList." + val + ".Name")
		partition := variable.ConfigKafkaV2Yml.GetInt("KafkaV2.TopicList." + val + ".Partition")
		replicationFactor := variable.ConfigKafkaV2Yml.GetInt("KafkaV2.TopicList." + val + ".ReplicationFactor")

		st := kafka.TopicConfig{
			Topic:             name,
			NumPartitions:     partition,
			ReplicationFactor: replicationFactor,
		}
		topicConfigs = append(topicConfigs, st)
	}

	fmt.Println(topicConfigs)
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
