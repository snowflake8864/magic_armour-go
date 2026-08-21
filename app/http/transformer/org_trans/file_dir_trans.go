package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type FileDirTrans struct {
	transformer.Common
	Lang string
}

// TranslPage 转换成名单分页需要的数据
func (t *FileDirTrans) TranslPage(total int64, models []*dir_models.DirModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)

			data["file_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.FileType, enum_tamper.EnumMapInt)

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
