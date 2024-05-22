package agent_post_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/service/agent/port_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/md5_encrypt"
)

func CreatePortFactory() *PortPost {
	return &PortPost{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type PortPost struct {
	queryClient *producer_kafka.ProducerClient
}

// MakeId 生成iD
func (crd *PortPost) MakeId(model port_models.PortWorkingModel) string {
	str := int_helper.Int64ToStr(model.OrgId)
	str += model.AgentId
	str += model.CreatedAt
	str += model.Port
	str += model.ProcDir
	str += int_helper.IntRandom2Str(1000, 9999)

	return md5_encrypt.MD5(str)
}

// ToQuery 格式化参数后推送给队列
func (crd *PortPost) ToQuery(orgId int64, agentId, inputJson string) error {
	var params []portWorking
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()

	var models []port_models.PortWorkingModel
	for _, item := range params {
		model := port_models.PortWorkingModel{
			OrgId:     orgId,
			AgentId:   agentId,
			Type:      item.Type,
			Ip:        item.Ip,
			Port:      json_helper.Number2String(item.Port),
			ProcId:    json_helper.Number2Int(item.Pid),
			ProcDir:   path_helper.ToStandard(item.Process),
			CreatedAt: now,
		}
		model.Id = crd.MakeId(model)

		models = append(models, model)
	}
	if !port_service.CreateWorkingFactory().BatchAdd(orgId, agentId, &models) {
		variable.ZapLog.Sugar().Debug("保存失败")
	}

	return nil
}
