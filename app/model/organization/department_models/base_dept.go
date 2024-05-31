package department_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId    int64
	Id       int64
	Ids      []int64
	FullCode string
}

type baseDept struct {
}

// conditions2Where 查询条件转化
func (crd *baseDept) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
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
	if !string_helper.IsEmpty(condition.FullCode) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "full_code findInSetLater",
			Val: condition.FullCode,
		})
	}
	return wheres
}
