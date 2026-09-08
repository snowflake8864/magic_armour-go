package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type TaskTrans struct {
	transformer.Common
	Lang string
}

func (t *TaskTrans) pageGeneral(model *task_models.TaskModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["id"] = fmt.Sprintf("%d", model.Id)
	if t.Lang == consts.LanguageZh {
		data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_task.EnumMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_task.EnumMap)
	} else {
		data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_task.EnumMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_task.EnumMap)
	}

	return data
}

// TransPage 转换成管理页需要的数据
func (t *TaskTrans) TransPage(total int64, models []*task_models.TaskModel) interface{} {
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
