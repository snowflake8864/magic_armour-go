package usb_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	UsbHash     string
	UsbHashList []string
	Status      int
	IsOnline    int
}

type baseUsb struct {
}

// conditions2Where 查询条件转化
func (m *baseUsb) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id",
			Val: condition.AgentId,
		})
	}
	if !string_helper.IsEmpty(condition.UsbHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "usb_hash",
			Val: condition.UsbHash,
		})

	}
	if len(condition.UsbHashList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "usb_hash in",
			Val: condition.UsbHashList,
		})
	}
	if condition.Status > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status",
			Val: condition.Status,
		})
	}
	if condition.IsOnline != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_online",
			Val: condition.IsOnline,
		})
	}

	return wheres
}
