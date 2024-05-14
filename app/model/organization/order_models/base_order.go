package order_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers"
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
	IsExpire    int
	CreateAtGte string
	CreateAtLte string
	ExpireAtGte string
	StatusList  []int
}

type baseOrder struct {
}

// conditions2Where 查询条件转化
func (crd *baseOrder) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !helpers.StringIsEmpty(search.OrderNo) {
		wheres["order_no"] = search.OrderNo
	}
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !helpers.StringIsEmpty(search.UserName) {
		wheres["user_name"] = search.UserName
	}

	if !helpers.StringIsEmpty(search.UserPhone) {
		wheres["user_phone"] = search.UserPhone
	}
	if !helpers.StringIsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !helpers.StringIsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}
	if !helpers.StringIsEmpty(search.OraName) {
		wheres["org_name like"] = "%" + search.OraName + "%"
	}
	if !helpers.StringIsEmpty(search.ExpireAtGte) {
		wheres["expire_at >="] = search.ExpireAtGte
	}
	if !helpers.StringIsEmpty(search.ProductName) {
		wheres["product_title like"] = "%" + search.ProductName + "%"
	}
	if search.GrantType > 0 && search.GrantType != enums.AllDefaultInt {
		wheres["grant_type"] = search.GrantType
	}
	if search.OrderStatus > 0 && search.OrderStatus != enums.AllDefaultInt {
		wheres["order_status"] = search.OrderStatus
	}
	if search.IsExpire != enums.AllDefaultInt {
		wheres["is_expire"] = search.IsExpire
	}
	if search.StatusList != nil && len(search.StatusList) > 0 {
		wheres["order_status in"] = search.StatusList
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
