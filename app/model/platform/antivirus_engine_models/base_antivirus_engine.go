package antivirus_engine_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
)

type SearchConditions struct {
	Id           int64
	PlatformType int
	IsPublish    *int
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
	if condition.IsPublish != nil && *condition.IsPublish != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_publish",
			Val: *condition.IsPublish,
		})
	}
	return wheres
}
