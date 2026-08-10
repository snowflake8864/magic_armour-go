package antivirus_engine_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	Id           int64
	PlatformType int
	IsNewest     *int
	IsPublish    *int
	Ver          string
	CreateAtGte  string
	CreateAtLte  string
}
type baseEngine struct {
}

// search2Where 查询条件转化
func (crd *baseEngine) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {

	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if condition.PlatformType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "platform_type",
			Val: condition.PlatformType,
		})
	}
	if !string_helper.IsEmpty(condition.Ver) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "ver",
			Val: condition.Ver,
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
	if condition.IsPublish != nil && *condition.IsPublish != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_publish",
			Val: *condition.IsPublish,
		})
	}
	if condition.IsNewest != nil && *condition.IsPublish != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_newest",
			Val: *condition.IsNewest,
		})
	}
	return wheres
}
