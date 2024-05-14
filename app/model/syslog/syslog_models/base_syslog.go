package syslog_models

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	CreateAtGte string
	CreateAtLte string
}

type baseSyslog struct {
}

// conditions2Where 查询条件转化
func (m *baseSyslog) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	return wheres
}
