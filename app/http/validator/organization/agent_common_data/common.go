package agent_common_data

// CommonSwitch 开关
type CommonSwitch struct {
	Switch *int `form:"switch" json:"switch" binding:"required,oneof=-100 0 1"`
}

// CommonMode 模式
type CommonMode struct {
	Mode int `form:"mode" json:"mode" binding:"required,oneof=-100 1 2"`
}
