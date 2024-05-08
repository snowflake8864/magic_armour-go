package org_stat_service

import (
	"goskeleton/app/model/model_organization"
	"goskeleton/app/utils/helpers"
	"strings"
)

func CreateStatDailyFactory() *StatDaily {
	return &StatDaily{masterModel: model_organization.CreateOrgStatDailyFactory()}
}

type StatDaily struct {
	masterModel *model_organization.OrgStatDailyModel
}

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	StatMonth   int
	StatDay     string
	CreateAtGte string
	CreateAtLte string
}

// search2Where 查询条件转化
func (crd *StatDaily) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !helpers.StringIsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if search.StatMonth > 0 {
		wheres["stat_month"] = search.StatMonth
	}
	if !helpers.StringIsEmpty(search.StatDay) {
		wheres["stat_day"] = search.StatDay
	}
	if !helpers.StringIsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !helpers.StringIsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}
	return wheres
}

// ColumnOfUser 统计表字段:用户相关
func (crd *StatDaily) ColumnOfUser(total, online, add, statDay bool) string {
	var columns []string
	if total {
		columns = append(columns, "SUM(user_total) AS user_total")
	}
	if online {
		columns = append(columns, "SUM(user_online) AS user_online")
	}
	if add {
		columns = append(columns, "SUM(user_add) AS user_add")
	}

	if statDay {
		columns = append(columns, "stat_day")
	}
	return strings.Join(columns, ", ")
}

// ColumnOfAgent 统计表字段: 终端相关
func (crd *StatDaily) ColumnOfAgent(total, online, add, linuxTotal, linuxOnline, linuxAdd, winTotal, winOnline, winAdd, statDay bool) string {
	var columns []string
	if total {
		columns = append(columns, "SUM(agent_total) AS agent_total")
	}
	if online {
		columns = append(columns, "SUM(agent_online) AS agent_online")
	}
	if add {
		columns = append(columns, "SUM(agent_add) AS agent_add")
	}
	if linuxTotal {
		columns = append(columns, "SUM(agent_linux_total) AS agent_linux_total")
	}
	if linuxOnline {
		columns = append(columns, "SUM(agent_linux_online) AS agent_linux_online")
	}
	if linuxAdd {
		columns = append(columns, "SUM(agent_linux_add) AS agent_linux_add")
	}
	if winTotal {
		columns = append(columns, "SUM(agent_win_total) AS agent_win_total")
	}
	if winOnline {
		columns = append(columns, "SUM(agent_win_online) AS agent_win_online")
	}
	if winAdd {
		columns = append(columns, "SUM(agent_win_add) AS agent_win_add")
	}

	if statDay {
		columns = append(columns, "stat_day")
	}
	return strings.Join(columns, ", ")
}

// Add 新增
func (crd *StatDaily) Add(param *model_organization.OrgStatDailyModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// SumByDay 计算： 每天
func (crd *StatDaily) SumByDay(day, column string) (result *model_organization.OrgStatDailyModel) {
	search := SearchConditions{
		StatDay: day,
	}
	where := crd.search2Where(search)
	result, err := crd.masterModel.CalcByConditions(where, column)
	if err != nil {
		return result
	}
	return result
}

// SumByMonthGroupDay 计算： 每月并且按照天分组
func (crd *StatDaily) SumByMonthGroupDay(month int, column string) (items []*model_organization.OrgStatDailyModel) {
	search := SearchConditions{
		StatMonth: month,
	}
	where := crd.search2Where(search)
	items, err := crd.masterModel.CalcByConditionsGroupDay(where, column)
	if err != nil {
		return nil
	}
	return items
}
