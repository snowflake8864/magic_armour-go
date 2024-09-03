package outreach_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	Id         int64
	Ip         string
	IpLike     string
	Addr       string
	Type       int
	OpenSwitch *int
}

type baseOutreach struct {
}

// conditions2Where 查询条件转化
func (m *baseOutreach) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if condition.Type > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: condition.Type,
		})
	}
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if !string_helper.IsEmpty(condition.Ip) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "ip",
			Val: condition.Ip,
		})
	}

	if !string_helper.IsEmpty(condition.IpLike) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "ip like",
			Val: condition.IpLike + "%",
		})
	}

	if condition.OpenSwitch != nil && *condition.OpenSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "open_switch",
			Val: *condition.OpenSwitch,
		})
	}

	if !string_helper.IsEmpty(condition.Addr) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "addr",
			Val: condition.Addr,
		})
	}
	return wheres
}
