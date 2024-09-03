package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/outreach_models"
)

type OutreachTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *OutreachTrans) TransGlobalPage(total int64, models []*outreach_models.GlobalModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["title"] = model.Title
			data["ip"] = model.Ip
			data["is_default"] = model.IsDefault
			data["open_switch"] = model.OpenSwitch
			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["update_id"] = model.UpdateId
			data["update_name"] = model.UpdateName
			data["updated_at"] = model.UpdatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransDetectList 转换成主动探测列表需要的数据d
func (t *OutreachTrans) TransDetectList(models []*outreach_models.DetectModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["title"] = model.Title
			data["addr"] = model.Addr
			data["method"] = model.Method
			data["open_switch"] = model.OpenSwitch
			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["update_id"] = model.UpdateId
			data["update_name"] = model.UpdateName
			data["updated_at"] = model.UpdatedAt

			items = append(items, data)
		}
	}

	return items
}
