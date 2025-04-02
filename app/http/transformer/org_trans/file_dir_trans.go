package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/dir_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type FileDirTrans struct {
	transformer.Common
}

// TranslPage 转换成名单分页需要的数据
func (t *FileDirTrans) TranslPage(total int64, models []*dir_models.DirModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
