package proc_models

type SearchConditionsForProcess struct {
	OrgId        int64
	AgentId      string
	ProcHash     string
	DllHash      string
	Status       int
	SearchStatus int
	Dir          string
	CreateAtGte  string
	CreateAtLte  string
}

type baseProc struct {
}
