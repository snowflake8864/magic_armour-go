package test

import (
	"fmt"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestAgentLogin(t *testing.T) {
	var param = model_agent.AgentModel{}

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
	param.AgentStartAt = helpers.TimestampToDatetime(int64(1632465410))
	param.OsStartAt = helpers.TimestampToDatetime(int64(1632465310))
	param.Auth = "123123"
	param.AdminId = 1
	param.NetIp = "192.168.4.6"
	if helpers.StringIsEmpty(param.AgentIp) {
		param.AgentIp = param.NetIp
	}

	//兼容处理，参考bug 6181
	if helpers.StringIsEmpty(param.OsName) {
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
	isOk, data := agent_service.CreateAuthFactory().RegisterOrLogin(&param)
	fmt.Println(isOk)
	fmt.Println(data)
}
func TestAgentDecrypt(t *testing.T) {
	d, _ := agent_service.CreateCacheFactory().DecryptByLogin("c94a8b121a3686e25685d59c64b33015", true)
	fmt.Println(d.AgentId)
}

func TestGetMapsByAgents(t *testing.T) {

	agentCurdFact := agent_service.CreateCurdFactory()
	agentModel := agentCurdFact.GetItemByOrgAndAgent(1, "win2010_test_001", agentCurdFact.ColumnConfig(true))
	fmt.Println(helpers.JsonEncode(agentModel))
	fmt.Println("==========================")
	fmt.Println(agentModel.ProcessSwitch)

}

func TestSetLivingByAgent(t *testing.T) {
	agent_service.CreateCurdFactory().SetLivingByAgent("win2004_test_001")
}
