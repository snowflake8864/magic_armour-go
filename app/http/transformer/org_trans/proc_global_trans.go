package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/proc_models"
)

type ProcGlobalTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *ProcGlobalTrans) TransPage(total int64, models []*proc_models.ProcGlobalModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["dir"] = model.Dir
			data["type"] = model.Type
			data["data_from"] = model.DataFrom
			data["open_switch"] = model.OpenSwitch
			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["updated_at"] = model.UpdatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
