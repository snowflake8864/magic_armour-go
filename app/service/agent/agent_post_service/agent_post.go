package agent_post_service

import (
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateAgentFactory() *AgentPost {
	return &AgentPost{}
}

type AgentPost struct {
}

// MakeId 生成iD
func (crd *AgentPost) MakeId(model port_models.PortWorkingModel) string {
	str := int_helper.Int64ToStr(model.OrgId)
	str += model.AgentId
	str += model.CreatedAt
	str += model.Port
	str += model.ProcDir
	str += int_helper.IntRandom2Str(1000, 9999)

	return md5_encrypt.MD5(str)
}

// MonitorToQuery 保存监控信息
func (crd *AgentPost) MonitorToQuery(orgId int64, agentId, inputJson string) error {
	var param agentMonitor
	if err := json_helper.ToStruct(inputJson, &param); err != nil {
		return err
	}
	now := date_helper.CurDatetime()

	return nil
}
