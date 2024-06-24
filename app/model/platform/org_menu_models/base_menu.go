package org_menu_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
)

type SearchConditions struct {
	Id           int64
	Ids          []int64
	IsRoot       *int
	UsageVerList []int
}

type baseMenu struct {
}

// conditions2Where 查询条件转化
func (m *baseMenu) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}

	if len(condition.Ids) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id in",
			Val: condition.Ids,
		})
	}
	if condition.IsRoot != nil && *condition.IsRoot != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_root",
			Val: *condition.IsRoot,
		})
	}
	if len(condition.UsageVerList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "usage_ver in",
			Val: condition.UsageVerList,
		})
	}
	return wheres
}
