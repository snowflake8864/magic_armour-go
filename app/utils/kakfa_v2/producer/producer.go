package producer

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
)

var configYml ymlconfig_interf.YmlConfigInterf

var writerCommon *kafka.Reader

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	fmt.Println("kakfa_v2 producer---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
}
