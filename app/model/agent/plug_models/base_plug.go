package plug_models

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	Ip         string
	OpenSwitch int
}

type basePlug struct {
}

// conditions2Where 查询条件转化
func (m *basePlug) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	return wheres
}
