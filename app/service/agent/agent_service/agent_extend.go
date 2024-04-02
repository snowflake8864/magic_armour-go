package agent_service

import (
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers"
)

func CreateExtendFactory() *AgentExtend {
	return &AgentExtend{masterModel: model_agent.CreateAgentExtendFactory()}
}

type AgentExtend struct {
	AgentBase
	masterModel *model_agent.AgentExtendModel
}

// AddDefault 存在登录，不存在注册
func (crd *AgentExtend) AddDefault(agentId string) bool {

	//新终端默认使用全局配置
	dictIds := []string{
		enum_dictionary.AgentGlobalPlugPortConfig,
		enum_dictionary.AgentGlobalPlugLoginConfig,
		enum_dictionary.AgentGlobalLoginSolveConfig,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
	var isExist bool
	jsonPlugPort, isExist := mapDict[enum_dictionary.AgentGlobalPlugPortConfig]
	if !isExist || helpers.StringIsEmpty(jsonPlugPort) {
		plugPortConfig := enum_plug.PlugPortConfig{}
		plugPortConfig.Format()
		jsonPlugPort = helpers.JsonEncode(plugPortConfig)
	}
	jsonPlugLogin, isExist := mapDict[enum_dictionary.AgentGlobalPlugLoginConfig]
	if !isExist || helpers.StringIsEmpty(jsonPlugLogin) {
		plugLoginConfig := enum_plug.PlugLoginConfig{}
		plugLoginConfig.Format()
		jsonPlugLogin = helpers.JsonEncode(plugLoginConfig)
	}
	jsonLoginSolve, isExist := mapDict[enum_dictionary.AgentGlobalPlugLoginConfig]
	if !isExist || helpers.StringIsEmpty(jsonPlugLogin) {
		loginSolveConfig := enum_login_solve.LoginSolveConfig{}
		loginSolveConfig.Format()
		jsonLoginSolve = helpers.JsonEncode(loginSolveConfig)
	}

	var param = model_agent.AgentExtendModel{}
	param.AgentId = agentId
	param.PlugPortConfig = jsonPlugPort
	param.PlugLoginConfig = jsonPlugLogin
	param.LoginSolveConfig = jsonLoginSolve
	isOk := crd.masterModel.InsertItem(&param)
	if !isOk {
		return false
	}
	return true
}
