package usb_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId    int64
	AgentId  string
	AgentIds []string
	Id       string
	IdList   []string
	UsbHash  string
	Status   *int
	IsOnline *int
}

// Calc 统计
type Calc struct {
	AgentId     string `json:"agent_id"`
	TotalNumber int    `json:"total_number"`
	WhiteNumber int    `json:"white_number"`
	BlackNumber int    `json:"black_number"`
	UnsetNumber int    `json:"unset_number"`
}

type baseUsb struct {
}

// ColumnCalc 字段：统计
func (m *baseUsb) ColumnCalc() string {
	var columns []string
	columns = append(columns, "agent_id AS agent_id")
	columns = append(columns, "COUNT(id) AS total_number")
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS unset_number", enum_usb.StatusDefault))
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS white_number", enum_usb.StatusWhite))
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS black_number", enum_usb.StatusBlack))

	return strings.Join(columns, ", ")
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
	if len(condition.AgentIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id in",
			Val: condition.AgentIds,
		})
	}
	if !string_helper.IsEmpty(condition.Id) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})

	}
	if len(condition.IdList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id in",
			Val: condition.IdList,
		})
	}
	if !string_helper.IsEmpty(condition.UsbHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "usb_hash",
			Val: condition.UsbHash,
		})

	}

	if condition.Status != nil && *condition.Status != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status",
			Val: *condition.Status,
		})
	}
	if condition.IsOnline != nil && *condition.IsOnline != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_online",
			Val: *condition.IsOnline,
		})
	}

	return wheres
}
