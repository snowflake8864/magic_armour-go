package organize

// CommonReason 原因
type CommonReason struct {
	Reason string `form:"reason" json:"reason" binding:"required,max=200"`
}
