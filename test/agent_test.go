package test

import (
	"fmt"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/agent_cmd"
	"testing"
)

func TestAgentLogin(t *testing.T) {
	var param = agent_models.AgentModel{}

	param.AgentId = "test_111"
	param.HostName = "中文名称"
	param.MacId = "E0:BE:03:42:A2:72,"
	param.AgentIp = "192.168.44.66"
	param.Ver = "2.0.0.20"
	param.PlatformType = 2
	param.OsName = ""
	param.MemDesc = "16g"
	param.CpuDesc = "16u"
	param.HdDesc = "931g"
	param.AgentStartAt = date_helper.TimestampToDatetime(int64(1632465410))
	param.OsStartAt = date_helper.TimestampToDatetime(int64(1632465310))
	param.Auth = "123123"
	param.OrgId = 1
	param.UserId = 1
	param.NetIp = "192.168.4.6"
	if string_helper.IsEmpty(param.AgentIp) {
		param.AgentIp = param.NetIp
	}

	//兼容处理，参考bug 6181
	if string_helper.IsEmpty(param.OsName) {
		if param.PlatformType == enum_software.PlatformTypeWindows {
			param.OsName = "windows"
		} else {
			param.OsName = "centos"
		}
	}

	//系统类型
	param.OsType = enum_software.GetOsTypeByName(param.OsName)
	if param.OsType == 0 {
		fmt.Println("未知的系统类型")
		return
	}
	fmt.Println("param---------------------------------------------")
	//data, _ := json.Marshal(param)
	fmt.Println(param)
	fmt.Println("RegisterOrLogin---------------------------------------------")
	isOk, data := agent_service.CreateOnlineFactory().RegisterOrLogin(&param)
	fmt.Println(isOk)
	fmt.Println(data)
}
func TestAgentDecrypt(t *testing.T) {
	d, _ := agent_service.CreateCacheFactory().DecryptByLogin("c94a8b121a3686e25685d59c64b33015", true)
	fmt.Println(d.AgentId)
}

func TestAgentGrantExpired(t *testing.T) {
	agent_cmd.CreateAgentFactory().GrantExpired()
}

func TestAgentOffline(t *testing.T) {
	agent_cmd.CreateAgentFactory().Offline()
}

func TestAgentMonitorToQuery(t *testing.T) {
	inputJson := "{\"hd_size\":\"147315MB\",\"hd_usage\":\"18\",\"cpu_number\":\"4\",\"cpu_usage\":\"1.50981\",\"cpu_tops\":[{\"id\":\"1155\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"11.9172\",\"mem_size\":\"1529044KB\",\"user\":\"elsearch\"},{\"id\":\"1154\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"3.7617\",\"mem_size\":\"1075740KB\",\"user\":\"logstash\"},{\"id\":\"1159\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"1.85666\",\"mem_size\":\"2360112KB\",\"user\":\"kafka\"},{\"id\":\"19109\",\"dir\":\"/opt/osec/MagicArmor_0\",\"hash\":\"6c0b9df8b7370cacc2728d0c1eff7cd1\",\"cpu_usage\":\"0.13308\",\"mem_size\":\"7268KB\",\"user\":\"root\"},{\"id\":\"1160\",\"dir\":\"/usr/sbin/zabbix_agent2\",\"hash\":\"82dc49046011e8e0a39d9f603a8f2f34\",\"cpu_usage\":\"0.0972695\",\"mem_size\":\"15628KB\",\"user\":\"zabbix\"}],\"mem_size\":\"7970416KB\",\"mem_usage\":\"69\",\"mem_tops\":[{\"id\":\"1159\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"1.85666\",\"mem_size\":\"2360112KB\",\"user\":\"kafka\"},{\"id\":\"1155\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"11.9172\",\"mem_size\":\"1529044KB\",\"user\":\"elsearch\"},{\"id\":\"1154\",\"dir\":\"/kingdee/cosmic/jdk/jdk1.8.0_201/bin/java\",\"hash\":\"7ce2c372228aa741a17f92cdc10c6883\",\"cpu_usage\":\"3.7617\",\"mem_size\":\"1075740KB\",\"user\":\"logstash\"},{\"id\":\"1715\",\"dir\":\"/CloudResetPwdUpdateAgent/depend/jre1.8.0_131/bin/java\",\"hash\":\"a4520ff15f0b61db0b6c566274e03d05\",\"cpu_usage\":\"0.0817988\",\"mem_size\":\"69496KB\",\"user\":\"root\"},{\"id\":\"758\",\"dir\":\"/usr/bin/python2.7\",\"hash\":\"d9ba93fc1778448ad966b5f613806caf\",\"cpu_usage\":\"0.00196214\",\"mem_size\":\"28192KB\",\"user\":\"root\"}],\"self\":{\"mem_size\":\"3980KB\",\"cpu_usage\":\"0\"}}"
	err := agent_post_service.CreateAgentFactory().MonitorToQuery(1, "ddddddd", inputJson)
	if err != nil {
		fmt.Println()
		return
	}
}
