package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_blackmail"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/backup_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type BackupTrans struct {
	transformer.Common
}

// TransHistoryPage 转换成备份历史分页需要的数据
func (t *BackupTrans) TransHistoryPage(total int64, models []*backup_models.BackupModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["backup_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("backup_status", model.BackupStatus, enum_blackmail.EnumMapInt)
			data["rollback_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("rollback_status", model.RollbackStatus, enum_blackmail.EnumMapInt)

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
