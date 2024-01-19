package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/kafka_v2/model"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
	"strings"
	"time"
)

var configYml ymlconfig_interf.YmlConfigInterf
var saramaConfig = sarama.NewConfig()

var producerCommon *sarama.SyncProducer

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	fmt.Println("kafka_sarama producer---------------------------------------------------------")
	getConfig()

	addr := configYml.GetString("KafkaV2.Addr")
	brokers := strings.Split(addr, ",")

	producerCommon = newProducer(brokers)
}

func getConfig() {
	// 设置生产者发布消息时需要等待的确认级别
	saramaConfig.Producer.RequiredAcks = sarama.NoResponse
	// 设置等待服务器响应的最长时间
	saramaConfig.Producer.Timeout = 5 * time.Second
	// 设置分区选择器, 可以根据需要自定义分区选择器
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.Return.Successes = true
}

func newProducer(brokers []string) *sarama.SyncProducer {
	saramaProducer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("NewSyncProducer失败，err:%v", err))
	}
	return &saramaProducer
}

// 生产消息,发送user信息
func produceMessage(p sarama.SyncProducer, msg *sarama.ProducerMessage) error {
	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		// 错误处理逻辑
		variable.ZapLog.Error(fmt.Sprintf("produceMessage fail，Tocic:%v,err:%v", msg.Topic, err))
		return err
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

func ProduceCommon(data *model.MessageStruct) {
	msgContent, err := json.Marshal(&data)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf(
			"ProduceCommon json marshal user err，data:%v,err:%v", data, err))
	}
	//variable.ZapLog.Info(string(msgContent))

	message := &sarama.ProducerMessage{
		Topic: "common",
		Value: sarama.StringEncoder(msgContent),
	}

	err = produceMessage(*producerCommon, message)
	if err != nil {
		fmt.Println(fmt.Sprintf("ProduceCommon 写入kafka失败，data:%v ,err:%v", data, err))
	}
}
