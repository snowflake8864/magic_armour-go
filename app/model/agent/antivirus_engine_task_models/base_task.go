package antivirus_engine_task_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId         int64
	Id            int64
	Ids           []int64
	Title         string
	UpgradeMode   int
	TargetType    int
	TargetAgentId string
	TargetDeptId  int64

	ExecStatus     int
	ExecStatusList []int
	ExecNextTsLte  int64
	ExecLiveTsLte  int64

	TaskId       int64
	CreateAtGte  string
	CreateAtLte  string
	UpdatedAtGte string
	UpdatedAtLte string
}

type baseOnline struct {
}

// conditions2Where 查询条件转化
func (m *baseOnline) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}

	if condition.UpgradeMode > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "upgrade_mode",
			Val: condition.UpgradeMode,
		})
	}
	if condition.TargetType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "target_type",
			Val: condition.TargetType,
		})
	}
	if condition.ExecStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "exec_status",
			Val: condition.ExecStatus,
		})
	}
	if len(condition.ExecStatusList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "exec_status in",
			Val: condition.ExecStatusList,
		})
	}

	if condition.ExecNextTsLte > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "exec_next_ts <=",
			Val: condition.ExecNextTsLte,
		})
	}

	if condition.ExecLiveTsLte > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "exec_live_ts <=",
			Val: condition.ExecLiveTsLte,
		})
	}

	if !string_helper.IsEmpty(condition.Title) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "title LIKE",
			Val: "%" + condition.Title + "%",
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
	if !string_helper.IsEmpty(condition.UpdatedAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "updated_at >=",
			Val: condition.UpdatedAtGte,
		})
	}
	if !string_helper.IsEmpty(condition.UpdatedAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "updated_at <=",
			Val: condition.UpdatedAtLte,
		})
	}

	return wheres
}
