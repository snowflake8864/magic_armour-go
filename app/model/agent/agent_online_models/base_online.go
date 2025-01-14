package agent_online_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	HostName    string
	AgentIp     string
	NetIp       string
	IsOnline    *int
	Alias       string
	ContactName string
	DeptIds     []int64
	CreateAtGte string
	CreateAtLte string
}

type baseOnline struct {
}

// conditions2Where 查询条件转化
func (m *baseOnline) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.AgentId,
		})
	}
	if !string_helper.IsEmpty(condition.HostName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "host_name",
			Val: condition.HostName,
		})
	}
	if !string_helper.IsEmpty(condition.NetIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "net_ip",
			Val: condition.NetIp,
		})
	}
	if !string_helper.IsEmpty(condition.AgentIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_ip like",
			Val: "%" + condition.AgentIp + "%",
		})
	}
	if condition.IsOnline != nil && *condition.IsOnline != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_online",
			Val: *condition.IsOnline,
		})
	}

	if !string_helper.IsEmpty(condition.ContactName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "contact_name",
			Val: condition.ContactName,
		})
	}
	if !string_helper.IsEmpty(condition.Alias) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "alias",
			Val: condition.Alias,
		})
	}
	if condition.DeptIds != nil && len(condition.DeptIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dept_id in",
			Val: condition.DeptIds,
		})
	}
	if !string_helper.IsEmpty(condition.CreateAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at >=",
			Val: condition.CreateAtGte,
		})
	}
	if !string_helper.IsEmpty(condition.CreateAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at <=",
			Val: condition.CreateAtLte,
		})
	}
	return wheres
}
