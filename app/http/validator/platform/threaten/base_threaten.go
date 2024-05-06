package threaten

// threatenSingle 软件单个
type threatenSingle struct {
	Id string `form:"id" json:"id" binding:"required,max=64"`
}

type threatenSearch struct {
	Id  string `form:"id" json:"id" binding:"omitempty,max=64"`
	Dir string `form:"dir" json:"dir" binding:"omitempty,min=2"`
}

type PushBatch struct {
	List string `form:"list" json:"list" binding:"required,min=16"`
}
