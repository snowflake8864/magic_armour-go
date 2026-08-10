package anti_engine

type TaskId struct {
	TaskId int `form:"TaskId" json:"TaskId" binding:"required"`
}
