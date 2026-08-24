package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/response"
)

// DictionaryController 首页相关
type DictionaryController struct {
	ControllerBase
}

// List 保存
func (u *DictionaryController) List(ctx *gin.Context) {
	models := dictionary_service.CreateCurdFactory().GetAll()
	response.SuccessReturnData(ctx, models)
}

// Save 保存
func (u *DictionaryController) Save(ctx *gin.Context) {
	id := ctx.GetString(consts.ValidatorPrefixPlatform + "id")
	val := ctx.GetString(consts.ValidatorPrefixPlatform + "val")

	if !dictionary_service.CreateCurdFactory().UpdateValById(id, val) {
		response.FailSimplify(ctx, my_errors.CurdSaveFailed)
		return
	}
	response.SuccessSimplify(ctx)
}
