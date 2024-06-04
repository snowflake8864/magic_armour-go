package process

type ProcHash struct {
	ProcHash string `form:"procHash" json:"procHash" binding:"omitempty,min=8,max=32"`
}

type DllHash struct {
	DllHash string `form:"dllHash" json:"dllHash" binding:"omitempty,min=8,max=32"`
}

type Dir struct {
	Dir string `form:"dir" json:"dir" binding:"omitempty,min=2,max=512"`
}

type SearchStatus struct {
	SearchStatus int `form:"searchStatus" json:"searchStatus" binding:"required"`
}
