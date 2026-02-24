package alarm_stat_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	AgentIds   []string
	StatMonth  int
	StatDay    string
	StatDayGte string
	StatDayLte string
}

const (
	PrefixTotal     string = ""
	PrefixProcess   string = "process"
	PrefixDll       string = "dll"
	PrefixBlackmail string = "blackmail"
	PrefixTamper    string = "tamper"
	PrefixPort      string = "port"
	PrefixPlug      string = "plug"
	PrefixUsb       string = "usb"
	PrefixOutreach  string = "outreach"
	PrefixLogin     string = "login"
	PrefixSelf      string = "self"
	PrefixVirus     string = "virus"
)

type baseAlarm struct {
}

// search2Where 查询条件转化
func (crd *baseAlarm) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if len(condition.AgentIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id in",
			Val: condition.AgentIds,
		})
	}
	if !string_helper.IsEmpty(condition.StatDayGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "stat_day >=",
			Val: condition.StatDayGte,
		})
	}
	if !string_helper.IsEmpty(condition.StatDayLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "stat_day <=",
			Val: condition.StatDayLte,
		})
	}
	return wheres
}

// StatColumn 统计表字段
// @use: StatColumn(PrefixTotal, true, true,true,true)
func (crd *baseAlarm) StatColumn(prefix string, total, handled, addTotal, addHandled, statDay bool) string {
	var columns []string
	prefix = strings.Trim(prefix, "_")
	if prefix != PrefixTotal {
		prefix += "_"
	}

	if total {
		field := prefix + "total"
		columns = append(columns, "SUM("+field+") AS "+field)
	}
	if handled {
		field := prefix + "handled"
		columns = append(columns, "SUM("+field+") AS "+field)
	}
	if addTotal {
		field := prefix + "add_total"
		columns = append(columns, "SUM("+field+") AS "+field)

	}
	if addHandled {
		field := prefix + "add_handled"
		columns = append(columns, "SUM("+field+") AS "+field)
	}

	if statDay {
		columns = append(columns, "stat_day")
	}
	return strings.Join(columns, ", ")
}

// StatColumnOfAllType 统计表字段: 全部类型之和
func (crd *baseAlarm) StatColumnOfAllType(total, handled, addTotal, addHandled, statDay bool) string {
	var columns []string
	prefixes := []string{
		PrefixTotal,
		PrefixProcess,
		PrefixDll,
		PrefixBlackmail,
		PrefixTamper,
		PrefixPort,
		PrefixPlug,
		PrefixUsb,
		PrefixOutreach,
		PrefixLogin,
		PrefixSelf,
		PrefixVirus,
	}

	for _, val := range prefixes {
		columns = append(columns, crd.StatColumn(val, total, handled, addTotal, addHandled, statDay))
	}

	return strings.Join(columns, ", ")
}
