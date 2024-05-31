package dir_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	SearchDir    string
	FileName     string
	FileType     int
	StartTimeGte string
	StartTimeLte string
	CreateAtGte  string
}

type baseDir struct {
}

// conditions2Where 查询条件转化
func (m *baseDir) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id",
			Val: condition.AgentId,
		})
	}

	if !string_helper.IsEmpty(condition.SearchDir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "search_dir",
			Val: condition.SearchDir,
		})
	}
	if !string_helper.IsEmpty(condition.FileName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "file_name like",
			Val: "%" + condition.FileName + "%",
		})
	}
	if condition.FileType != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "file_type",
			Val: condition.FileType,
		})
	}
	if !string_helper.IsEmpty(condition.StartTimeGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "start_time >=",
			Val: condition.StartTimeGte,
		})
	}
	if !string_helper.IsEmpty(condition.StartTimeLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "start_time <=",
			Val: condition.StartTimeLte,
		})
	}
	if !string_helper.IsEmpty(condition.CreateAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at >=",
			Val: condition.CreateAtGte,
		})
	}

	return wheres
}
