package org_stat_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
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
func (crd *baseStat) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if condition.StatMonth > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "stat_month",
			Val: condition.StatMonth,
		})
	}
	if !string_helper.IsEmpty(condition.StatDay) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "stat_day",
			Val: condition.StatDay,
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
