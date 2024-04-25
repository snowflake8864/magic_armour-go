package software

// SoftSingle 软件单个
type SoftSingle struct {
	SoftId int `form:"softId" json:"softId" binding:"required"`
}

type PlatformType struct {
	PlatformType int `form:"platformType" json:"platformType" binding:"required,oneof=-100 1 2"`
}
