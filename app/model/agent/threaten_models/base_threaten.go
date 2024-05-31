package threaten_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	Id          string
	Ids         []string
	Dir         string
	CreateAtGte string
	CreateAtLte string
	Level       int
	LevelList   []int
}

type baseThreaten struct {
	Id          string
	Ids         []string
	Dir         string
	CreateAtGte string
	CreateAtLte string
	Level       int
	LevelList   []int
}

// conditions2Where 查询条件转化
func (m *baseThreaten) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if !string_helper.IsEmpty(condition.Id) {
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
	if condition.Level != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "level",
			Val: condition.Level,
		})
	}

	if len(condition.LevelList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "level in",
			Val: condition.LevelList,
		})
	}

	if !string_helper.IsEmpty(condition.Dir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dir like",
			Val: "%" + condition.Dir + "%",
		})
	}

	return wheres
}
