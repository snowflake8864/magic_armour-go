package org_stat_models

import (
	"goskeleton/app/utils/helpers"
	"strings"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	StatMonth   int
	StatDay     string
	CreateAtGte string
	CreateAtLte string
}

type baseStat struct {
	OrgId       int64
	AgentId     string
	StatMonth   int
	StatDay     string
	CreateAtGte string
	CreateAtLte string
}

// conditions2Where 查询条件转化
func (crd *baseStat) conditions2Where(search SearchConditions) map[string]interface{} {
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
func (crd *baseStat) ColumnOfUser(total, online, add, statDay bool) string {
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
func (crd *baseStat) ColumnOfAgent(total, online, add, linuxTotal, linuxOnline, linuxAdd, winTotal, winOnline, winAdd, statDay bool) string {
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
