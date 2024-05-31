package org_menu_models

import (
	"goskeleton/app/model/model_comm"
)

type SearchConditions struct {
	Id     int64
	Ids    []int64
	IsRoot int
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
	if condition.IsRoot > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_root",
			Val: condition.IsRoot,
		})
	}
	return wheres
}
