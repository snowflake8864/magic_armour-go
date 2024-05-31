package order_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrderNo     string
	OrgId       int64
	OraName     string
	UserName    string
	UserPhone   string
	ProductName string
	GrantType   int
	OrderStatus int
	IsExpire    *int
	CreateAtGte string
	CreateAtLte string
	ExpireAtLte string
	StatusList  []int
}

type baseOrder struct {
}

// conditions2Where 查询条件转化
func (crd *baseOrder) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if !string_helper.IsEmpty(condition.OrderNo) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "order_no",
			Val: condition.OrderNo,
		})
	}
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
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
			Val: condition.UserName,
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
	if !string_helper.IsEmpty(condition.ExpireAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "expire_at <=",
			Val: condition.ExpireAtLte,
		})
	}
	if !string_helper.IsEmpty(condition.OraName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_name like",
			Val: "%" + condition.OraName + "%",
		})
	}

	if !string_helper.IsEmpty(condition.ProductName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "product_title like",
			Val: "%" + condition.ProductName + "%",
		})
	}
	if condition.GrantType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "grant_type",
			Val: condition.GrantType,
		})
	}
	if condition.OrderStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "order_status",
			Val: condition.OrderStatus,
		})
	}
	if condition.StatusList != nil && len(condition.StatusList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "order_status in",
			Val: condition.StatusList,
		})
	}
	if condition.IsExpire != nil && *condition.IsExpire != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_expire",
			Val: *condition.IsExpire,
		})
	}

	return wheres
}

// ColumnCalcOfUsage 字段： 统计使用数
func (crd *baseOrder) ColumnCalcOfUsage() string {
	var columns []string
	columns = append(columns, "SUM(buy_number) AS buy_number")
	columns = append(columns, "SUM(usage_number) AS usage_number")
	return strings.Join(columns, ", ")
}
