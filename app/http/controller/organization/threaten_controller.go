package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/slice_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// ThreatenController 威胁情报相关
type ThreatenController struct {
	ControllerBase
}

// GetLevels 通过hashes获取集合
func (u *ThreatenController) GetLevels(ctx *gin.Context) {
	jsonStr := ctx.GetString(consts.ValidatorPrefixOrg + "hashes")
	hashes := json_helper.ToSliceString(jsonStr)
	if hashes == nil || len(hashes) == 0 {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	CacheFact := threaten_models.CreateCacheFactory()
	maps := CacheFact.GetLevelsByIds(slice_helper.StringRemoveRepeatByMap(hashes))
	var result []map[string]interface{}
	if maps != nil && len(maps) > 0 {
		for id, level := range maps {
			result = append(result, map[string]interface{}{
				"id":    id,
				"level": level,
				"label": enums.Get1stIsStrAnd2ndIsIntReturnStr("level", level, enum_threaten.EnumMap),
			})
		}
	}

	response.SuccessReturnData(ctx, result)
}

// GetAttr 获取属性
func (u *ThreatenController) GetAttr(ctx *gin.Context) {
	hash := ctx.GetString(consts.ValidatorPrefixOrg + "hash")

	model := threaten_models.CreateCacheFactory().GetItemById(hash, "*")
	if model == nil || string_helper.IsEmpty(model.Id) {
		response.FailSimplify(ctx, my_errors.CurdDataNotExist)
		return
	}
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).TransThreatenInfo(model))
}
