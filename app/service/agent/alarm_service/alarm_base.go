package alarm_service

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
)

type AlarmBase struct {
}

// search2Where 查询条件转化
func (crd *AlarmBase) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !helpers.StringIsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !helpers.StringIsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !helpers.StringIsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}
	return wheres
}

// StatColumn 统计表字段
// @use: StatColumn(PrefixTotal, true, true,true,true)
func (crd *AlarmBase) StatColumn(prefix string, total, handled, addTotal, addHandled bool) string {
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
	return strings.Join(columns, ", ")
}

// StatColumnOfAllType 统计表字段: 全部类型之和
func (crd *AlarmBase) StatColumnOfAllType(total, handled, addTotal, addHandled bool) string {
	var columns []string
	prefixes := []string{
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
	}

	for _, val := range prefixes {
		columns = append(columns, crd.StatColumn(val, total, handled, addTotal, addHandled))
	}

	return strings.Join(columns, ", ")
}
