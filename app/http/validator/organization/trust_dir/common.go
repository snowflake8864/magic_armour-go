package trust_dir

type Dir struct {
	Dir string `form:"dir" json:"dir" binding:"required,max=512"`
}

type PlatformType struct {
	PlatformType int `form:"platformType" json:"platformType" binding:"required,oneof=1 2"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=1"`
}

type IsExtend struct {
	IsExtend *int `form:"isExtend" json:"isExtend" binding:"required,oneof=0 1"`
}

type Name struct {
	Name string `form:"name" json:"name" binding:"required,max=64"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}
