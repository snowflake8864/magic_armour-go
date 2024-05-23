package threaten_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	Id          string
	Ids         []string
	Dir         string
	GteCreateAt string
	LteCreateAt string
	Level       int
	LevelList   []int
}

type baseThreaten struct {
	Id          string
	Ids         []string
	Dir         string
	GteCreateAt string
	LteCreateAt string
	Level       int
	LevelList   []int
}

// conditions2Where 查询条件转化
func (m *baseThreaten) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !string_helper.IsEmpty(search.Id) {
		wheres["id"] = strings.ToLower(search.Id)
	}
	if search.Ids != nil && len(search.Ids) > 0 {
		wheres["id in"] = search.Ids
	}
	if !string_helper.IsEmpty(search.GteCreateAt) {
		wheres["created_at >="] = search.GteCreateAt
	}
	if !string_helper.IsEmpty(search.LteCreateAt) {
		wheres["created_at <="] = search.LteCreateAt
	}
	if search.Level > 0 && search.Level != enums.AllDefaultInt {
		wheres["level"] = search.Level
	}
	if !string_helper.IsEmpty(search.Dir) {
		wheres["dir like"] = "%" + search.Dir + "%"
	}
	if search.LevelList != nil && len(search.LevelList) > 0 {
		wheres["level in"] = search.LevelList
	}

	return wheres
}
