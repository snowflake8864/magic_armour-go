package antivirus_engine

// EngineSingle 软件单个
type EngineSingle struct {
	EngineId int `form:"engineId" json:"engineId" binding:"required"`
}

type PlatformType struct {
	PlatformType int `form:"platformType" json:"platformType" binding:"required,oneof=-100 1 2"`
}
