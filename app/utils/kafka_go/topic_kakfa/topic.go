package topic_kakfa

import (
	"github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"net"
	"strconv"
)

func CreateTopicFactory() *TopicClient {
	return &TopicClient{
		Addr: variable.ConfigYml.GetString("KafkaV2.Addr"),
		topicList: []string{
			kafka_go.TopicCommon,
			kafka_go.TopicProcess,
			kafka_go.TopicThreaten,
			kafka_go.TopicAlarm,
			kafka_go.TopicSyslog,
			kafka_go.TopicCrontab,
		},
	}

}

// 结构体
type TopicClient struct {
	Addr      string
	topicList []string
}

func (k *TopicClient) Init() {
	conn, err := kafka.Dial("tcp", k.Addr)
	if err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	variable.ZapLog.Info(json_helper.Encode(controller))
	if err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
	defer controllerConn.Close()

	var topicConfigs []kafka.TopicConfig
	for _, val := range k.topicList {
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

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
