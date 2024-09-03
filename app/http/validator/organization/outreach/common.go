package outreach

type AddIp struct {
	Title string `form:"title" json:"title" binding:"required,max=64"`
	Ip    string `form:"ip" json:"ip" binding:"required,max=128"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=1"`
}

type SIp struct {
	Ip string `form:"ip" json:"ip" binding:"omitempty,max=128"`
}

type AddAddr struct {
	Title string `form:"title" json:"title" binding:"required,max=64"`
	Addr  string `form:"addr" json:"addr" binding:"required,max=128"`
}
