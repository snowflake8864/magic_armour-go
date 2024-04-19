package organize

// CommonReason 原因
type CommonReason struct {
	Reason string `form:"reason" json:"reason" binding:"required,max=200"`
}

// CommonSwitch 开关
type CommonSwitch struct {
	Switch *int `form:"switch" json:"switch" binding:"required,oneof=0 1"`
}
