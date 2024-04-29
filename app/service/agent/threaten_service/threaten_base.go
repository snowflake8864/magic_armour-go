package threaten_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers"
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

type ThreatenBase struct {
}

// search2Where 查询条件转化
func (crd *ThreatenBase) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !helpers.StringIsEmpty(search.Id) {
		wheres["id"] = search.Id
	}
	if search.Ids != nil && len(search.Ids) > 0 {
		wheres["id in"] = search.Ids
	}
	if !helpers.StringIsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !helpers.StringIsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}
	if search.Level > 0 && search.Level != enums.AllDefaultInt {
		wheres["level"] = search.Level
	}
	if !helpers.StringIsEmpty(search.Dir) {
		wheres["dir like"] = "%" + search.Dir + "%"
	}
	if search.LevelList != nil && len(search.LevelList) > 0 {
		wheres["level in"] = search.LevelList
	}

	return wheres
}
