package login_solve

type AddConfig struct {
	Gap       int `form:"gap" json:"gap" binding:"required,min=1,max=100"`
	Threshold int `form:"threshold" json:"threshold" binding:"required,min=1,max=100"`
}
