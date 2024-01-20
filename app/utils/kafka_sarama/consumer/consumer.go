package consumer

import (
	"fmt"
	"github.com/IBM/sarama"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
)

var configYml ymlconfig_interf.YmlConfigInterf
var saramaConfig = sarama.NewConfig()

func init() {
	fmt.Println("kakfa_v2 consumer---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	saramaConfig.Version
	sarama.NewConsumerGroupFromClient()
}
