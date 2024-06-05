package process

type HashOmitempty struct {
	Hash string `form:"hash" json:"hash" binding:"omitempty,min=8,max=32"`
}

type Dir struct {
	Dir string `form:"dir" json:"dir" binding:"omitempty,min=2,max=512"`
}

type Hash struct {
	Hash string `form:"hash" json:"hash" binding:"required,min=8,max=32"`
}

type Hashes struct {
	Hashes string `form:"hashes" json:"hashes" binding:"required"`
}

type SearchStatus struct {
	Status int `form:"status" json:"status" binding:"required"`
}

type ChangeSelfStatus struct {
	Status int `form:"status" json:"status" binding:"required,oneof=2 4 8"`
}

type ChangeGlobalStatus struct {
	Status int `form:"status" json:"status" binding:"required,oneof=16 32"`
}
