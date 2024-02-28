package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// TaskController 任务相关
type TaskController struct {
}

// Get 获取任务
func (u *TaskController) Get(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
