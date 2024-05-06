package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// DictionaryController 首页相关
type DictionaryController struct {
	ControllerBase
}

// List 保存
func (u *DictionaryController) List(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "List")
}

// Save 保存
func (u *DictionaryController) Save(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Save")
}
