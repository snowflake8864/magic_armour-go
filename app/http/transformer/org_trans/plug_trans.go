package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/plug_models"
)

type PlugTrans struct {
	transformer.Common
}

// TranIpPolicyPage 转换成终端ip名单分页需要的数据
func (t *PlugTrans) TranIpPolicyPage(total int64, models []*plug_models.IpPolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["ip"] = model.Ip
			data["remark"] = model.Remark
			data["direction"] = model.Direction
			data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)
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
