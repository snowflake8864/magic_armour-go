package proc_models

type SearchConditionsForProcess struct {
	OrgId        int64
	AgentId      string
	ProcHash     string
	DllHash      string
	Status       int
	SearchStatus int
	Dir          string
	GteCreateAt  string
	LteCreateAt  string
}

type SearchConditionsForSample struct {
	OrgId   int64
	AgentId string
	Id      string
	Status  int
}

type baseProc struct {
}
