package pw_jump

type Config struct {
	Switch  *int `form:"switch" json:"switch" binding:"required,oneof=0 1"`
	Combine int  `form:"combine" json:"combine" binding:"required,min=1,max=15"`
	Size    int  `form:"size" json:"size" binding:"required,min=6,max=24"`
	Cycle   int  `form:"cycle" json:"cycle" binding:"required,min=1,max=3"`
	//Timing  int  `form:"timing" json:"timing" binding:"required,min=1,max=30"`
}
