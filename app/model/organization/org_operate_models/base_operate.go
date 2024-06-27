package org_operate_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId       int64
	CreateId    int64
	Module      string
	NetIp       string
	CreateAtGte string
	CreateAtLte string
}

type baseOperate struct {
}

// conditions2Where 查询条件转化
func (m *baseOperate) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.Module) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "module",
			Val: condition.Module,
		})
	}
	if !string_helper.IsEmpty(condition.NetIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "net_ip",
			Val: condition.NetIp,
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
	if condition.CreateId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "create_id",
			Val: condition.CreateId,
		})
	}

	return wheres
}
