package plug

type IpAdd struct {
	Ip     string `form:"ip" json:"ip" binding:"required,max=128"`
	Remark string `form:"remark" json:"remark" binding:"required,max=128"`
}

type ConfigPortAdd struct {
	Weight    int `form:"weight" json:"weight" binding:"required,min=1,max=100"`
	CountTime int `form:"countTime" json:"countTime" binding:"required,min=1,max=120"`
	Type      int `form:"type" json:"type" binding:"required,oneof=1 2"`
	PlugTime  int `form:"plugTime" json:"plugTime" binding:"required_if=Type 1,min=0,max=1440"`
}

type ConfigLoginAdd struct {
	Type     int `form:"type" json:"type" binding:"required,oneof=1 2"`
	PlugTime int `form:"plugTime" json:"plugTime" binding:"required_if=Type 1,min=0,max=1440"`
}

type Id struct {
	Id int64 `form:"id" json:"id" binding:"required,min=1"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}
