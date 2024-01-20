package kafka_go

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"net"
	"strconv"
)

const (
	TopicCommon   string = "TopicCommon"
	TopicProcess  string = "TopicProcess"
	TopicThreaten string = "TopicThreaten"
	TopicAlarm    string = "TopicAlarm"
	TopicSyslog   string = "TopicSyslog"
)

func CreateTopic() {
	addr := variable.ConfigYml.GetString("KafkaV2.Addr")
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
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

	topics := []string{TopicCommon, TopicProcess, TopicThreaten, TopicAlarm, TopicSyslog}
	for _, val := range topics {
		name := variable.ConfigYml.GetString("KafkaV2." + val + ".Name")
		partition := variable.ConfigYml.GetInt("KafkaV2." + val + ".Partition")
		replicationFactor := variable.ConfigYml.GetInt("KafkaV2." + val + ".ReplicationFactor")

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
