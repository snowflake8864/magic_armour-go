package blackmail

type PlatformType struct {
	PlatformType int `form:"platformType" json:"platformType" binding:"required,oneof=1 2"`
}

type AddPolicy struct {
	Type    int    `form:"type" json:"type" binding:"required,oneof=1 2"`
	TypeVal string `form:"typeVal" json:"typeVal" binding:"required,max=512"`
	GroupId string `form:"groupId" json:"groupId" binding:"required_if=Type 1"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=1"`
}

type Remark struct {
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=128"`
}

type Processes struct {
	Processes string `form:"processes" json:"processes" binding:"required,min=32"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}

type Title struct {
	Title string `form:"title" json:"title" binding:"omitempty,max=128"`
}

type Suffix struct {
	Suffix string `form:"suffix" json:"suffix" binding:"required,max=64"`
}
