package consumer

import (
	"fmt"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
)

var configYml ymlconfig_interf.YmlConfigInterf

func init() {
	fmt.Println("kakfa_v2 init---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	
}
