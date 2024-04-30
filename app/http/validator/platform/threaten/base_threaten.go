package threaten

// threatenSingle 软件单个
type threatenSingle struct {
	Id string `form:"id" json:"id" binding:"required,max=64"`
}

type threatenSearch struct {
	Dir   string `form:"dir" json:"dir" binding:"omitempty,min=2"`
	Level *int   `form:"level" json:"level" binding:"required,max=20"`
}

type PushBatch struct {
	List string `form:"List" json:"List" binding:"required,min=16"`
}
