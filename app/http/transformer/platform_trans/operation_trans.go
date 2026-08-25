package platform_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Operation struct {
	transformer.Common
	Lang string
}

func (t *Operation) pageGeneral(model *operate_models.LogModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["id"] = fmt.Sprintf("%d", model.Id)

	if t.Lang == consts.LanguageZh {
		_, data["path_label"], _ = enum_platform_operate.GetModuleAndLabelByPath(model.Path)
		data["module_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("module", model.Module, enum_platform_operate.EnumMap)
	} else {
		_, _, data["path_label"] = enum_platform_operate.GetModuleAndLabelByPath(model.Path)
		data["module_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("module_en", model.Module, enum_platform_operate.EnumMap)
	}

	return data
}

// TransPage 转换成管理页需要的数据
func (t *Operation) TransPage(total int64, models []*operate_models.LogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
