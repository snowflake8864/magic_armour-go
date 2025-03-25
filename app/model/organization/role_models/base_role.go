package role_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	Id          int64
	OrgId       int64
	Name        string
	CreateAtGte string
	CreateAtLte string
}

type baseRole struct {
}

// conditions2Where 查询条件转化
func (crd *baseRole) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if condition.OrgIds != nil && len(condition.OrgIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id in",
			Val: condition.OrgIds,
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
	if !string_helper.IsEmpty(condition.OrgName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_name like",
			Val: "%" + condition.OrgName + "%",
		})
	}

	if !string_helper.IsEmpty(condition.ProductTitle) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "product_title like",
			Val: "%" + condition.ProductTitle + "%",
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
	if !string_helper.IsEmpty(condition.PayOrderNo) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "pay_order_no",
			Val: condition.PayOrderNo,
		})
	}

	return wheres
}

// ColumnCalcOfUsage 字段： 统计使用数
func (crd *baseRole) ColumnCalcOfUsage() string {
	var columns []string
	columns = append(columns, "SUM(buy_number) AS buy_number")
	columns = append(columns, "SUM(usage_number) AS usage_number")
	return strings.Join(columns, ", ")
}

type RespRemind struct {
	TotalNumber     int64 `json:"total_number"`
	UnexpiredNumber int64 `json:"unexpired_number"`
	GoldNumber      int64 `json:"gold_number"`
}

// ColumnCalcOfRemind 字段： 统计提醒数
func (crd *baseRole) ColumnCalcOfRemind() string {
	var columns []string
	columns = append(columns, "COUNT(*) as total_number")
	columns = append(columns, fmt.Sprintf("SUM(if(is_expire=%d, 1,0)) as unexpired_number", enums.IsNot))
	columns = append(columns, fmt.Sprintf("SUM(if(grant_type=%d, 1,0)) as gold_number", enum_product.TypeGold))
	return strings.Join(columns, ", ")
}
