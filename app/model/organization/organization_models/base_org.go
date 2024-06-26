package organization_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	OrgType    int
	OraName    string
	LogOraName string
	UserId     int64
	UserName   string
	UserPhone  string
	StatusList []int
}

type baseOrg struct {
}

// search2Where 查询条件转化
func (crd *baseOrg) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if condition.OrgType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: condition.OrgType,
		})
	}
	if condition.UserId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "user_id",
			Val: condition.UserId,
		})
	}
	if !string_helper.IsEmpty(condition.UserName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "user_name",
			Val: condition.UserName,
		})
	}
	if !string_helper.IsEmpty(condition.UserPhone) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "user_phone",
			Val: condition.UserPhone,
		})
	}

	if !string_helper.IsEmpty(condition.OraName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "name like",
			Val: "%" + condition.OraName + "%",
		})
	}
	if !string_helper.IsEmpty(condition.LogOraName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_name like",
			Val: "%" + condition.LogOraName + "%",
		})
	}

	if condition.StatusList != nil && len(condition.StatusList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status in",
			Val: condition.StatusList,
		})
	}
	return wheres
}
