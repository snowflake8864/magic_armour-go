package prodcut_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	Id        int64
	GrantType int
	Title     string
	IsPublish *int
}
type baseProduct struct {
}

// search2Where 查询条件转化
func (crd *baseProduct) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {

	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if condition.GrantType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "grant_type",
			Val: condition.GrantType,
		})
	}
	if !string_helper.IsEmpty(condition.Title) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "title like",
			Val: "%" + condition.Title + "%",
		})
	}
	if condition.IsPublish != nil && *condition.IsPublish != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_publish",
			Val: *condition.IsPublish,
		})
	}
	return wheres
}
