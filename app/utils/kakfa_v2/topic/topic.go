package topic

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
	"net"
	"strconv"
)

var configYml ymlconfig_interf.YmlConfigInterf

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	configYml = yml_config.CreateYamlFactory("kafka_v2")
}

func CreateTopic() {
	addr := configYml.GetString("KafkaV2.Producer.Addr")
	fmt.Println(addr)
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

	topics := []string{"TopicComm", "TopicProcess", "TopicThreaten", "TopicAlarm", "TopicSyslog"}
	for _, val := range topics {
		name := configYml.GetString("KafkaV2.TopicList." + val + ".Name")
		partition := configYml.GetInt("KafkaV2.TopicList." + val + ".Partition")
		replicationFactor := configYml.GetInt("KafkaV2.TopicList." + val + ".ReplicationFactor")

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
