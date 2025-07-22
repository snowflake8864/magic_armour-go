package ip_jump

type GlobalConfig struct {
	Switch  *int   `form:"switch" json:"switch" binding:"required,oneof=0 1"`
	Ip      string `form:"ip" json:"ip" binding:"required,ipv4"`
	Mask    int    `form:"mask" json:"mask" binding:"required,min=8,max=32"`
	Gateway string `form:"gateway" json:"gateway" binding:"required,ipv4"`
	Cycle   int    `form:"cycle" json:"cycle" binding:"required,min=1,max=3"`
	Timing  int    `form:"timing" json:"timing" binding:"required,min=1,max=30"`
	Way     int    `form:"way" json:"way" binding:"required,min=1,max=2"`
}

type AgentConfig struct {
	Switch *int `form:"switch" json:"switch" binding:"required,oneof=0 1"`
	Cycle  int  `form:"cycle" json:"cycle" binding:"required,min=1,max=3"`
	Timing int  `form:"timing" json:"timing" binding:"required,min=1,max=30"`
	Way    int  `form:"way" json:"way" binding:"required,min=1,max=2"`
}
