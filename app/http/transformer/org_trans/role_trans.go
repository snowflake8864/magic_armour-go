package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Role struct {
	transformer.Common
	Lang string
}

func (t *Role) pageGeneral(model *role_models.RoleModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	return data
}

// TransPage 转换成管理页需要的数据
func (t *Role) TransPage(total int64, models []*role_models.RoleModel) interface{} {
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

// TransList 转换成管理页需要的数据
func (t *Role) TransList(models []*role_models.RoleModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return items
}
