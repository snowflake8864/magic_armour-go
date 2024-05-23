package dir_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId         int64
	AgentId       string
	SearchDir     string
	FileName      string
	FileType      int
	GteCStartTime string
	LteStartTime  string
	GteCreateAt   string
}

type baseDir struct {
}

// conditions2Where 查询条件转化
func (m *baseDir) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.SearchDir) {
		wheres["search_dir"] = search.SearchDir
	}
	if !string_helper.IsEmpty(search.FileName) {
		wheres["file_name like"] = "%" + search.FileName + "%"
	}
	if search.FileType != enums.AllDefaultInt {
		wheres["file_type"] = search.FileType
	}
	if !string_helper.IsEmpty(search.GteCStartTime) {
		wheres["start_time >="] = search.GteCStartTime
	}
	if !string_helper.IsEmpty(search.LteStartTime) {
		wheres["start_time <="] = search.LteStartTime
	}
	if !string_helper.IsEmpty(search.GteCreateAt) {
		wheres["created_at >="] = search.GteCreateAt
	}

	return wheres
}
